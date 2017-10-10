package two_shafts

import (
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/nodes/compose"
	"github.com/Sovianum/turbocycle/impl/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/nodes/source"
	"github.com/Sovianum/turbocycle/schemes"
	"github.com/Sovianum/turbocycle/application"
	"errors"
)

const (
	tAtm = 288
	pAtm = 1e5
	sigmaInlet = 0.98
	etaComp = 0.86
	piComp = 11
	tGas = 1223
	tFuel = 300
	sigmaBurn = 0.99
	etaBurn = 0.98
	initAlpha = 3
	t0 = 300
	etaCompTurbine = 0.9
	lambdaOut = 0.3
	dgCompTurbine = -0.01
	etaM = 0.99

	sigmaCompTurbinePipe = 0.98

	etaFreeTurbine = 0.92
	dgFreeTurbine = -0.01
	freeTurbinePressureLossSigma = 0.93

	precision = 0.05
)

func GetInitedTwoShaftsScheme() schemes.TwoShaftsScheme {
	var gasSource = source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm)
	var inletPressureDrop = constructive.NewPressureLossNode(sigmaInlet)
	var gasGenerator = compose.NewGasGeneratorNode(
		etaComp, piComp, fuel.GetCH4(),
		tGas, tFuel, sigmaBurn, etaBurn, initAlpha, t0,
		etaCompTurbine, lambdaOut, func(node constructive.TurbineNode) float64 {
			return dgCompTurbine
		},
		etaM, precision,
	)
	var compressorTurbinePipe = constructive.NewPressureLossNode(sigmaCompTurbinePipe)
	var freeTurbineBlock = compose.NewFreeTurbineBlock(
		pAtm,
		etaFreeTurbine, lambdaOut, precision, func(node constructive.TurbineNode) float64 {
			return dgFreeTurbine
		}, freeTurbinePressureLossSigma,
	)

	return schemes.NewTwoShaftsScheme(gasSource, inletPressureDrop, gasGenerator, compressorTurbinePipe, freeTurbineBlock)
}

func GetDataGenerator(
	scheme schemes.TwoShaftsScheme, power float64, relaxCoef float64, iterNum int,
) func(pi float64) (application.SingleCompressorDataPoint, error) {
	return func(pi float64) (application.SingleCompressorDataPoint, error) {
		scheme.Compressor().SetPiStag(pi)
		var converged, err = scheme.GetNetwork().Solve(relaxCoef, iterNum, 0.001)
		if err != nil {
			return application.SingleCompressorDataPoint{}, err
		}
		if !converged {
			return application.SingleCompressorDataPoint{}, errors.New("not converged")
		}

		return application.SingleCompressorDataPoint{
			Pi: pi,
			Efficiency:schemes.GetEfficiency(scheme),
			MassRate:schemes.GetMassRate(power, scheme),
			SpecificFuelRate:schemes.GetSpecificFuelRate(scheme),
		}, nil
	}
}
