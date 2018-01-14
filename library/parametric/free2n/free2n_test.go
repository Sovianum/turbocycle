package free2n

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	c "github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

const (
	cMassRate0 = 5
	piC0       = 8
	etaC0      = 0.85
	rpmC0      = 1e5

	cpSigma = 0.98

	ctMassRate0 = 12
	piCT0       = 1.4
	etaCT0      = 0.9
	dCTMean     = 0.7

	ctpSigma = 0.98

	ftMassRate0 = 6
	piFT0       = 14
	etaFT0      = 0.9
	dFTMean     = 0.7

	ftpSigma = 0.98

	t0 = 288
	p0 = 1e5

	tFuel = 300

	etaBurn          = 0.99
	lambdaIn0        = 0.2
	tStagIn0         = 500
	fuelMassRateRel0 = 0.03

	payloadPower0 = 2.3e6
	payloadRPM0   = 3e3

	precision = 1e-3

	pAtm = 1e5
	tAtm = 288
	tGas = 1700
)

func TestNewDoubleShaftFreeScheme_Smoke(t *testing.T) {
	var scheme = getUnitTestScheme()
	var network, err = scheme.GetNetwork()
	assert.Nil(t, err)

	var converged = false
	converged, err = network.Solve(1, 2, 100, 1e-3)

	assert.True(t, converged)
	assert.Nil(t, err)

	var sysCall = variator.SysCallFromNetwork(
		network, scheme.Assembler().GetVectorPort(), 1, 2, 100, 1e-3,
	)
	var variators = scheme.Variators()
	var solverGen = newton.NewUniformNewtonSolverGen(1e-4)

	var variatorSolver = variator.NewVariatorSolver(sysCall, variators, solverGen)

	_, err = variatorSolver.Solve(
		mat.NewVecDense(6, []float64{0.8, 0.8, 0.034, 1, 1, 1}),
		1e-6, 0.1, 1000,
	)
	assert.Nil(t, err)
}

func getUnitTestScheme() DoubleShaftFreeScheme {
	return getTestScheme(
		func(lambdaU, normPiStag float64) float64 {
			return 1 - (1-normPiStag)*5
		},
		func(lambdaU, normPiStag float64) float64 {
			return 1
		},
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
			return normRPM
		},
	)
}

func getTestScheme(
	ctNormMassRateChar, ctNormEtaChar,
	ftNormMassRateChar, ftNormEtaChar c.TurbineCharFunc,
	compressorNormEtaChar, compressorNormRpmChar c.CompressorCharFunc,
	burnerSigmaFunc func(lambda float64) float64,
	powerChar func(normRPM float64) float64,
) DoubleShaftFreeScheme {
	var compressor = c.NewParametricCompressorNode(
		cMassRate0, piC0, rpmC0, etaC0, t0, p0, precision,
		compressorNormEtaChar, compressorNormRpmChar,
	)
	var compressorPipe = c.NewPressureLossNode(cpSigma)

	var burner = c.NewParametricBurnerNode(
		fuel.GetCH4(), tFuel, t0, etaBurn, lambdaIn0, p0*piC0, tStagIn0,
		cMassRate0, fuelMassRateRel0, precision, burnerSigmaFunc,
	)

	var cTurbine = c.NewSimpleParametricTurbineNode(
		ctMassRate0, piCT0, etaCT0, t0, p0, dCTMean, precision,
		0, 0, 0,
		ctNormMassRateChar, ctNormEtaChar,
	)
	var ctPipe = c.NewPressureLossNode(ctpSigma)

	var fTurbine = c.NewSimpleParametricTurbineNode(
		ftMassRate0, piFT0, etaFT0, t0, p0, dFTMean, precision,
		0, 0, 0,
		ftNormMassRateChar, ftNormEtaChar,
	)
	var ftPipe = c.NewPressureLossNode(ftpSigma)

	var payload = c.NewPayload(payloadRPM0, -payloadPower0, powerChar)

	return NewDoubleShaftFreeScheme(
		gases.GetAir(), tAtm, pAtm, tGas,
		0.99, compressor, compressorPipe, burner, cTurbine, ctPipe,
		fTurbine, ftPipe, payload,
	)
}
