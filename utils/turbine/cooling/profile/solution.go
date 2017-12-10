package profile

import "fmt"

type TemperatureSolution struct {
	ParametricCoord  []float64
	X                []float64
	Y                []float64
	LengthCoord      []float64
	AlphaAir         []float64
	AlphaGas         []float64
	AirTemperature   []float64
	WallTemperature  []float64
	HeatTransferCoef []float64
}

func (s TemperatureSolution) String() string {
	var result = ""
	for i := range s.ParametricCoord {
		result += fmt.Sprintf(
			"t: %.4f \t l: %.4f \t alpha_gas: %.1f \t alpha_air: %.1f \t T_air: %.3f \t T_wall: %.3f \n",
			s.ParametricCoord[i],
			s.LengthCoord[i],
			s.AlphaGas[i],
			s.AlphaAir[i],
			s.AirTemperature[i],
			s.WallTemperature[i],
		)
	}
	return result
}
