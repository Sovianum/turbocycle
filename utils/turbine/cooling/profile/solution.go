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
	ParametricCoord       []float64 `json:"parametric_coord"`
	X                     []float64 `json:"x"`
	Y                     []float64 `json:"y"`
	LengthCoord           []float64 `json:"l"`
	AlphaAir              []float64 `json:"alpha_air"`
	AlphaGas              []float64 `json:"alpha_gas"`
	FilmTemperature       []float64 `json:"t_film"`
	AirTemperature        []float64 `json:"t_air"`
	WallTemperature       []float64 `json:"t_wall"`
	SmoothWallTemperature []float64 `json:"t_wall_smooth"`
	HeatTransferCoef      []float64 `json:"k_transfer"`
	FilmEfficiency        []float64 `json:"film_efficiency"`

	SlitsSolution []SlitSolution `json:"-"`
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
