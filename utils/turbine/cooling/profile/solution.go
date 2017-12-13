package profile

import "fmt"

func NewTemperatureSolution(length int) TemperatureSolution {
	return TemperatureSolution{
		ParametricCoord:  make([]float64, length),
		X:                make([]float64, length),
		Y:                make([]float64, length),
		LengthCoord:      make([]float64, length),
		AlphaAir:         make([]float64, length),
		AlphaGas:         make([]float64, length),
		FilmTemperature:  make([]float64, length),
		AirTemperature:   make([]float64, length),
		WallTemperature:  make([]float64, length),
		HeatTransferCoef: make([]float64, length),
		FilmEfficiency:   make([]float64, length),
	}
}

type SlitSolution struct {
	Id                int
	SlitInfo          SlitInfo
	Velocity          float64
	BlowingParameter  float64
	TemperatureFactor float64
	MassRate          float64
	MassRateRel       float64
}

type TemperatureSolution struct {
	ParametricCoord  []float64
	X                []float64
	Y                []float64
	LengthCoord      []float64
	AlphaAir         []float64
	AlphaGas         []float64
	FilmTemperature  []float64
	AirTemperature   []float64
	WallTemperature  []float64
	HeatTransferCoef []float64
	FilmEfficiency   []float64

	SlitsSolution []SlitSolution
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

func (s TemperatureSolution) ToMatrix() [][]float64 {
	var result = make([][]float64, len(s.ParametricCoord))
	for i := 0; i != len(s.ParametricCoord); i++ {
		result[i] = []float64{
			s.ParametricCoord[i],
			s.X[i],
			s.Y[i],
			s.LengthCoord[i],
			s.AlphaAir[i],
			s.AlphaGas[i],
			s.FilmTemperature[i],
			s.AirTemperature[i],
			s.WallTemperature[i],
			s.HeatTransferCoef[i],
			s.FilmEfficiency[i],
		}
	}
	return result
}
