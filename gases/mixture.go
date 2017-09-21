package gases

import (
	"fmt"
)

func NewMixture(gases []Gas, fractions []float64) Gas {
	if len(gases) != len(fractions) {
		panic(fmt.Sprintf("len(gases) == %d; len(fractions) == %d", len(gases), len(fractions)))
	}

	var fracSum float64 = 0
	for _, v := range fractions {
		fracSum += v
	}

	var normFracs = make([]float64, len(fractions))
	for i, v := range fractions {
		normFracs[i] = v / fracSum
	}

	return mixture{
		gases:     gases,
		fractions: normFracs,
	}
}

type mixture struct {
	gases     []Gas
	fractions []float64
}

func (m mixture) Cp(t float64) float64 {
	var result float64 = 0
	for i := range m.fractions {
		result += m.gases[i].Cp(t) * m.fractions[i]
	}

	return result
}

func (m mixture) R() float64 {
	var result float64 = 0
	for i := range m.fractions {
		result += m.gases[i].R() * m.fractions[i]
	}

	return result
}
