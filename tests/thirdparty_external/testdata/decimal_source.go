package main

import "github.com/shopspring/decimal"

func DecimalAdd(a, b string) string {
	da, _ := decimal.NewFromString(a)
	db, _ := decimal.NewFromString(b)
	return da.Add(db).StringFixed(2)
}

func DecimalSum() string {
	d1 := decimal.NewFromInt(1)
	d2 := decimal.NewFromInt(2)
	d3 := decimal.NewFromInt(3)
	return decimal.Sum(d1, d2, d3).String()
}

func DecimalAvg() string {
	d1 := decimal.NewFromInt(2)
	d2 := decimal.NewFromInt(4)
	d3 := decimal.NewFromInt(8)
	return decimal.Avg(d1, d2, d3).StringFixed(2)
}
