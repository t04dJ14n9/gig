package main

import "github.com/VividCortex/ewma"

func EwmaNewMovingAverage() float64 {
	a := ewma.NewMovingAverage()
	a.Add(10.0)
	a.Add(20.0)
	a.Add(30.0)
	return a.Value()
}

func EwmaNewMovingAverageWithAge() float64 {
	a := ewma.NewMovingAverage(5.0)
	a.Add(10.0)
	a.Add(20.0)
	return a.Value()
}

func EwmaSimpleEWMA() float64 {
	var e ewma.SimpleEWMA
	e.Add(10.0)
	e.Add(20.0)
	return e.Value()
}

func EwmaSimpleEWMASet() float64 {
	var e ewma.SimpleEWMA
	e.Set(42.0)
	return e.Value()
}

func EwmaMovingAverageInterface() float64 {
	var a ewma.MovingAverage = ewma.NewMovingAverage()
	a.Add(10.0)
	return a.Value()
}
