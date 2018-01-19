package gases

import (
	"fmt"

	"github.com/Sovianum/turbocycle/common"
)

type Oxidizer interface {
	OxygenMassFraction() float64
}

type noOxygen struct{}

func (noOxygen) OxygenMassFraction() float64 {
	return 0
}

func GetOxyFreeGas(gas Gas) Gas {
	var oxygenFraction = gas.OxygenMassFraction()
	return NewMixture(
		[]Gas{gas, GetOxygen()},
		[]float64{1, -oxygenFraction},
	)
}

type Gas interface {
	fmt.Stringer
	Oxidizer
	Cp(t float64) float64
	R() float64
	Mu(t float64) float64
	Lambda(t float64) float64
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
