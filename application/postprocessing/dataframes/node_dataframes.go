package dataframes

type GasDF struct {
	T       float64
	P       float64
	Density float64
	K       float64
	Cp      float64
	R       float64
}

type GasMeanDF struct {
	T1     float64
	T2     float64
	P      float64
	KMean  float64
	CpMean float64
	R      float64
}

type CompressorDF struct {
	PIn  float64
	POut float64

	TIn  float64
	TOut float64

	Pi  float64
	Eta float64

	Labour float64

	GasData GasMeanDF
}

type PressureDropDF struct {
	PIn  float64
	POut float64

	TIn  float64
	TOut float64

	Sigma float64
}

type RegeneratorDF struct {
	PColdIn  float64
	PColdOut float64
	PHotIn   float64
	PHotOut  float64

	TColdIn  float64
	TColdOut float64
	THotIn   float64
	THotOut  float64

	Sigma float64
}

type FuelDF struct {
	C      float64
	TInit  float64
	T0     float64
	QLower float64
	L0     float64
}

type BurnerDF struct {
	Tg              float64
	Eta             float64
	Alpha           float64
	FuelMassRateRel float64
	Sigma           float64

	Fuel FuelDF

	AirDataInlet  GasDF
	AirData0      GasDF
	GasData0      GasDF
	GasDataOutlet GasDF
}

type BlockedTurbineDF struct {
	PIn  float64
	POut float64

	TIn  float64
	TOut float64

	GasData GasMeanDF

	MassRateRel     float64
	LeakMassRateRel float64
	CoolMassRateRel float64

	Labour float64
	Eta    float64
}

type FreeTurbineDF struct {
}

type ShaftDF struct {
	Eta float64
}
