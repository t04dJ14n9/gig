package complex

// ============================================================================
// Complex Closure Tests
// ============================================================================

// ClosureCaptureLoop tests capturing loop variables correctly.
func ClosureCaptureLoop() int {
	fns := make([]func() int, 3)
	for i := 0; i < 3; i++ {
		v := i
		fns[i] = func() int { return v }
	}
	return fns[0]() + fns[1]()*10 + fns[2]()*100
}

// ClosureMutualRecursion tests closures that call each other.
func ClosureMutualRecursion() int {
	var isEven func(int) bool
	var isOdd func(int) bool
	isEven = func(n int) bool {
		if n == 0 {
			return true
		}
		return isOdd(n - 1)
	}
	isOdd = func(n int) bool {
		if n == 0 {
			return false
		}
		return isEven(n - 1)
	}
	if isEven(10) && isOdd(7) {
		return 1
	}
	return 0
}

// ClosureCurryAdd tests currying with closures.
func ClosureCurryAdd() int {
	add := func(a int) func(int) int {
		return func(b int) int {
			return a + b
		}
	}
	return add(10)(20)
}

// ClosureCounterChain tests chained counter closures.
func ClosureCounterChain() int {
	makeCounter := func(start int) func() int {
		count := start
		return func() int {
			count++
			return count
		}
	}
	c1 := makeCounter(0)
	c2 := makeCounter(100)
	return c1() + c1() + c2() + c2()
}

// ClosureMemoize tests memoization pattern.
func ClosureMemoize() int {
	memo := make(map[int]int)
	var fib func(int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		if v, ok := memo[n]; ok {
			return v
		}
		result := fib(n-1) + fib(n-2)
		memo[n] = result
		return result
	}
	return fib(20)
}

// ClosureYieldGenerator tests generator pattern.
func ClosureYieldGenerator() int {
	generate := func(start, end int) func() (int, bool) {
		current := start
		return func() (int, bool) {
			if current > end {
				return 0, false
			}
			v := current
			current++
			return v, true
		}
	}
	gen := generate(1, 5)
	sum := 0
	for {
		v, ok := gen()
		if !ok {
			break
		}
		sum += v
	}
	return sum
}

// ClosureFilter tests filter pattern.
func ClosureFilter() int {
	makeFilter := func(predicate func(int) bool) func([]int) []int {
		return func(nums []int) []int {
			result := []int{}
			for _, n := range nums {
				if predicate(n) {
					result = append(result, n)
				}
			}
			return result
		}
	}
	isEven := func(n int) bool { return n%2 == 0 }
	filter := makeFilter(isEven)
	filtered := filter([]int{1, 2, 3, 4, 5, 6})
	return len(filtered)
}

// ClosureReduce tests reduce pattern.
func ClosureReduce() int {
	makeReducer := func(accumulate func(int, int) int, initial int) func([]int) int {
		return func(nums []int) int {
			result := initial
			for _, n := range nums {
				result = accumulate(result, n)
			}
			return result
		}
	}
	sum := func(a, b int) int { return a + b }
	reducer := makeReducer(sum, 0)
	return reducer([]int{1, 2, 3, 4, 5})
}

// ClosureCompose tests function composition.
func ClosureCompose() int {
	compose := func(f, g func(int) int) func(int) int {
		return func(x int) int {
			return f(g(x))
		}
	}
	double := func(x int) int { return x * 2 }
	increment := func(x int) int { return x + 1 }
	composed := compose(double, increment)
	return composed(5)
}

// ClosurePartial tests partial application.
func ClosurePartial() int {
	partial := func(fn func(int, int) int, a int) func(int) int {
		return func(b int) int {
			return fn(a, b)
		}
	}
	add := func(a, b int) int { return a + b }
	add5 := partial(add, 5)
	return add5(10)
}

// ClosureOnce tests single execution pattern.
func ClosureOnce() int {
	once := func(fn func() int) func() int {
		called := false
		var result int
		return func() int {
			if !called {
				result = fn()
				called = true
			}
			return result
		}
	}
	expensive := func() int { return 42 }
	onceFn := once(expensive)
	return onceFn() + onceFn() + onceFn()
}

// ClosureState tests stateful closure.
func ClosureState() int {
	makeAccumulator := func() func(int) int {
		sum := 0
		return func(n int) int {
			sum += n
			return sum
		}
	}
	acc := makeAccumulator()
	return acc(1) + acc(2) + acc(3)
}

// ClosureDefer tests closure with defer.
func ClosureDefer() int {
	result := 0
	func() {
		defer func() {
			result += 10
		}()
		result += 1
	}()
	return result
}

// ClosureCaptureSlice tests capturing slice.
func ClosureCaptureSlice() int {
	s := []int{1, 2, 3}
	push := func(v int) {
		s = append(s, v)
	}
	push(4)
	push(5)
	return s[len(s)-1]
}

// ClosureCaptureMap tests capturing map.
func ClosureCaptureMap() int {
	m := make(map[string]int)
	set := func(k string, v int) {
		m[k] = v
	}
	set("a", 1)
	set("b", 2)
	return m["a"] + m["b"]
}

// ClosureRecursive tests recursive closure.
func ClosureRecursive() int {
	var factorial func(int) int
	factorial = func(n int) int {
		if n <= 1 {
			return 1
		}
		return n * factorial(n-1)
	}
	return factorial(6)
}

// ClosureInStruct tests closure in struct.
func ClosureInStruct() int {
	type Processor struct {
		Process func(int) int
	}
	p := Processor{
		Process: func(x int) int { return x * 2 },
	}
	return p.Process(21)
}

// ClosureSliceOfFuncs tests slice of closures.
func ClosureSliceOfFuncs() int {
	fns := make([]func(int) int, 5)
	for i := 0; i < 5; i++ {
		factor := i + 1
		fns[i] = func(x int) int { return x * factor }
	}
	sum := 0
	for _, fn := range fns {
		sum += fn(10)
	}
	return sum
}

// ClosureMapOfFuncs tests map of closures.
func ClosureMapOfFuncs() int {
	ops := map[string]func(int, int) int{
		"add": func(a, b int) int { return a + b },
		"mul": func(a, b int) int { return a * b },
	}
	return ops["add"](1, 2) + ops["mul"](3, 4)
}

// ============================================================================
// Complex Control Flow Tests
// ============================================================================

func ControlNestedIf() int {
	x, y, z := 1, 2, 3
	result := 0
	if x > 0 {
		if y > 1 {
			if z > 2 {
				result = 100
			} else {
				result = 10
			}
		} else {
			result = 1
		}
	}
	return result
}

func ControlNestedLoop() int {
	sum := 0
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			sum += i * j
		}
	}
	return sum
}

func ControlTripleLoop() int {
	count := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			for k := 0; k < 5; k++ {
				if i+j+k == 6 {
					count++
				}
			}
		}
	}
	return count
}

// ============================================================================
// Complex Defer Tests
// ============================================================================

func DeferBasic() int { result := 0; defer func() { result += 10 }(); result += 1; return result }

func DeferMultiple() int {
	result := 0
	defer func() { result += 100 }()
	defer func() { result += 10 }()
	return result
}

func DeferClosureCapture() int {
	result := 0
	x := 10
	defer func() { result += x }()
	x = 20
	return result
}

func DeferRecover() int { defer func() { recover() }(); panic("test") }

func DeferRecoverCheck() int { DeferRecover(); return 1 }

// ============================================================================
// Complex Map Operations Tests
// ============================================================================

// MapNested tests nested maps.
func MapNested() int {
	m := map[string]map[string]int{
		"a": {"x": 1, "y": 2},
		"b": {"x": 3, "y": 4},
	}
	return m["a"]["x"] + m["b"]["y"]
}

// MapMerge tests merging maps.
func MapMerge() int {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}
	result := make(map[string]int)
	for k, v := range m1 {
		result[k] = v
	}
	for k, v := range m2 {
		result[k] = v
	}
	return result["a"] + result["b"] + result["c"]
}

// MapInvert tests inverting key-value.
func MapInvert() int {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	result := make(map[string]int)
	for k, v := range m {
		result[v] = k
	}
	return result["a"] + result["b"]*10 + result["c"]*100
}

// MapKeys tests getting keys.
func MapKeys() int {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	keys := []int{}
	for k := range m {
		keys = append(keys, k)
	}
	return len(keys)
}

// MapValues tests getting values.
func MapValues() int {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	values := []string{}
	for _, v := range m {
		values = append(values, v)
	}
	return len(values)
}

// MapFilterKeys tests filtering by key.
func MapFilterKeys() int {
	m := map[int]int{1: 10, 2: 20, 3: 30, 4: 40, 5: 50}
	result := make(map[int]int)
	for k, v := range m {
		if k%2 == 0 {
			result[k] = v
		}
	}
	return len(result)
}

// MapFilterValues tests filtering by value.
func MapFilterValues() int {
	m := map[int]int{1: 10, 2: 20, 3: 30, 4: 40, 5: 50}
	result := make(map[int]int)
	for k, v := range m {
		if v > 25 {
			result[k] = v
		}
	}
	return len(result)
}

// MapMapKeys tests mapping keys.
func MapMapKeys() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	result := make(map[int]int)
	for k, v := range m {
		result[k*10] = v
	}
	return result[10] + result[20] + result[30]
}

// MapMapValues tests mapping values.
func MapMapValues() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	result := make(map[int]int)
	for k, v := range m {
		result[k] = v * 2
	}
	return result[1] + result[2] + result[3]
}

// MapCounter tests using map as counter.
func MapCounter() int {
	str := "hello world"
	counts := make(map[rune]int)
	for _, r := range str {
		counts[r]++
	}
	return counts['l']
}

// MapHistogram tests histogram.
func MapHistogram() int {
	arr := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	hist := make(map[int]int)
	for _, v := range arr {
		hist[v]++
	}
	return hist[4]
}

// MapGroupBy tests grouping.
func MapGroupBy() int {
	arr := []string{"apple", "banana", "apricot", "cherry", "avocado"}
	groups := make(map[byte][]string)
	for _, s := range arr {
		key := s[0]
		groups[key] = append(groups[key], s)
	}
	return len(groups['a'])
}

// MapSet tests set operations with map.
func MapSet() int {
	set := make(map[int]bool)
	add := func(v int) { set[v] = true }
	contains := func(v int) bool { return set[v] }
	remove := func(v int) { delete(set, v) }
	add(1)
	add(2)
	add(3)
	if contains(2) {
		remove(2)
	}
	return len(set)
}

// MapMultiSet tests multiset.
func MapMultiSet() int {
	mset := make(map[int]int)
	add := func(v int) { mset[v]++ }
	remove := func(v int) {
		if mset[v] > 0 {
			mset[v]--
			if mset[v] == 0 {
				delete(mset, v)
			}
		}
	}
	add(1)
	add(1)
	add(2)
	remove(1)
	return mset[1]
}

// MapBiMap tests bidirectional map.
func MapBiMap() int {
	forward := make(map[string]int)
	backward := make(map[int]string)
	insert := func(k string, v int) {
		forward[k] = v
		backward[v] = k
	}
	insert("a", 1)
	insert("b", 2)
	return forward["a"] + len(backward[2])
}

// MapUpdateNested tests updating nested map.
func MapUpdateNested() int {
	m := map[string]map[string]int{
		"a": {"x": 1},
	}
	if inner, ok := m["a"]; ok {
		inner["y"] = 2
	}
	return m["a"]["x"] + m["a"]["y"]
}

// MapDeleteWhileIterating tests safe deletion pattern.
func MapDeleteWhileIterating() int {
	m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
	toDelete := []int{}
	for k, v := range m {
		if v%2 == 0 {
			toDelete = append(toDelete, k)
		}
	}
	for _, k := range toDelete {
		delete(m, k)
	}
	return len(m)
}

// MapComplexKey tests complex key.
func MapComplexKey() int {
	type Key struct {
		X, Y int
	}
	m := make(map[Key]int)
	m[Key{1, 2}] = 10
	m[Key{3, 4}] = 20
	return m[Key{1, 2}] + m[Key{3, 4}]
}

// MapDefaultValue tests default value pattern.
func MapDefaultValue() int {
	m := map[int]int{1: 10, 2: 20}
	get := func(k int) int {
		if v, ok := m[k]; ok {
			return v
		}
		return 0
	}
	return get(1) + get(3)
}

// MapAccumulate tests accumulation.
func MapAccumulate() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}

// MapFindKey tests finding key by value.
func MapFindKey() int {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	findKey := func(v string) int {
		for k, val := range m {
			if val == v {
				return k
			}
		}
		return -1
	}
	return findKey("b")
}

// MapFrequency tests frequency count.
func MapFrequency() int {
	str := "mississippi"
	freq := make(map[rune]int)
	for _, r := range str {
		freq[r]++
	}
	return freq['s']*10 + freq['i']
}

// MapMemoize tests memoization with map.
func MapMemoize() int {
	cache := make(map[int]int)
	var fib func(int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		if v, ok := cache[n]; ok {
			return v
		}
		result := fib(n-1) + fib(n-2)
		cache[n] = result
		return result
	}
	return fib(20)
}

// MapIncrement tests increment operations.
func MapIncrement() int {
	m := make(map[int]int)
	for i := 0; i < 10; i++ {
		m[i%3]++
	}
	return m[0] + m[1] + m[2]
}

// MapCopy tests copying map.
func MapCopy() int {
	m1 := map[int]int{1: 10, 2: 20, 3: 30}
	m2 := make(map[int]int)
	for k, v := range m1 {
		m2[k] = v
	}
	m2[1] = 100
	return m1[1] + m2[1]
}

// ============================================================================
// Complex Pointer Tests
// ============================================================================

func PointerBasic() int { x := 42; p := &x; return *p }

func PointerModify() int { x := 10; p := &x; *p = 20; return x }

func PointerSwap() int { a, b := 1, 2; pa, pb := &a, &b; *pa, *pb = *pb, *pa; return a*10 + b }

func PointerToPointer() int { x := 42; p := &x; pp := &p; return **pp }

func PointerSlice() int { arr := []int{1, 2, 3}; p := &arr; *p = append(*p, 4); return len(arr) }

// ============================================================================
// Complex Recursion Tests
// ============================================================================

// RecursionFib tests Fibonacci recursion.
func RecursionFib(n int) int {
	if n <= 1 {
		return n
	}
	return RecursionFib(n-1) + RecursionFib(n-2)
}

// RecursionFibCheck tests Fibonacci.
func RecursionFibCheck() int {
	return RecursionFib(15)
}

// RecursionAckermann tests Ackermann function.
func RecursionAckermann(m, n int) int {
	if m == 0 {
		return n + 1
	}
	if n == 0 {
		return RecursionAckermann(m-1, 1)
	}
	return RecursionAckermann(m-1, RecursionAckermann(m, n-1))
}

// RecursionAckermannCheck tests Ackermann.
func RecursionAckermannCheck() int {
	return RecursionAckermann(3, 4)
}

// RecursionGCD tests GCD via recursion.
func RecursionGCD(a, b int) int {
	if b == 0 {
		return a
	}
	return RecursionGCD(b, a%b)
}

// RecursionGCDCheck tests GCD.
func RecursionGCDCheck() int {
	return RecursionGCD(48, 18)
}

// RecursionSum tests sum via recursion.
func RecursionSum(n int) int {
	if n <= 0 {
		return 0
	}
	return n + RecursionSum(n-1)
}

// RecursionSumCheck tests sum recursion.
func RecursionSumCheck() int {
	return RecursionSum(100)
}

// RecursionPower tests power via recursion.
func RecursionPower(base, exp int) int {
	if exp == 0 {
		return 1
	}
	if exp%2 == 0 {
		half := RecursionPower(base, exp/2)
		return half * half
	}
	return base * RecursionPower(base, exp-1)
}

// RecursionPowerCheck tests power.
func RecursionPowerCheck() int {
	return RecursionPower(2, 10)
}

// RecursionBinarySearch tests binary search recursion.
func RecursionBinarySearch(arr []int, target, low, high int) int {
	if low > high {
		return -1
	}
	mid := (low + high) / 2
	if arr[mid] == target {
		return mid
	}
	if arr[mid] > target {
		return RecursionBinarySearch(arr, target, low, mid-1)
	}
	return RecursionBinarySearch(arr, target, mid+1, high)
}

// RecursionBinarySearchCheck tests binary search.
func RecursionBinarySearchCheck() int {
	arr := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
	return RecursionBinarySearch(arr, 11, 0, len(arr)-1)
}

// RecursionReverse tests string reversal via recursion.
func RecursionReverse(s string) string {
	if len(s) <= 1 {
		return s
	}
	return RecursionReverse(s[1:]) + string(s[0])
}

// RecursionReverseCheck tests reversal.
func RecursionReverseCheck() int {
	r := RecursionReverse("hello")
	if r == "olleh" {
		return 1
	}
	return 0
}

// RecursionPalindrome tests palindrome check via recursion.
func RecursionPalindrome(s string) bool {
	if len(s) <= 1 {
		return true
	}
	if s[0] != s[len(s)-1] {
		return false
	}
	return RecursionPalindrome(s[1 : len(s)-1])
}

// RecursionPalindromeCheck tests palindrome.
func RecursionPalindromeCheck() int {
	if RecursionPalindrome("racecar") && !RecursionPalindrome("hello") {
		return 1
	}
	return 0
}

// RecursionMergeSort tests merge sort recursion.
func RecursionMergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := RecursionMergeSort(arr[:mid])
	right := RecursionMergeSort(arr[mid:])
	return RecursionMerge(left, right)
}

// RecursionMerge merges two sorted arrays.
func RecursionMerge(left, right []int) []int {
	result := []int{}
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

// RecursionMergeSortCheck tests merge sort.
func RecursionMergeSortCheck() int {
	arr := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}
	sorted := RecursionMergeSort(arr)
	for i := 1; i < len(sorted); i++ {
		if sorted[i] < sorted[i-1] {
			return 0
		}
	}
	return 1
}

// RecursionQuickSort tests quick sort recursion.
func RecursionQuickSort(arr []int, low, high int) {
	if low < high {
		pivot := RecursionPartition(arr, low, high)
		RecursionQuickSort(arr, low, pivot-1)
		RecursionQuickSort(arr, pivot+1, high)
	}
}

// RecursionPartition partitions array for quick sort.
func RecursionPartition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1
	for j := low; j < high; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// RecursionQuickSortCheck tests quick sort.
func RecursionQuickSortCheck() int {
	arr := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}
	RecursionQuickSort(arr, 0, len(arr)-1)
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return 0
		}
	}
	return 1
}

// RecursionCountDigits tests digit counting via recursion.
func RecursionCountDigits(n int) int {
	if n < 10 {
		return 1
	}
	return 1 + RecursionCountDigits(n/10)
}

// RecursionCountDigitsCheck tests digit counting.
func RecursionCountDigitsCheck() int {
	return RecursionCountDigits(12345)
}

// RecursionSumDigits tests digit sum via recursion.
func RecursionSumDigits(n int) int {
	if n == 0 {
		return 0
	}
	return n%10 + RecursionSumDigits(n/10)
}

// RecursionSumDigitsCheck tests digit sum.
func RecursionSumDigitsCheck() int {
	return RecursionSumDigits(12345)
}

// RecursionHanoi tests Tower of Hanoi move count.
func RecursionHanoi(n int) int {
	if n == 1 {
		return 1
	}
	return 2*RecursionHanoi(n-1) + 1
}

// RecursionHanoiCheck tests Hanoi.
func RecursionHanoiCheck() int {
	return RecursionHanoi(10)
}

// RecursionStaircase tests ways to climb staircase.
func RecursionStaircase(n int) int {
	if n <= 1 {
		return 1
	}
	return RecursionStaircase(n-1) + RecursionStaircase(n-2)
}

// RecursionStaircaseCheck tests staircase.
func RecursionStaircaseCheck() int {
	return RecursionStaircase(10)
}

// RecursionSubsetSum tests subset sum recursion.
func RecursionSubsetSum(arr []int, target, index int) bool {
	if target == 0 {
		return true
	}
	if index >= len(arr) || target < 0 {
		return false
	}
	return RecursionSubsetSum(arr, target-arr[index], index+1) ||
		RecursionSubsetSum(arr, target, index+1)
}

// RecursionSubsetSumCheck tests subset sum.
func RecursionSubsetSumCheck() int {
	arr := []int{3, 34, 4, 12, 5, 2}
	if RecursionSubsetSum(arr, 9, 0) {
		return 1
	}
	return 0
}

// RecursionPermuteCount counts permutations.
func RecursionPermuteCount(n int) int {
	if n <= 1 {
		return 1
	}
	return n * RecursionPermuteCount(n-1)
}

// RecursionPermuteCheck tests permutation count.
func RecursionPermuteCheck() int {
	return RecursionPermuteCount(7)
}

// RecursionCombination counts combinations.
func RecursionCombination(n, r int) int {
	if r == 0 || r == n {
		return 1
	}
	return RecursionCombination(n-1, r-1) + RecursionCombination(n-1, r)
}

// RecursionCombinationCheck tests combinations.
func RecursionCombinationCheck() int {
	return RecursionCombination(10, 4)
}

// ============================================================================
// Complex Slice Operations Tests
// ============================================================================

// SliceFlatten tests flattening nested slices.
func SliceFlatten() int {
	nested := [][]int{{1, 2}, {3, 4, 5}, {6}}
	flat := []int{}
	for _, inner := range nested {
		for _, v := range inner {
			flat = append(flat, v)
		}
	}
	sum := 0
	for _, v := range flat {
		sum += v
	}
	return sum
}

// SliceChunk tests chunking slice.
func SliceChunk() int {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	chunks := [][]int{}
	size := 3
	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}
		chunks = append(chunks, arr[i:end])
	}
	return len(chunks)
}

// SliceRotateLeft tests left rotation.
func SliceRotateLeft() int {
	arr := []int{1, 2, 3, 4, 5}
	k := 2
	k = k % len(arr)
	result := append(arr[k:], arr[:k]...)
	return result[0]*1000 + result[1]*100 + result[2]*10 + result[3]
}

// SliceRotateRight tests right rotation.
func SliceRotateRight() int {
	arr := []int{1, 2, 3, 4, 5}
	k := 2
	k = k % len(arr)
	result := append(arr[len(arr)-k:], arr[:len(arr)-k]...)
	return result[0]*1000 + result[1]*100 + result[2]*10 + result[3]
}

// SliceUnique tests unique elements.
func SliceUnique() int {
	arr := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	seen := make(map[int]bool)
	result := []int{}
	for _, v := range arr {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return len(result)
}

// SliceIntersect tests intersection.
func SliceIntersect() int {
	a := []int{1, 2, 3, 4, 5}
	b := []int{3, 4, 5, 6, 7}
	setA := make(map[int]bool)
	for _, v := range a {
		setA[v] = true
	}
	result := []int{}
	for _, v := range b {
		if setA[v] {
			result = append(result, v)
		}
	}
	return len(result)
}

// SliceUnion tests union.
func SliceUnion() int {
	a := []int{1, 2, 3}
	b := []int{3, 4, 5}
	set := make(map[int]bool)
	for _, v := range a {
		set[v] = true
	}
	for _, v := range b {
		set[v] = true
	}
	return len(set)
}

// SliceDifference tests difference.
func SliceDifference() int {
	a := []int{1, 2, 3, 4, 5}
	b := []int{3, 4, 5, 6, 7}
	setB := make(map[int]bool)
	for _, v := range b {
		setB[v] = true
	}
	result := []int{}
	for _, v := range a {
		if !setB[v] {
			result = append(result, v)
		}
	}
	return len(result)
}

// SliceZip tests zipping slices.
func SliceZip() int {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	result := []int{}
	for i := 0; i < len(a) && i < len(b); i++ {
		result = append(result, a[i], b[i])
	}
	return result[0] + result[1] + result[2] + result[3]
}

// SlicePartition tests partitioning by predicate.
func SlicePartition() int {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evens := []int{}
	odds := []int{}
	for _, v := range arr {
		if v%2 == 0 {
			evens = append(evens, v)
		} else {
			odds = append(odds, v)
		}
	}
	return len(evens)*10 + len(odds)
}

// SliceTake tests taking first n elements.
func SliceTake() int {
	arr := []int{1, 2, 3, 4, 5}
	n := 3
	if n > len(arr) {
		n = len(arr)
	}
	result := arr[:n]
	return result[0] + result[1] + result[2]
}

// SliceDrop tests dropping first n elements.
func SliceDrop() int {
	arr := []int{1, 2, 3, 4, 5}
	n := 2
	if n > len(arr) {
		n = len(arr)
	}
	result := arr[n:]
	return result[0] + result[1] + result[2]
}

// SliceTakeWhile tests taking while predicate is true.
func SliceTakeWhile() int {
	arr := []int{2, 4, 6, 8, 1, 3, 5}
	result := []int{}
	for _, v := range arr {
		if v%2 != 0 {
			break
		}
		result = append(result, v)
	}
	return len(result)
}

// SliceDropWhile tests dropping while predicate is true.
func SliceDropWhile() int {
	arr := []int{2, 4, 6, 8, 1, 3, 5}
	result := []int{}
	dropping := true
	for _, v := range arr {
		if dropping && v%2 != 0 {
			dropping = false
		}
		if !dropping {
			result = append(result, v)
		}
	}
	return len(result)
}

// SliceGroupBy tests grouping by key.
func SliceGroupBy() int {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	groups := make(map[int][]int)
	for _, v := range arr {
		key := v % 3
		groups[key] = append(groups[key], v)
	}
	return len(groups[0]) + len(groups[1])*10 + len(groups[2])*100
}

// SliceWindow tests sliding window.
func SliceWindow() int {
	arr := []int{1, 2, 3, 4, 5}
	size := 3
	windows := [][]int{}
	for i := 0; i <= len(arr)-size; i++ {
		windows = append(windows, arr[i:i+size])
	}
	return windows[0][0] + windows[1][1] + windows[2][2]
}

// SlicePairwise tests pairwise iteration.
func SlicePairwise() int {
	arr := []int{1, 2, 3, 4, 5}
	sum := 0
	for i := 0; i < len(arr)-1; i++ {
		sum += arr[i] + arr[i+1]
	}
	return sum
}

// SliceCartesian tests cartesian product.
func SliceCartesian() int {
	a := []int{1, 2}
	b := []int{3, 4}
	pairs := [][]int{}
	for _, x := range a {
		for _, y := range b {
			pairs = append(pairs, []int{x, y})
		}
	}
	return len(pairs)
}

// SliceTranspose tests transpose.
func SliceTranspose() int {
	matrix := [][]int{{1, 2, 3}, {4, 5, 6}}
	result := [][]int{}
	for j := 0; j < len(matrix[0]); j++ {
		row := []int{}
		for i := 0; i < len(matrix); i++ {
			row = append(row, matrix[i][j])
		}
		result = append(result, row)
	}
	return result[0][0] + result[0][1] + result[1][0] + result[1][1]
}

// SliceReverseInPlace tests in-place reversal.
func SliceReverseInPlace() int {
	arr := []int{1, 2, 3, 4, 5}
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr[0]*10000 + arr[1]*1000 + arr[2]*100 + arr[3]*10 + arr[4]
}

// SliceCompact tests removing consecutive duplicates.
func SliceCompact() int {
	arr := []int{1, 1, 2, 2, 2, 3, 3, 2, 2, 1}
	result := []int{}
	for i, v := range arr {
		if i == 0 || v != arr[i-1] {
			result = append(result, v)
		}
	}
	return len(result)
}

// SliceProduct tests computing product.
func SliceProduct() int {
	arr := []int{1, 2, 3, 4, 5}
	product := 1
	for _, v := range arr {
		product *= v
	}
	return product
}

// SliceCumSum tests cumulative sum.
func SliceCumSum() int {
	arr := []int{1, 2, 3, 4, 5}
	result := []int{}
	sum := 0
	for _, v := range arr {
		sum += v
		result = append(result, sum)
	}
	return result[len(result)-1]
}

// SliceCumProd tests cumulative product.
func SliceCumProd() int {
	arr := []int{1, 2, 3, 4}
	result := []int{}
	product := 1
	for _, v := range arr {
		product *= v
		result = append(result, product)
	}
	return result[len(result)-1]
}

// SliceFind tests finding element.
func SliceFind() int {
	arr := []int{1, 2, 3, 4, 5}
	target := 3
	for i, v := range arr {
		if v == target {
			return i
		}
	}
	return -1
}

// SliceFindLast tests finding last element.
func SliceFindLast() int {
	arr := []int{1, 2, 3, 2, 1}
	target := 2
	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] == target {
			return i
		}
	}
	return -1
}

// ============================================================================
// Complex Variadic Tests
// ============================================================================

func VariadicBasic(nums ...int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

func VariadicBasicCheck() int {
	return VariadicBasic(1, 2, 3, 4, 5)
}

func VariadicEmpty() int {
	return VariadicBasic()
}

func VariadicWithRegular(base int, nums ...int) int {
	sum := base
	for _, n := range nums {
		sum += n
	}
	return sum
}

func VariadicWithRegularCheck() int {
	return VariadicWithRegular(100, 1, 2, 3)
}

func VariadicSpread() int {
	nums := []int{1, 2, 3, 4, 5}
	return VariadicBasic(nums...)
}

func VariadicOneArg() int {
	return VariadicBasic(42)
}
