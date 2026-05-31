package bytecode

// unknownOpName is the string returned for unrecognized opcodes.
const unknownOpName = "UNKNOWN"

type opcodeSpec struct {
	name  string
	width int
}

// opcodeSpecs is the authoritative opcode metadata table. It feeds both
// human-readable names and operand widths so the two cannot drift apart.
var opcodeSpecs = [256]opcodeSpec{
	OpAdd:                          {name: "ADD"},
	OpAddLocalConst:                {name: "ADDLOCALCONST", width: 4},
	OpAddLocalLocal:                {name: "ADDLOCALLOCAL", width: 4},
	OpAddSetLocal:                  {name: "ADDSETLOCAL", width: 2},
	OpAddr:                         {name: "ADDR", width: 2},
	OpAnd:                          {name: "AND"},
	OpAndNot:                       {name: "ANDNOT"},
	OpAppend:                       {name: "APPEND"},
	OpAssert:                       {name: "ASSERT", width: 2},
	OpCall:                         {name: "CALL", width: 3},
	OpCallExternal:                 {name: "CALLEXTERNAL", width: 3},
	OpCallIndirect:                 {name: "CALLINDIRECT", width: 1},
	OpCap:                          {name: "CAP"},
	OpChangeType:                   {name: "CHANGETYPE", width: 4},
	OpClose:                        {name: "CLOSE"},
	OpClosure:                      {name: "CLOSURE", width: 3},
	OpComplex:                      {name: "COMPLEX"},
	OpConst:                        {name: "CONST", width: 2},
	OpConvert:                      {name: "CONVERT", width: 2},
	OpCopy:                         {name: "COPY"},
	OpDefer:                        {name: "DEFER", width: 2},
	OpDeferExternal:                {name: "DEFEREXTERNAL", width: 3},
	OpDeferIndirect:                {name: "DEFERINDIRECT", width: 2},
	OpDelete:                       {name: "DELETE"},
	OpDeref:                        {name: "DEREF"},
	OpDiv:                          {name: "DIV"},
	OpDup:                          {name: "DUP"},
	OpEqual:                        {name: "EQUAL"},
	OpFalse:                        {name: "FALSE"},
	OpField:                        {name: "FIELD", width: 2},
	OpFieldAddr:                    {name: "FIELDADDR", width: 2},
	OpFree:                         {name: "FREE", width: 1},
	OpGlobal:                       {name: "GLOBAL", width: 2},
	OpGoCall:                       {name: "GOCALL", width: 3},
	OpGoCallExternal:               {name: "GOCALLEXTERNAL", width: 3},
	OpGoCallIndirect:               {name: "GOCALLINDIRECT", width: 1},
	OpGreater:                      {name: "GREATER"},
	OpGreaterEq:                    {name: "GREATEREQ"},
	OpGreaterLocalLocalJumpTrue:    {name: "GREATERLOCALLOCALJUMPTRUE", width: 6},
	OpHalt:                         {name: "HALT"},
	OpImag:                         {name: "IMAG"},
	OpIndex:                        {name: "INDEX"},
	OpIndexAddr:                    {name: "INDEXADDR"},
	OpIndexOk:                      {name: "INDEXOK"},
	OpIntGreaterLocalLocalJumpTrue: {name: "INTGREATERLOCALLOCALJUMPTRUE", width: 6},
	OpIntLessEqLocalConstJumpFalse: {name: "INTLESSEQLOCALCONSTJUMPFALSE", width: 6},
	OpIntLessEqLocalConstJumpTrue:  {name: "INTLESSEQLOCALCONSTJUMPTRUE", width: 6},
	OpIntLessLocalConstJumpFalse:   {name: "INTLESSLOCALCONSTJUMPFALSE", width: 6},
	OpIntLessLocalConstJumpTrue:    {name: "INTLESSLOCALCONSTJUMPTRUE", width: 6},
	OpIntLessLocalLocalJumpFalse:   {name: "INTLESSLOCALLOCALJUMPFALSE", width: 6},
	OpIntLessLocalLocalJumpTrue:    {name: "INTLESSLOCALLOCALJUMPTRUE", width: 6},
	OpIntLocal:                     {name: "INTLOCAL", width: 2},
	OpIntLocalConstAddSetLocal:     {name: "INTLOCALCONSTADDSETLOCAL", width: 6},
	OpIntLocalConstMulSetLocal:     {name: "INTLOCALCONSTMULSETLOCAL", width: 6},
	OpIntLocalConstSubSetLocal:     {name: "INTLOCALCONSTSUBSETLOCAL", width: 6},
	OpIntLocalLocalAddSetLocal:     {name: "INTLOCALLOCALADDSETLOCAL", width: 6},
	OpIntLocalLocalMulSetLocal:     {name: "INTLOCALLOCALMULSETLOCAL", width: 6},
	OpIntLocalLocalSubSetLocal:     {name: "INTLOCALLOCALSUBSETLOCAL", width: 6},
	OpIntMoveLocal:                 {name: "INTMOVELOCAL", width: 4},
	OpIntSetLocal:                  {name: "INTSETLOCAL", width: 2},
	OpIntSliceGet:                  {name: "INTSLICEGET", width: 6},
	OpIntSliceSet:                  {name: "INTSLICESET", width: 6},
	OpIntSliceSetConst:             {name: "INTSLICESETCONST", width: 6},
	OpJump:                         {name: "JUMP", width: 2},
	OpJumpFalse:                    {name: "JUMPFALSE", width: 2},
	OpJumpTrue:                     {name: "JUMPTRUE", width: 2},
	OpLen:                          {name: "LEN"},
	OpLess:                         {name: "LESS"},
	OpLessEq:                       {name: "LESSEQ"},
	OpLessEqLocalConstJumpFalse:    {name: "LESSEQLOCALCONSTJUMPFALSE", width: 6},
	OpLessEqLocalConstJumpTrue:     {name: "LESSEQLOCALCONSTJUMPTRUE", width: 6},
	OpLessLocalConstJumpFalse:      {name: "LESSLOCALCONSTJUMPFALSE", width: 6},
	OpLessLocalConstJumpTrue:       {name: "LESSLOCALCONSTJUMPTRUE", width: 6},
	OpLessLocalLocalJumpFalse:      {name: "LESSLOCALLOCALJUMPFALSE", width: 6},
	OpLessLocalLocalJumpTrue:       {name: "LESSLOCALLOCALJUMPTRUE", width: 6},
	OpLocal:                        {name: "LOCAL", width: 2},
	OpLocalConstAddSetLocal:        {name: "LOCALCONSTADDSETLOCAL", width: 6},
	OpLocalConstMulSetLocal:        {name: "LOCALCONSTMULSETLOCAL", width: 6},
	OpLocalConstSubSetLocal:        {name: "LOCALCONSTSUBSETLOCAL", width: 6},
	OpLocalLocalAddSetLocal:        {name: "LOCALLOCALADDSETLOCAL", width: 6},
	OpLocalLocalMulSetLocal:        {name: "LOCALLOCALMULSETLOCAL", width: 6},
	OpLocalLocalSubSetLocal:        {name: "LOCALLOCALSUBSETLOCAL", width: 6},
	OpLsh:                          {name: "LSH"},
	OpMakeChan:                     {name: "MAKECHAN"},
	OpMakeInterface:                {name: "MAKEINTERFACE", width: 4},
	OpMakeMap:                      {name: "MAKEMAP"},
	OpMakeSlice:                    {name: "MAKESLICE"},
	OpMod:                          {name: "MOD"},
	OpMul:                          {name: "MUL"},
	OpMulLocalLocal:                {name: "MULLOCALLOCAL", width: 4},
	OpNeg:                          {name: "NEG"},
	OpNew:                          {name: "NEW", width: 2},
	OpNil:                          {name: "NIL"},
	OpNot:                          {name: "NOT"},
	OpNotEqual:                     {name: "NOTEQUAL"},
	OpOr:                           {name: "OR"},
	OpPack:                         {name: "PACK", width: 2},
	OpPanic:                        {name: "PANIC"},
	OpPop:                          {name: "POP"},
	OpPrint:                        {name: "PRINT", width: 1},
	OpPrintln:                      {name: "PRINTLN", width: 1},
	OpRange:                        {name: "RANGE"},
	OpRangeNext:                    {name: "RANGENEXT"},
	OpReal:                         {name: "REAL"},
	OpRecover:                      {name: "RECOVER"},
	OpRecv:                         {name: "RECV"},
	OpRecvOk:                       {name: "RECVOK"},
	OpReturn:                       {name: "RETURN"},
	OpReturnVal:                    {name: "RETURNVAL"},
	OpRsh:                          {name: "RSH"},
	OpRunDefers:                    {name: "RUNDEFERS"},
	OpSelect:                       {name: "SELECT", width: 2},
	OpSend:                         {name: "SEND"},
	OpSetDeref:                     {name: "SETDEREF"},
	OpSetField:                     {name: "SETFIELD", width: 2},
	OpSetFree:                      {name: "SETFREE", width: 1},
	OpSetGlobal:                    {name: "SETGLOBAL", width: 2},
	OpSetIndex:                     {name: "SETINDEX"},
	OpSetLocal:                     {name: "SETLOCAL", width: 2},
	OpSlice:                        {name: "SLICE"},
	OpSub:                          {name: "SUB"},
	OpSubLocalConst:                {name: "SUBLOCALCONST", width: 4},
	OpSubLocalLocal:                {name: "SUBLOCALLOCAL", width: 4},
	OpSubSetLocal:                  {name: "SUBSETLOCAL", width: 2},
	OpTrue:                         {name: "TRUE"},
	OpUnpack:                       {name: "UNPACK"},
	OpXor:                          {name: "XOR"},
}

var opNameTable, operandWidthTable = buildOpcodeTables()

func buildOpcodeTables() ([256]string, [256]int) {
	var names [256]string
	for i := range names {
		names[i] = unknownOpName
	}

	var widths [256]int
	for op, spec := range opcodeSpecs {
		if spec.name != "" {
			names[op] = spec.name
		}
		widths[op] = spec.width
	}
	return names, widths
}

// String returns the name of the opcode as a human-readable string.
func (op OpCode) String() string {
	return opNameTable[op]
}

// OperandWidth returns the operand byte width for an opcode using O(1) array lookup.
func OperandWidth(op OpCode) int {
	return operandWidthTable[op]
}
