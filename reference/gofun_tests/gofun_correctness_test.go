package gofun

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	newgofun "git.code.oa.com/datacenter/onefun/gofun"
	_ "git.code.oa.com/datacenter/onefun/gofun/packages"
)

// ============================================================================
// 使用 Gig 的 correctness_test.go 测试用例跑 gofun，统计失败数量
// ============================================================================

// testFunc 定义一个测试函数
type testFunc struct {
	pkg      string
	funcName string
}

// extractExportedFuncs 从源码中提取所有导出函数名
func extractExportedFuncs(src string) []string {
	var funcs []string
	for _, line := range strings.Split(src, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "func ") && !strings.HasPrefix(line, "func (") {
			// 提取函数名
			rest := line[5:]
			idx := strings.IndexByte(rest, '(')
			if idx <= 0 {
				continue
			}
			name := rest[:idx]
			// 跳过 main, init
			if name == "main" || name == "init" {
				continue
			}
			// 只要导出函数 (首字母大写)
			if len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z' {
				// 跳过带参数的函数（只测试无参函数）
				params := rest[idx+1:]
				closeIdx := strings.IndexByte(params, ')')
				if closeIdx >= 0 {
					paramStr := strings.TrimSpace(params[:closeIdx])
					if paramStr == "" {
						funcs = append(funcs, name)
					}
				}
			}
		}
	}
	return funcs
}

// TestGofunCorrectness 使用 gig 的 testdata 跑 gofun，统计 Build 和 Run 失败
func TestGofunCorrectness(t *testing.T) {
	// testdata 目录
	testdataDir := "../../tests/testdata"

	// 要测试的包（跳过 goroutine/channels/panic_recover 等并发/panic 相关的）
	packages := []string{
		"advanced", "algorithms", "arithmetic", "bitwise",
		"closures", "closures_advanced",
		"controlflow", "cornercases", "edgecases",
		"functions",
		"leetcode_hard", "mapadvanced", "maps", "multiassign",
		"namedreturn", "recursion",
		"scope", "slices", "slicing", "strings_pkg",
		"switch", "typeconv", "variables",
	}
	// 跳过: complex, structs, tricky (unexported struct fields → panic in typeChange)
	// 跳过: goroutine, channels, panic_recover (并发/panic 相关)
	// 跳过: external (用了 strings.ReplaceAll，gofun 未注册)
	// 跳过: initialize (unexported struct fields → panic)
	// 跳过: resolved_issue (用了 encoding/json，gofun 未注册)

	var totalTests int
	var buildFails int
	var runFails int
	var passed int

	type failInfo struct {
		pkg      string
		funcName string
		phase    string // "build" or "run"
		err      string
	}
	var failures []failInfo

	for _, pkg := range packages {
		srcPath := filepath.Join(testdataDir, pkg, "main.go")
		srcBytes, err := os.ReadFile(srcPath)
		if err != nil {
			t.Logf("跳过 %s: %v", pkg, err)
			continue
		}
		src := string(srcBytes)

		funcs := extractExportedFuncs(src)

		// 尝试 Build (加 recover 防止 panic 崩溃进程)
		var program *newgofun.Program
		var buildErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					buildErr = fmt.Errorf("panic: %v", r)
				}
			}()
			program, buildErr = newgofun.Build(src)
		}()
		if buildErr != nil {
			for _, fn := range funcs {
				totalTests++
				buildFails++
				failures = append(failures, failInfo{pkg, fn, "build", buildErr.Error()})
			}
			continue
		}

		// Build 成功，尝试 Run 每个无参导出函数
		for _, fn := range funcs {
			totalTests++
			var runErr error
			done := make(chan struct{})
			go func() {
				defer func() {
					if r := recover(); r != nil {
						runErr = fmt.Errorf("panic: %v", r)
					}
					close(done)
				}()
				_, runErr = program.Run(fn)
			}()
			select {
			case <-done:
			case <-time.After(3 * time.Second):
				runErr = fmt.Errorf("timeout (3s)")
			}
			if runErr != nil {
				runFails++
				failures = append(failures, failInfo{pkg, fn, "run", runErr.Error()})
			} else {
				passed++
			}
		}
	}

	// 按包排序
	sort.Slice(failures, func(i, j int) bool {
		if failures[i].pkg != failures[j].pkg {
			return failures[i].pkg < failures[j].pkg
		}
		return failures[i].funcName < failures[j].funcName
	})

	// 输出统计
	t.Logf("")
	t.Logf("========== gofun 正确性测试统计 ==========")
	t.Logf("总测试数: %d", totalTests)
	t.Logf("通过:     %d (%.1f%%)", passed, float64(passed)*100/float64(totalTests))
	t.Logf("Build失败: %d", buildFails)
	t.Logf("Run失败:   %d", runFails)
	t.Logf("总失败:   %d (%.1f%%)", buildFails+runFails, float64(buildFails+runFails)*100/float64(totalTests))
	t.Logf("")

	// 输出失败列表
	if len(failures) > 0 {
		t.Logf("========== 失败列表 ==========")
		lastPkg := ""
		for _, f := range failures {
			if f.pkg != lastPkg {
				t.Logf("")
				t.Logf("--- %s ---", f.pkg)
				lastPkg = f.pkg
			}
			// 截断错误信息
			errMsg := f.err
			if len(errMsg) > 120 {
				errMsg = errMsg[:120] + "..."
			}
			t.Logf("  [%s] %s: %s", f.phase, f.funcName, errMsg)
		}
	}

	// 按错误类型统计
	t.Logf("")
	t.Logf("========== 错误类型统计 ==========")
	errTypes := make(map[string]int)
	for _, f := range failures {
		errKey := f.err
		// 提取有意义的错误类别
		if strings.Contains(errKey, "undefined:") {
			idx := strings.Index(errKey, "undefined:")
			errKey = errKey[idx:]
			if endIdx := strings.Index(errKey, "\n"); endIdx > 0 {
				errKey = errKey[:endIdx]
			}
		} else if strings.Contains(errKey, "unexpected instruction:") {
			errKey = "unexpected instruction: *ssa.Index"
		} else if strings.Contains(errKey, "unexported but missing PkgPath") {
			errKey = "reflect.StructOf: unexported field missing PkgPath"
		} else if strings.Contains(errKey, "reflect: call of") {
			idx := strings.Index(errKey, "reflect: call of")
			errKey = errKey[idx:]
		} else if strings.Contains(errKey, "timeout") {
			errKey = "timeout (3s)"
		} else if len(errKey) > 80 {
			errKey = errKey[:80]
		}
		errTypes[errKey]++
	}
	type errCount struct {
		err   string
		count int
	}
	var sortedErrs []errCount
	for k, v := range errTypes {
		sortedErrs = append(sortedErrs, errCount{k, v})
	}
	sort.Slice(sortedErrs, func(i, j int) bool {
		return sortedErrs[i].count > sortedErrs[j].count
	})
	for _, e := range sortedErrs {
		t.Logf("  %3d × %s", e.count, e.err)
	}

	// 输出 Markdown 表格
	t.Logf("")
	t.Logf("========== 按包统计 ==========")
	t.Logf("| 包 | 总数 | 通过 | 失败 | 通过率 |")
	t.Logf("|-----|------|------|------|--------|")
	for _, pkg := range packages {
		pkgTotal := 0
		pkgFail := 0
		for _, f := range failures {
			if f.pkg == pkg {
				pkgFail++
			}
		}
		srcPath := filepath.Join(testdataDir, pkg, "main.go")
		srcBytes, err := os.ReadFile(srcPath)
		if err != nil {
			continue
		}
		funcs := extractExportedFuncs(string(srcBytes))
		pkgTotal = len(funcs)
		pkgPass := pkgTotal - pkgFail
		if pkgTotal == 0 {
			continue
		}
		pct := float64(pkgPass) * 100 / float64(pkgTotal)
		marker := ""
		if pkgFail > 0 {
			marker = " ❌"
		} else {
			marker = " ✅"
		}
		t.Logf("| %s | %d | %d | %d | %.0f%%%s |", pkg, pkgTotal, pkgPass, pkgFail, pct, marker)
	}

	// 最终摘要
	t.Logf("")
	t.Logf("gofun 通过 %d/%d 测试 (%.1f%%)", passed, totalTests, float64(passed)*100/float64(totalTests))

	// 写结果到文件
	resultPath := "gofun_correctness_result.md"
	var sb strings.Builder
	sb.WriteString("# gofun 正确性测试结果\n\n")
	sb.WriteString(fmt.Sprintf("使用 Gig 的 %d 个 correctness 测试用例（无参导出函数）跑新 gofun (onefun/gofun)。\n\n", totalTests))
	sb.WriteString(fmt.Sprintf("- **通过**: %d (%.1f%%)\n", passed, float64(passed)*100/float64(totalTests)))
	sb.WriteString(fmt.Sprintf("- **Build 失败**: %d\n", buildFails))
	sb.WriteString(fmt.Sprintf("- **Run 失败**: %d\n", runFails))
	sb.WriteString(fmt.Sprintf("- **总失败**: %d (%.1f%%)\n\n", buildFails+runFails, float64(buildFails+runFails)*100/float64(totalTests)))

	sb.WriteString("## 按包统计\n\n")
	sb.WriteString("| 包 | 总数 | 通过 | 失败 | 通过率 |\n")
	sb.WriteString("|-----|------|------|------|--------|\n")
	for _, pkg := range packages {
		pkgFail := 0
		for _, f := range failures {
			if f.pkg == pkg {
				pkgFail++
			}
		}
		srcPath := filepath.Join(testdataDir, pkg, "main.go")
		srcBytes, err := os.ReadFile(srcPath)
		if err != nil {
			continue
		}
		funcs := extractExportedFuncs(string(srcBytes))
		pkgTotal := len(funcs)
		pkgPass := pkgTotal - pkgFail
		if pkgTotal == 0 {
			continue
		}
		pct := float64(pkgPass) * 100 / float64(pkgTotal)
		marker := "✅"
		if pkgFail > 0 {
			marker = "❌"
		}
		sb.WriteString(fmt.Sprintf("| %s | %d | %d | %d | %.0f%% %s |\n", pkg, pkgTotal, pkgPass, pkgFail, pct, marker))
	}

	sb.WriteString("\n## 失败列表\n\n")
	lastPkg := ""
	for _, f := range failures {
		if f.pkg != lastPkg {
			sb.WriteString(fmt.Sprintf("\n### %s\n\n", f.pkg))
			lastPkg = f.pkg
		}
		errMsg := f.err
		if len(errMsg) > 150 {
			errMsg = errMsg[:150] + "..."
		}
		sb.WriteString(fmt.Sprintf("- `%s` [%s]: %s\n", f.funcName, f.phase, errMsg))
	}

	sb.WriteString("\n## 错误类型统计\n\n")
	for _, e := range sortedErrs {
		sb.WriteString(fmt.Sprintf("- **%d** × `%s`\n", e.count, e.err))
	}

	os.WriteFile(resultPath, []byte(sb.String()), 0644)
	t.Logf("结果已写入 %s", resultPath)
}
