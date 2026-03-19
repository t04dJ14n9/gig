package initialize

// Complex initialization test with cache and multiple init functions

// Cache simulates a complex data structure initialized at package load time
type Cache struct {
	data        map[string]int
	order       []string
	sum         int
	initialized bool
}

var (
	// Global cache initialized via init()
	globalCache *Cache

	// Multiple interdependent variables
	a int
	b int
	c int

	// Complex data structures
	lookupTable map[int]string
	fibonacci   []int
)

// First init - sets up base values
func init() {
	a = 10
	b = 20
}

// Second init - uses values from first init, initializes cache
func init() {
	c = a + b // c = 30

	// Initialize complex cache
	globalCache = &Cache{
		data:        make(map[string]int),
		order:       make([]string, 0),
		sum:         0,
		initialized: true,
	}

	// Populate cache with computed values
	for i := 1; i <= 5; i++ {
		key := string(rune('A' + i - 1))
		value := a*i + b
		globalCache.data[key] = value
		globalCache.order = append(globalCache.order, key)
		globalCache.sum += value
	}
}

// Third init - depends on cache being initialized
func init() {
	// Initialize lookup table using cache values
	lookupTable = make(map[int]string)
	for k, v := range globalCache.data {
		lookupTable[v] = k
	}

	// Precompute fibonacci sequence up to sum
	fibonacci = computeFibonacci(globalCache.sum)
}

// computeFibonacci generates fibonacci sequence up to max
func computeFibonacci(max int) []int {
	if max <= 0 {
		return []int{}
	}

	result := []int{0, 1}
	for {
		next := result[len(result)-1] + result[len(result)-2]
		if next > max {
			break
		}
		result = append(result, next)
	}
	return result
}

// GetA returns the value of a
func GetA() int {
	return a
}

// GetB returns the value of b
func GetB() int {
	return b
}

// GetC returns the computed value of c
func GetC() int {
	return c
}

// GetCacheSum returns the sum of all cached values
func GetCacheSum() int {
	if globalCache == nil {
		return -1
	}
	return globalCache.sum
}

// GetCacheValue returns a specific value from cache
func GetCacheValue(key string) int {
	if globalCache == nil {
		return -1
	}
	if v, ok := globalCache.data[key]; ok {
		return v
	}
	return -1
}

// GetCacheSize returns the number of entries in cache
func GetCacheSize() int {
	if globalCache == nil {
		return -1
	}
	return len(globalCache.data)
}

// GetCacheOrder returns the insertion order count
func GetCacheOrder() int {
	if globalCache == nil {
		return -1
	}
	return len(globalCache.order)
}

// LookupByValue performs reverse lookup in the lookup table
func LookupByValue(v int) string {
	if s, ok := lookupTable[v]; ok {
		return s
	}
	return ""
}

// GetFibonacciCount returns the length of precomputed fibonacci sequence
func GetFibonacciCount() int {
	return len(fibonacci)
}

// GetFibonacciSum returns the sum of fibonacci sequence
func GetFibonacciSum() int {
	sum := 0
	for _, v := range fibonacci {
		sum += v
	}
	return sum
}

// GetFibonacciAt returns the fibonacci number at index
func GetFibonacciAt(index int) int {
	if index < 0 || index >= len(fibonacci) {
		return -1
	}
	return fibonacci[index]
}

// ComplexInitTest runs all initialization tests
func ComplexInitTest() int {
	// Verify all init() functions ran in correct order
	if a != 10 || b != 20 || c != 30 {
		return 1
	}

	// Verify cache was initialized
	if globalCache == nil || !globalCache.initialized {
		return 2
	}

	// Verify cache contents
	// A = 10*1 + 20 = 30
	// B = 10*2 + 20 = 40
	// C = 10*3 + 20 = 50
	// D = 10*4 + 20 = 60
	// E = 10*5 + 20 = 70
	expectedSum := 30 + 40 + 50 + 60 + 70 // 250
	if globalCache.sum != expectedSum {
		return 3
	}

	// Verify lookup table
	if lookupTable == nil || len(lookupTable) != 5 {
		return 4
	}
	if lookupTable[30] != "A" || lookupTable[70] != "E" {
		return 5
	}

	// Verify fibonacci sequence was computed
	// fib(250) = [0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233]
	if len(fibonacci) < 2 || fibonacci[0] != 0 || fibonacci[1] != 1 {
		return 6
	}

	// Verify fibonacci contains expected values
	expectedFib := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233}
	if len(fibonacci) != len(expectedFib) {
		return 7
	}
	for i, v := range expectedFib {
		if fibonacci[i] != v {
			return 8
		}
	}

	return 0 // All tests passed
}

// InitOrderTest verifies init() execution order
func InitOrderTest() int {
	// a should be 10 (set in first init)
	// b should be 20 (set in first init)
	// c should be 30 (a + b, set in second init)

	if GetA() != 10 {
		return 1
	}
	if GetB() != 20 {
		return 2
	}
	if GetC() != 30 {
		return 3
	}

	return 0
}

// CacheInitTest tests cache initialization
func CacheInitTest() int {
	// Test cache values
	if GetCacheValue("A") != 30 {
		return 1
	}
	if GetCacheValue("B") != 40 {
		return 2
	}
	if GetCacheValue("C") != 50 {
		return 3
	}
	if GetCacheValue("D") != 60 {
		return 4
	}
	if GetCacheValue("E") != 70 {
		return 5
	}

	// Test cache metadata
	if GetCacheSize() != 5 {
		return 6
	}
	if GetCacheOrder() != 5 {
		return 7
	}
	if GetCacheSum() != 250 {
		return 8
	}

	return 0
}

// LookupTableInitTest tests lookup table initialization
func LookupTableInitTest() int {
	// Reverse lookup should work
	if LookupByValue(30) != "A" {
		return 1
	}
	if LookupByValue(50) != "C" {
		return 2
	}
	if LookupByValue(70) != "E" {
		return 3
	}

	// Non-existent value should return empty string
	if LookupByValue(999) != "" {
		return 4
	}

	return 0
}

// FibonacciInitTest tests fibonacci sequence initialization
func FibonacciInitTest() int {
	// Should have precomputed sequence
	if GetFibonacciCount() != 14 {
		return 1
	}

	// Check specific values
	if GetFibonacciAt(0) != 0 {
		return 2
	}
	if GetFibonacciAt(1) != 1 {
		return 3
	}
	if GetFibonacciAt(7) != 13 {
		return 4
	}
	if GetFibonacciAt(13) != 233 {
		return 5
	}

	// Out of bounds should return -1
	if GetFibonacciAt(-1) != -1 {
		return 6
	}
	if GetFibonacciAt(100) != -1 {
		return 7
	}

	// Sum should be correct
	expectedSum := 0 + 1 + 1 + 2 + 3 + 5 + 8 + 13 + 21 + 34 + 55 + 89 + 144 + 233
	if GetFibonacciSum() != expectedSum {
		return 8
	}

	return 0
}

var (
	counter = 0
)

func IncCounter1() int {
	counter++
	return counter
}

func IncCounter2() int {
	counter++
	return counter
}
