package complex_tests

func DeferBasic() int { result := 0; defer func() { result += 10 }(); result += 1; return result }

func DeferMultiple() int { result := 0; defer func() { result += 100 }(); defer func() { result += 10 }(); return result }

func DeferClosureCapture() int { result := 0; x := 10; defer func() { result += x }(); x = 20; return result }

func DeferRecover() int { defer func() { recover() }(); panic("test") }

func DeferRecoverCheck() int { DeferRecover(); return 1 }
