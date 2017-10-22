package fuel

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/helpers/gases"
)

type GasFuel interface {
	Cp(t float64) float64
	AirMassTheory() float64
	QLower() float64
	GetCombustionGas(alpha float64) gases.Gas
}

func CpMean(fuel GasFuel, t1, t2 float64, n int) float64 {
	return common.Average(fuel.Cp, t1, t2, n)
}
