package divergence_hunt56

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Round 56: Real-world data processing - ETL, aggregation, reporting
// ============================================================================

func SalesReport() float64 {
	type Sale struct {
		Product string
		Amount  float64
		Region  string
	}
	sales := []Sale{
		{"A", 100, "East"},
		{"B", 200, "West"},
		{"A", 150, "East"},
		{"C", 300, "North"},
		{"B", 250, "West"},
	}
	total := 0.0
	for _, s := range sales { total += s.Amount }
	return total
}

func SalesByRegion() int {
	type Sale struct {
		Amount float64
		Region string
	}
	sales := []Sale{
		{100, "East"},
		{200, "West"},
		{150, "East"},
	}
	byRegion := map[string]float64{}
	for _, s := range sales {
		byRegion[s.Region] += s.Amount
	}
	return int(byRegion["East"])
}

func TopProducts() string {
	type Product struct {
		Name   string
		Sales  int
	}
	products := []Product{
		{"Widget", 100},
		{"Gadget", 250},
		{"Doohickey", 175},
	}
	sort.Slice(products, func(i, j int) bool {
		return products[i].Sales > products[j].Sales
	})
	return products[0].Name
}

func DataCleaning() int {
	raw := []string{"Alice", "", "Bob", "  ", "Charlie", ""}
	cleaned := []string{}
	for _, s := range raw {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return len(cleaned)
}

func DataNormalization() int {
	data := []float64{100, 200, 300, 400, 500}
	min, max := data[0], data[0]
	for _, v := range data {
		if v < min { min = v }
		if v > max { max = v }
	}
	// Normalize to 0-1 range
	normalized := make([]float64, len(data))
	for i, v := range data {
		normalized[i] = (v - min) / (max - min)
	}
	return int(normalized[0] * 100) // should be 0
}

func JSONDataExport() string {
	type Record struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Score int    `json:"score"`
	}
	records := []Record{
		{1, "Alice", 85},
		{2, "Bob", 92},
	}
	data, _ := json.Marshal(records)
	return string(data)
}

func PivotTable() int {
	type Entry struct {
		Category string
		Region   string
		Value    int
	}
	entries := []Entry{
		{"A", "East", 10},
		{"A", "West", 20},
		{"B", "East", 30},
		{"B", "West", 40},
	}
	pivot := map[string]map[string]int{}
	for _, e := range entries {
		if pivot[e.Category] == nil {
			pivot[e.Category] = map[string]int{}
		}
		pivot[e.Category][e.Region] += e.Value
	}
	return pivot["A"]["East"] + pivot["B"]["West"]
}

func PercentileCalc() int {
	data := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	sort.Ints(data)
	median := data[len(data)/2]
	return median
}

func MovingAverage() float64 {
	data := []float64{10, 20, 30, 40, 50}
	window := 3
	sum := 0.0
	for i := 0; i < window; i++ { sum += data[i] }
	return sum / float64(window)
}

func FrequencyDistribution() string {
	data := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	freq := map[int]int{}
	for _, v := range data { freq[v]++ }
	return fmt.Sprintf("%d:%d", freq[3], freq[4])
}

func DataMerge() int {
	left := map[string]int{"a": 1, "b": 2}
	right := map[string]int{"b": 3, "c": 4}
	merged := map[string]int{}
	for k, v := range left { merged[k] = v }
	for k, v := range right { merged[k] = v }
	return merged["b"]
}

func StringReport() string {
	data := []string{"apple", "banana", "avocado", "blueberry", "cherry"}
	sort.Strings(data)
	return strings.Join(data, ",")
}
