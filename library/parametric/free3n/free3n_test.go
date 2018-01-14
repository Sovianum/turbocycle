package free3n

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
	mpcMassRate0 = 50
	mpcPi0       = 8
	mpcEta0      = 0.85
	mpcRPM0      = 1e5
	mpcSigma     = 0.98

	hpcMassRate0 = 50
	hpcPi0       = 8
	hpcEta0      = 0.85
	hpcRPM0      = 1e5
	hpcSigma     = 0.98

	hptMassRate0 = 30
	hptPi0       = 5
	hptEta0      = 0.9
	hptDMean     = 0.7
	hptSigma     = 0.98

	mptMassRate0 = 20
	mptPi0       = 5
	mptEta0      = 0.9
	mptDMean     = 0.9
	mptSigma     = 0.98

	lptMassRate0 = 12
	lptPi0       = 2
	lptEta0      = 0.9
	lptDMean     = 1.1
	lptSigma     = 0.98

	mpEtaM = 0.99
	hpEtaM = 0.99

	t0 = 288
	p0 = 1e5

	tFuel = 300

	etaBurn          = 0.99
	lambdaIn0        = 0.2
	tStagIn0         = 500
	fuelMassRateRel0 = 0.03

	payloadPower0 = 3e6
	payloadRPM0   = 3e3

	precision = 1e-3

	pAtm = 1e5
	tAtm = 288
	tGas = 1200
)

func TestNewTripleShaftFreeScheme_Smoke(t *testing.T) {
	var scheme = getUnit3nFreeTestScheme()
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
		mat.NewVecDense(9, []float64{
			1, 1, 1, 1,
			1, 1, 1,
			0.034, 1,
		}),
		1e-8, 0.5, 100,
	)
	assert.Nil(t, err)

	//fmt.Println()
	//fmt.Println(scheme.LPT().PiTStag())
	//fmt.Println(scheme.MPC().PiStag(), scheme.MPT().PiTStag())
	//fmt.Println(scheme.HPC().PiStag(), scheme.HPT().PiTStag())
	//fmt.Println(scheme.Payload().RPM(), scheme.MPC().MassRate())
	//fmt.Println(scheme.Burner().Alpha())

	var delta = 1e-7
	assert.InDelta(
		t,
		scheme.MPC().MassRateOutput().GetState().Value().(float64),
		scheme.HPC().MassRateInput().GetState().Value().(float64),
		delta,
	)
	assert.InDelta(
		t,
		scheme.Burner().MassRateOutput().GetState().Value().(float64),
		scheme.HPT().MassRateInput().GetState().Value().(float64),
		delta,
	)
	assert.InDelta(
		t,
		scheme.HPT().MassRateOutput().GetState().Value().(float64),
		scheme.MPT().MassRateInput().GetState().Value().(float64),
		delta,
	)
	assert.InDelta(
		t,
		scheme.MPT().MassRateOutput().GetState().Value().(float64),
		scheme.LPT().MassRateInput().GetState().Value().(float64),
		delta,
	)
}

func getUnit3nFreeTestScheme() ThreeShaftFreeScheme {
	return get3nFreeTestScheme(
		func(normMassRate, normPiStag float64) float64 {
			return 1
		},
		func(normMassRate, normPiStag float64) float64 {
			return 1
		},
		func(normMassRate, normPiStag float64) float64 {
			return 1
		},
		func(normMassRate, normPiStag float64) float64 {
			return 1
		},
		func(lambdaU, normPiStag float64) float64 {
			return 1
		},
		func(lambdaU, normPiStag float64) float64 {
			return 1
		},
		func(lambdaU, normPiStag float64) float64 {
			return 1
		},
		func(lambdaU, normPiStag float64) float64 {
			return 1
		},
		func(lambdaU, normPiStag float64) float64 {
			return 1
		},
		func(lambdaU, normPiStag float64) float64 {
			return 1
		},
		func(lambda float64) float64 {
			return 1
		},
		func(normRPM float64) float64 {
			return normRPM * normRPM
		},
	)
}

func get3nFreeTestScheme(
	mpcNormEtaChar, mpcNormRpmChar,
	hpcNormEtaChar, hpcNormRpmChar c.CompressorCharFunc,
	hptNormMassRateChar, hptNormEtaChar,
	mptNormMassRateChar, mptNormEtaChar,
	lptNormMassRateChar, lptNormEtaChar c.TurbineCharFunc,
	burnerSigmaFunc func(lambda float64) float64,
	powerChar func(normRPM float64) float64,
) ThreeShaftFreeScheme {
	var mpc = c.NewParametricCompressorNode(
		mpcMassRate0, mpcPi0, mpcRPM0, mpcEta0, t0, p0, precision,
		mpcNormEtaChar, mpcNormRpmChar,
	)
	var mpcPipe = c.NewPressureLossNode(mpcSigma)

	var hpc = c.NewParametricCompressorNode(
		hpcMassRate0, hpcPi0, hpcRPM0, hpcEta0, t0, p0, precision,
		hpcNormEtaChar, hpcNormRpmChar,
	)
	var hpcPipe = c.NewPressureLossNode(hpcSigma)

	var burner = c.NewParametricBurnerNode(
		fuel.GetCH4(), tFuel, t0, etaBurn, lambdaIn0, p0*mpcPi0*hpcPi0, tStagIn0,
		hpcMassRate0, fuelMassRateRel0, precision, burnerSigmaFunc,
	)

	var hpt = c.NewSimpleParametricTurbineNode(
		hptMassRate0, hptPi0, hptEta0, t0, p0, hptDMean, precision,
		0, 0, 0,
		hptNormMassRateChar, hptNormEtaChar,
	)
	var hptPipe = c.NewPressureLossNode(hptSigma)

	var mpt = c.NewSimpleParametricTurbineNode(
		mptMassRate0, mptPi0, mptEta0, t0, p0, mptDMean, precision,
		0, 0, 0,
		mptNormMassRateChar, mptNormEtaChar,
	)
	var mptPipe = c.NewPressureLossNode(mptSigma)

	var lpt = c.NewSimpleParametricTurbineNode(
		lptMassRate0, lptPi0, lptEta0, t0, p0, lptDMean, precision,
		0, 0, 0,
		lptNormMassRateChar, lptNormEtaChar,
	)
	var lptPipe = c.NewPressureLossNode(lptSigma)

	var payload = c.NewPayload(payloadRPM0, -payloadPower0, powerChar)

	return NewThreeShaftFreeScheme(
		gases.GetAir(), tAtm, pAtm, tGas,
		mpc, mpcPipe, mpt, mptPipe, mpEtaM,
		hpc, hpcPipe, hpt, hptPipe, hpEtaM,
		lpt, lptPipe, burner, payload,
	)
}
