package gig

import (
	"bytes"
	"fmt"
	"go/token"
	"go/types"
	"reflect"
	"sort"
	"strings"

	xssa "golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/compiler"
	"github.com/t04dJ14n9/gig/importer"
	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/external"
)

// DebugDump compiles sourceCode and returns a readable SSA and bytecode dump.
//
// Unlike Build, DebugDump does not execute init functions. It only runs the
// compile pipeline through SSA construction and bytecode generation.
func DebugDump(sourceCode string, opts ...BuildOption) (string, error) {
	cfg := debugBuildConfig(opts)
	result, err := compiler.Build(sourceCode, cfg.registry, debugCompilerOptions(cfg)...)
	if err != nil {
		return "", err
	}

	var out strings.Builder
	out.WriteString("# Gig Debug Dump\n\n")
	writeSSADump(&out, result.SSAPkg)
	writeBytecodeDump(&out, result.Program)
	return out.String(), nil
}

func debugBuildConfig(opts []BuildOption) buildConfig {
	cfg := buildConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.registry == nil {
		cfg.registry = importer.GlobalRegistry()
	}
	return cfg
}

func debugCompilerOptions(cfg buildConfig) []compiler.BuildOption {
	var opts []compiler.BuildOption
	if cfg.allowPanic {
		opts = append(opts, compiler.WithAllowPanic())
	}
	if cfg.allowUnsafeTypePass {
		opts = append(opts, compiler.WithAllowUnsafeTypePass())
	}
	return opts
}

func writeSSADump(out *strings.Builder, pkg *xssa.Package) {
	out.WriteString("## SSA\n\n")
	if pkg == nil {
		out.WriteString("(no SSA package)\n\n")
		return
	}

	fmt.Fprintf(out, "package %s", pkg.Pkg.Name())
	if path := pkg.Pkg.Path(); path != "" && path != pkg.Pkg.Name() {
		fmt.Fprintf(out, " (%s)", path)
	}
	out.WriteString("\n\n")

	for _, member := range sortedSSAMembers(pkg) {
		writeSSAMember(out, pkg, member)
	}
}

func sortedSSAMembers(pkg *xssa.Package) []xssa.Member {
	names := make([]string, 0, len(pkg.Members))
	for name := range pkg.Members {
		names = append(names, name)
	}
	sort.Strings(names)

	members := make([]xssa.Member, 0, len(names))
	for _, name := range names {
		members = append(members, pkg.Members[name])
	}
	return members
}

func writeSSAMember(out *strings.Builder, pkg *xssa.Package, member xssa.Member) {
	switch m := member.(type) {
	case *xssa.Function:
		writeSSAFunction(out, m)
	default:
		fmt.Fprintf(out, "### %s %s\n\n", tokenName(member.Token()), member.Name())
		fmt.Fprintf(out, "%s: %s\n\n", member.RelString(pkg.Pkg), member.Type())
	}
}

func writeSSAFunction(out *strings.Builder, fn *xssa.Function) {
	fmt.Fprintf(out, "### Function %s\n\n", fn.Name())
	if fn.Synthetic != "" {
		fmt.Fprintf(out, "synthetic: %s\n\n", fn.Synthetic)
	}

	var buf bytes.Buffer
	_, _ = fn.WriteTo(&buf)
	out.Write(buf.Bytes())
	if !strings.HasSuffix(out.String(), "\n") {
		out.WriteByte('\n')
	}
	out.WriteByte('\n')
}

func tokenName(tok token.Token) string {
	if tok == token.ILLEGAL {
		return "member"
	}
	return strings.ToLower(tok.String())
}

func writeBytecodeDump(out *strings.Builder, prog *bytecode.CompiledProgram) {
	out.WriteString("## Bytecode\n\n")
	if prog == nil {
		out.WriteString("(no bytecode program)\n\n")
		return
	}

	writeFunctionIndex(out, prog)
	for _, fn := range sortedCompiledFunctions(prog) {
		writeCompiledFunction(out, prog, fn)
	}
	writeConstants(out, prog.Constants)
	writeGlobals(out, prog.Globals)
	writeTypes(out, prog.Types)
}

func writeFunctionIndex(out *strings.Builder, prog *bytecode.CompiledProgram) {
	out.WriteString("### Function index\n\n")
	if len(prog.FuncByIndex) == 0 {
		out.WriteString("(empty)\n\n")
		return
	}
	for idx, fn := range prog.FuncByIndex {
		if fn == nil {
			fmt.Fprintf(out, "[%d] <nil>\n", idx)
			continue
		}
		fmt.Fprintf(out, "[%d] %s\n", idx, fn.Name)
	}
	out.WriteByte('\n')
}

func sortedCompiledFunctions(prog *bytecode.CompiledProgram) []*bytecode.CompiledFunction {
	funcs := make([]*bytecode.CompiledFunction, 0, len(prog.Functions))
	seen := make(map[*bytecode.CompiledFunction]bool, len(prog.Functions))
	for _, fn := range prog.Functions {
		if fn == nil || seen[fn] {
			continue
		}
		funcs = append(funcs, fn)
		seen[fn] = true
	}
	sort.Slice(funcs, func(i, j int) bool {
		if funcs[i].Name == funcs[j].Name {
			return funcs[i].FuncIdx < funcs[j].FuncIdx
		}
		return funcs[i].Name < funcs[j].Name
	})
	return funcs
}

func writeCompiledFunction(out *strings.Builder, prog *bytecode.CompiledProgram, fn *bytecode.CompiledFunction) {
	fmt.Fprintf(out, "### Function %s\n\n", fn.Name)
	fmt.Fprintf(
		out,
		"locals: %d, params: %d, free: %d, maxStack: %d, funcIdx: %d, hasIntLocals: %t\n",
		fn.NumLocals,
		fn.NumParams,
		fn.NumFreeVars,
		fn.MaxStack,
		fn.FuncIdx,
		fn.HasIntLocals,
	)
	if fn.HasReceiver {
		fmt.Fprintf(out, "receiver: %s, pointer: %t\n", fn.ReceiverTypeName, fn.ReceiverIsPointer)
	}
	if len(fn.ParamTypes) > 0 {
		out.WriteString("params:\n")
		for idx, typ := range fn.ParamTypes {
			fmt.Fprintf(out, "  [%d] %s\n", idx, typ)
		}
	}
	out.WriteString("\ninstructions:\n")
	writeInstructions(out, prog, fn.Instructions)
	out.WriteByte('\n')
}

func writeInstructions(out *strings.Builder, prog *bytecode.CompiledProgram, code []byte) {
	for pc := 0; pc < len(code); {
		next := writeInstruction(out, prog, code, pc)
		if next <= pc {
			next = pc + 1
		}
		pc = next
	}
}

func writeInstruction(out *strings.Builder, prog *bytecode.CompiledProgram, code []byte, pc int) int {
	op := bytecode.OpCode(code[pc])
	width := bytecode.OperandWidth(op)
	end := pc + 1 + width

	if end > len(code) {
		fmt.Fprintf(out, "  %04d  %-34s <truncated: need %d operand bytes, have %d>\n", pc, op, width, len(code)-pc-1)
		return len(code)
	}

	operands := code[pc+1 : end]
	operandText := formatOperands(operands)
	annotation := instructionAnnotation(op, operands, prog)
	if annotation != "" {
		annotation = " ; " + annotation
	}
	if operandText == "" && annotation == "" {
		fmt.Fprintf(out, "  %04d  %s\n", pc, op)
	} else if operandText == "" {
		fmt.Fprintf(out, "  %04d  %-34s%s\n", pc, op, annotation)
	} else {
		fmt.Fprintf(out, "  %04d  %-34s %s%s\n", pc, op, operandText, annotation)
	}
	return end
}

func formatOperands(operands []byte) string {
	switch len(operands) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d", operands[0])
	case 2:
		return fmt.Sprintf("%d", readUint16(operands, 0))
	case 3:
		return fmt.Sprintf("%d, %d", readUint16(operands, 0), operands[2])
	default:
		parts := make([]string, 0, len(operands)/2+1)
		i := 0
		for ; i+1 < len(operands); i += 2 {
			parts = append(parts, fmt.Sprintf("%d", readUint16(operands, i)))
		}
		if i < len(operands) {
			parts = append(parts, fmt.Sprintf("%d", operands[i]))
		}
		return strings.Join(parts, ", ")
	}
}

func instructionAnnotation(op bytecode.OpCode, operands []byte, prog *bytecode.CompiledProgram) string {
	switch op {
	case bytecode.OpConst:
		return constantAnnotation(readUint16(operands, 0), prog)
	case bytecode.OpLocal, bytecode.OpSetLocal, bytecode.OpAddr, bytecode.OpIntLocal, bytecode.OpIntSetLocal:
		return localAnnotation(readUint16(operands, 0))
	case bytecode.OpFree, bytecode.OpSetFree:
		return freeAnnotation(operands[0])
	case bytecode.OpAddLocalLocal, bytecode.OpSubLocalLocal, bytecode.OpMulLocalLocal:
		return localPairAnnotation(readUint16(operands, 0), readUint16(operands, 2))
	case bytecode.OpAddLocalConst, bytecode.OpSubLocalConst:
		return localConstAnnotation(readUint16(operands, 0), readUint16(operands, 2), prog)
	case bytecode.OpLocalConstAddSetLocal, bytecode.OpLocalConstSubSetLocal, bytecode.OpLocalConstMulSetLocal,
		bytecode.OpIntLocalConstAddSetLocal, bytecode.OpIntLocalConstSubSetLocal,
		bytecode.OpIntLocalConstMulSetLocal:
		return localConstSetAnnotation(readUint16(operands, 0), readUint16(operands, 2), readUint16(operands, 4), prog)
	case bytecode.OpLocalLocalAddSetLocal, bytecode.OpLocalLocalSubSetLocal, bytecode.OpLocalLocalMulSetLocal,
		bytecode.OpIntLocalLocalAddSetLocal, bytecode.OpIntLocalLocalSubSetLocal,
		bytecode.OpIntLocalLocalMulSetLocal:
		return localLocalSetAnnotation(readUint16(operands, 0), readUint16(operands, 2), readUint16(operands, 4))
	case bytecode.OpGlobal, bytecode.OpSetGlobal:
		return globalAnnotation(readUint16(operands, 0), prog)
	case bytecode.OpCall, bytecode.OpGoCall:
		return functionCallAnnotation(readUint16(operands, 0), operands[2], prog)
	case bytecode.OpClosure:
		return functionAnnotation(readUint16(operands, 0), prog)
	case bytecode.OpCallExternal, bytecode.OpGoCallExternal:
		return externalCallAnnotation(readUint16(operands, 0), operands[2], prog)
	case bytecode.OpDefer:
		return functionAnnotation(readUint16(operands, 0), prog)
	case bytecode.OpDeferExternal:
		return externalCallAnnotation(readUint16(operands, 0), operands[len(operands)-1], prog)
	case bytecode.OpJump, bytecode.OpJumpTrue, bytecode.OpJumpFalse:
		return fmt.Sprintf("target=%d", readUint16(operands, 0))
	case bytecode.OpConvert, bytecode.OpAssert, bytecode.OpNew:
		return typeAnnotation(readUint16(operands, 0), prog)
	case bytecode.OpMakeInterface, bytecode.OpChangeType:
		return typePairAnnotation(readUint16(operands, 0), readUint16(operands, 2), prog)
	default:
		return ""
	}
}

func localAnnotation(idx uint16) string {
	return fmt.Sprintf("local[%d]", idx)
}

func freeAnnotation(idx byte) string {
	return fmt.Sprintf("free[%d]", idx)
}

func localPairAnnotation(first, second uint16) string {
	return fmt.Sprintf("local[%d], local[%d]", first, second)
}

func localConstAnnotation(localIdx, constIdx uint16, prog *bytecode.CompiledProgram) string {
	return fmt.Sprintf("local[%d], %s", localIdx, constantAnnotation(constIdx, prog))
}

func localConstSetAnnotation(localIdx, constIdx, dstIdx uint16, prog *bytecode.CompiledProgram) string {
	return fmt.Sprintf("local[%d], %s -> local[%d]", localIdx, constantAnnotation(constIdx, prog), dstIdx)
}

func localLocalSetAnnotation(leftIdx, rightIdx, dstIdx uint16) string {
	return fmt.Sprintf("local[%d], local[%d] -> local[%d]", leftIdx, rightIdx, dstIdx)
}

func readUint16(bytes []byte, offset int) uint16 {
	return uint16(bytes[offset])<<8 | uint16(bytes[offset+1])
}

func constantAnnotation(idx uint16, prog *bytecode.CompiledProgram) string {
	if prog == nil || int(idx) >= len(prog.Constants) {
		return fmt.Sprintf("const[%d] <out of range>", idx)
	}
	return fmt.Sprintf("const[%d] = %s", idx, describeConstant(prog.Constants[idx]))
}

func globalAnnotation(idx uint16, prog *bytecode.CompiledProgram) string {
	if prog == nil {
		return fmt.Sprintf("global[%d]", idx)
	}
	names := make([]string, 0, len(prog.Globals))
	for name, globalIdx := range prog.Globals {
		if globalIdx == int(idx) {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	if len(names) == 0 {
		return fmt.Sprintf("global[%d]", idx)
	}
	return fmt.Sprintf("global[%d] = %s", idx, strings.Join(names, ", "))
}

func functionCallAnnotation(idx uint16, numArgs byte, prog *bytecode.CompiledProgram) string {
	return fmt.Sprintf("%s, args=%d", functionAnnotation(idx, prog), numArgs)
}

func functionAnnotation(idx uint16, prog *bytecode.CompiledProgram) string {
	if prog == nil || int(idx) >= len(prog.FuncByIndex) || prog.FuncByIndex[idx] == nil {
		return fmt.Sprintf("fn[%d] <out of range>", idx)
	}
	return fmt.Sprintf("fn[%d] = %s", idx, prog.FuncByIndex[idx].Name)
}

func externalCallAnnotation(idx uint16, numArgs byte, prog *bytecode.CompiledProgram) string {
	return fmt.Sprintf("%s, args=%d", constantAnnotation(idx, prog), numArgs)
}

func typeAnnotation(idx uint16, prog *bytecode.CompiledProgram) string {
	if prog == nil || int(idx) >= len(prog.Types) {
		return fmt.Sprintf("type[%d] <out of range>", idx)
	}
	return fmt.Sprintf("type[%d] = %s", idx, prog.Types[idx])
}

func typePairAnnotation(first, second uint16, prog *bytecode.CompiledProgram) string {
	return fmt.Sprintf("%s, %s", typeAnnotation(first, prog), typeAnnotation(second, prog))
}

func describeConstant(v any) string {
	switch c := v.(type) {
	case nil:
		return "nil"
	case string:
		return fmt.Sprintf("%q (string)", c)
	case *external.ExternalFuncInfo:
		return fmt.Sprintf("external func %s.%s direct=%t variadic=%t numIn=%d", c.PkgPath, c.FuncName, c.DirectCall != nil, c.IsVariadic, c.NumIn)
	case *external.ExternalMethodInfo:
		return fmt.Sprintf("external method %s.%s direct=%t receiver=%s", c.PkgPath, c.MethodName, c.DirectCall != nil, c.ReceiverTypeName)
	default:
		if reflect.TypeOf(v) != nil && reflect.TypeOf(v).Kind() == reflect.Func {
			return fmt.Sprintf("%T", v)
		}
		return truncateDebugText(fmt.Sprintf("%#v (%T)", v, v), 160)
	}
}

func writeConstants(out *strings.Builder, constants []any) {
	out.WriteString("## Constants\n\n")
	if len(constants) == 0 {
		out.WriteString("(empty)\n\n")
		return
	}
	for idx, c := range constants {
		fmt.Fprintf(out, "[%d] %s\n", idx, describeConstant(c))
	}
	out.WriteByte('\n')
}

func writeGlobals(out *strings.Builder, globals map[string]int) {
	out.WriteString("## Globals\n\n")
	if len(globals) == 0 {
		out.WriteString("(empty)\n\n")
		return
	}
	type globalEntry struct {
		name string
		idx  int
	}
	entries := make([]globalEntry, 0, len(globals))
	for name, idx := range globals {
		entries = append(entries, globalEntry{name: name, idx: idx})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].idx == entries[j].idx {
			return entries[i].name < entries[j].name
		}
		return entries[i].idx < entries[j].idx
	})
	for _, entry := range entries {
		fmt.Fprintf(out, "[%d] %s\n", entry.idx, entry.name)
	}
	out.WriteByte('\n')
}

func writeTypes(out *strings.Builder, types []types.Type) {
	out.WriteString("## Types\n\n")
	if len(types) == 0 {
		out.WriteString("(empty)\n\n")
		return
	}
	for idx, typ := range types {
		fmt.Fprintf(out, "[%d] %s\n", idx, typ)
	}
	out.WriteByte('\n')
}

func truncateDebugText(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", `\n`)
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}
