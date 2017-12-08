package profile

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
