package gases

import "github.com/Sovianum/turbocycle/common"

type Gas interface {
	Cp(t float64) float64
	R() float64
}

func Density(gas Gas, t float64, p float64) float64 {
	return p / (gas.R() * t)
}

func CpMean(gas Gas, t1 float64, t2 float64, n int) float64 {
	return common.Average(gas.Cp, t1, t2, n)
}

func K(gas Gas, t float64) float64 {
	return gas.Cp(t) / (gas.Cp(t) - gas.R())
}

func KMean(gas Gas, t1 float64, t2 float64, n int) float64 {
	var cpMean = CpMean(gas, t1, t2, n)
	return cpMean / (cpMean - gas.R())
}
