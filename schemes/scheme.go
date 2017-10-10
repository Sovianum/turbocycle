package schemes

import "github.com/Sovianum/turbocycle/core"

type Scheme interface {
	GetNetwork() core.Network
	GetSpecificPower() float64
	GetFuelMassRateRel() float64
	GetQLower() float64
}

func GetMassRate(power float64, scheme Scheme) float64 {
	return power / scheme.GetSpecificPower()
}

func GetSpecificFuelRate(scheme Scheme) float64 {
	return scheme.GetFuelMassRateRel() * 3600 / scheme.GetSpecificPower()
}

func GetEfficiency(scheme Scheme) float64 {
	return 3600 / (scheme.GetSpecificPower() * scheme.GetQLower())
}
