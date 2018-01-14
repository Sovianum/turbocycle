package fuel

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/material/gases"
)

func NewHydroCarbon(c, h int) HydroCarbon {
	return HydroCarbon{
		C: float64(c),
		H: float64(h),
	}
}

type HydroCarbon struct {
	H float64
	C float64
}

func (hc HydroCarbon) GetCombustionGas(gas gases.Gas, alpha float64) gases.Gas {
	if gas.OxygenMassFraction() < 1e-6 {
		panic("gas can not be burnt")
	}

	var exhaustComplex = hc.getExhaustComplex(gas, alpha)

	var h2oFraction = hc.getH2OComplex(gas, alpha) / exhaustComplex
	var co2Fraction = hc.getCO2Complex(gas, alpha) / exhaustComplex
	var o2Fraction = hc.getO2Complex(gas, alpha) / exhaustComplex
	var restFraction = 1 - h2oFraction - co2Fraction - o2Fraction

	return gases.NewMixture(
		[]gases.Gas{
			gases.GetH2OVapour(), gases.GetCO2(), gases.GetOxygen(),
			gases.GetOxyFreeGas(gas),
		},
		[]float64{
			h2oFraction, co2Fraction, o2Fraction, restFraction,
		},
	)
}

func (hc HydroCarbon) getO2Complex(gas gases.Gas, alpha float64) float64 {
	return (1 - hc.getAlphaFunc(alpha)) * gas.OxygenMassFraction()
}

func (hc HydroCarbon) getH2OComplex(gas gases.Gas, alpha float64) float64 {
	var factors = []float64{
		2 * hc.H / (4*hc.C + hc.H),
		common.H2OWeight / common.O2Weight,
		gas.OxygenMassFraction(),
		hc.getAlphaFunc(alpha),
	}
	return common.Product(factors)
}

func (hc HydroCarbon) getCO2Complex(gas gases.Gas, alpha float64) float64 {
	var factors = []float64{
		4 * hc.C / (4*hc.C + hc.H),
		common.CO2Weight / common.O2Weight,
		gas.OxygenMassFraction(),
		hc.getAlphaFunc(alpha),
	}
	return common.Product(factors)
}

func (hc HydroCarbon) getExhaustComplex(gas gases.Gas, alpha float64) float64 {
	var factors = []float64{
		1 / alpha,
		4 / (4*hc.C + hc.H),
		(common.CWeight*hc.C + common.HWeight*hc.H) / common.O2Weight,
		gas.OxygenMassFraction(),
	}
	return 1 + common.Product(factors)
}

func (hc HydroCarbon) getAlphaFunc(alpha float64) float64 {
	if alpha <= 1 {
		return 1
	}
	return 1 / alpha
}
