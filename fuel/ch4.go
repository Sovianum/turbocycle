package fuel

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/gases"
)

func GetCH4() GasFuel {
	return ch4{}
}

type ch4 struct {}

func (gas ch4) Cp(t float64) float64 {
	var tArr = []float64{
		200, 225, 250, 275, 300,
		325, 350, 375, 400, 450,
		500, 550, 600, 650, 700,
		750, 800, 850, 900, 950,
		1000, 1050, 1100,
	}

	var cpArr = []float64{
		2087, 2121, 2156, 2191, 2226,
		2293, 2365, 2442, 2525, 2703,
		2889, 3074, 3256, 3432, 3602,
		3766, 3923, 4072, 4214, 4348,
		4475, 4595, 4708,
	}

	var cp, err = common.Interp(t, tArr, cpArr)
	if err != nil {
		panic(err)
	}

	return cp
}

func (gas ch4) AirMassTheory() float64 {
	var oxygenMassFraction = 0.2315
	return 2 / oxygenMassFraction * common.O2Weight / common.CH4Weight
}

func (gas ch4) QLower() float64 {
	return 49030e3
}

func (gas ch4) GetCombustionGas(alpha float64) gases.Gas {
	var factor = 1 / (1 + 2 * common.O2AirFraction * common.CH4Weight / common.O2Weight)

	var omegaN2 = factor * (common.N2AirFraction / alpha)
	var omegaCO2 = factor * (common.O2AirFraction / alpha * common.CO2Weight / common.O2Weight)
	var omegaH2O = factor * (2 * common.O2AirFraction / alpha * common.H2OWeight / common.O2Weight)
	var omegaAir = factor * (1 - 1 / alpha)

	return gases.NewMixture(
		[]gases.Gas{gases.GetNitrogen(), gases.GetCO2(), gases.GetH2OVapour(), gases.GetAir()},
		[]float64{omegaN2, omegaCO2, omegaH2O, omegaAir},
	)
}
