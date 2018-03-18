package gases

type TestGas struct {
	CpVal     float64
	LambdaVal float64
	MuVal     float64
	RVal      float64
	OFraction float64
}

func (gas TestGas) String() string {
	return "test_gas"
}

func (gas TestGas) OxygenMassFraction() float64 {
	return gas.OFraction
}

func (gas TestGas) Cp(t float64) float64 {
	return gas.CpVal
}

func (gas TestGas) R() float64 {
	return gas.RVal
}

func (gas TestGas) Mu(t float64) float64 {
	return gas.MuVal
}

func (gas TestGas) Lambda(t float64) float64 {
	return gas.LambdaVal
}
