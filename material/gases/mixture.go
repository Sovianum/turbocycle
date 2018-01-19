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

func (m mixture) String() string {
	result := "Mixture: "
	for i := range m.gases {
		result += fmt.Sprintf("%s (%.1f%)", m.gases[i].String(), m.fractions[i]*100)
	}
	return result
}

func (m mixture) OxygenMassFraction() float64 {
	return m.getParameter(func(gas Gas) float64 {
		return gas.OxygenMassFraction()
	})
}

func (m mixture) Cp(t float64) float64 {
	return m.getParameter(func(gas Gas) float64 {
		return gas.Cp(t)
	})
}

func (m mixture) Mu(t float64) float64 {
	return m.getParameter(func(gas Gas) float64 {
		return gas.Mu(t)
	})
}

func (m mixture) Lambda(t float64) float64 {
	return m.getParameter(func(gas Gas) float64 {
		return gas.Lambda(t)
	})
}

func (m mixture) R() float64 {
	return m.getParameter(func(gas Gas) float64 {
		return gas.R()
	})
}

func (m mixture) getParameter(f func(gas Gas) float64) float64 {
	var result float64 = 0
	for i := range m.fractions {
		result += f(m.gases[i]) * m.fractions[i]
	}

	return result
}
