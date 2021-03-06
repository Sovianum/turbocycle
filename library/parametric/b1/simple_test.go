package b1

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	c "github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

const (
	cMassRate0       = 3
	tMassRate0       = 7
	piTC0            = 8
	piC0             = 8
	etaTC0           = 0.9
	etaC0            = 0.85
	rpm0             = 1e5
	t0               = 288
	p0               = 1e5
	tFuel            = 300
	dMean            = 0.7
	etaBurn          = 0.99
	lambdaIn0        = 0.2
	tStagIn0         = 500
	fuelMassRateRel0 = 0.03
	power0           = 1e6
	precision        = 1e-3

	pAtm = 1e5
	tAtm = 288
	tGas = 1700
)

func TestNewSingleShaftScheme_Smoke(t *testing.T) {
	var scheme = getUnitTestScheme()
	var network, err = scheme.GetNetwork()
	assert.Nil(t, err)

	err = network.Solve(1, 2, 100, 1e-3)
	assert.Nil(t, err)

	var sysCall = variator.SysCallFromNetwork(
		network, scheme.Assembler().GetVectorPort(), 1, 2, 100, 1e-3,
	)
	var variators = scheme.Variators()
	var solverGen = newton.NewUniformNewtonSolverGen(1e-4, newton.NoLog)

	var variatorSolver = variator.NewVariatorSolver(sysCall, variators, solverGen)

	_, err = variatorSolver.Solve(
		mat.NewVecDense(5, []float64{0.8, 0.8, 0.05, 1, 1}),
		1e-8, 1, 100,
	)
	assert.Nil(t, err)

	var compressorMassRate = scheme.Compressor().MassRateOutput().GetState().Value().(float64)
	var burnerMassRate = scheme.Burner().MassRateOutput().GetState().Value().(float64)
	var turbineMassRate = scheme.Turbine().MassRateInput().GetState().Value().(float64)

	assert.InDelta(t, 0, turbineMassRate-burnerMassRate, 1e-7)

	var compressorPower = scheme.Compressor().PowerOutput().GetState().Value().(float64) * compressorMassRate
	var turbinePower = scheme.Turbine().PowerOutput().GetState().Value().(float64) * turbineMassRate
	var payloadPower = scheme.Payload().PowerOutput().GetState().Value().(float64)

	assert.InDelta(t, 0, turbinePower+compressorPower+payloadPower, 1e-7)

	var turbinePressure = scheme.Turbine().PressureOutput().GetState().Value().(float64)

	assert.InDelta(t, 0, turbinePressure-pAtm, 1e-7)
}

func getUnitTestScheme() SingleShaftScheme {
	return getTestScheme(
		func(lambdaU, normPiStag float64) float64 {
			return 1 - (1-normPiStag)*5
		},
		func(lambdaU, normPiStag float64) float64 {
			return 1
		},
		func(normMassRate, normPiStag float64) float64 {
			return 1
		},
		func(normMassRate, normPiStag float64) float64 {
			return 1
		},
		func(lambda float64) float64 {
			return 1
		},
		func(normRPM float64) float64 {
			return 1
		},
	)
}

func getTestScheme(
	turbineNormMassRateChar, turbineNormEtaChar c.TurbineCharFunc,
	compressorNormEtaChar, compressorNormRpmChar c.CompressorCharFunc,
	burnerSigmaFunc func(lambda float64) float64,
	powerChar func(normRPM float64) float64,
) SingleShaftScheme {
	var compressor = c.NewParametricCompressorNode(
		cMassRate0, piC0, rpm0, etaC0, t0, p0, precision,
		compressorNormEtaChar, compressorNormRpmChar,
	)
	var burner = c.NewParametricBurnerNode(
		fuel.GetCH4(), tFuel, t0, etaBurn, lambdaIn0, p0*piC0, tStagIn0,
		cMassRate0, fuelMassRateRel0, precision, 1, nodes.DefaultN, burnerSigmaFunc,
	)
	var turbine = c.NewSimpleParametricTurbineNode(
		tMassRate0, piTC0, etaTC0, t0, p0, dMean, precision,
		0, 0, 0,
		turbineNormMassRateChar, turbineNormEtaChar,
	)
	var payload = c.NewPayload(rpm0, -power0, powerChar)

	return NewSingleShaftScheme(
		gases.GetAir(), pAtm, tAtm, tGas,
		compressor, burner, turbine, payload,
	)
}
