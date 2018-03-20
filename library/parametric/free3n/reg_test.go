package free3n

import (
	"fmt"
	"testing"

	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	c "github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

func TestNewTripleShaftRegFreeScheme_Smoke(t *testing.T) {
	t.Skip("skipping TestNewTripleShaftRegFreeScheme_Smoke; broken")

	var scheme = getUnit3nRegFreeTestScheme()
	var network, err = scheme.GetNetwork()
	assert.Nil(t, err)

	err = network.Solve(1, 2, 100, 1e-3)
	assert.Nil(t, err)

	burner := scheme.Burner()
	fmt.Println(burner.TemperatureInput().GetState(), burner.TemperatureOutput().GetState())

	var sysCall = variator.SysCallFromNetwork(
		network, scheme.Assembler().GetVectorPort(), 1, 2, 100, 1e-3,
	)
	var variators = scheme.Variators()
	var solverGen = newton.NewUniformNewtonSolverGen(1e-4, newton.DefaultLog)

	var variatorSolver = variator.NewVariatorSolver(sysCall, variators, solverGen)

	_, err = variatorSolver.Solve(
		mat.NewVecDense(9, []float64{
			0.8, 1, 1, 1,
			1, 1, 1,
			0.05, 1,
		}),
		1e-4, 0.2, 10000,
	)
	assert.Nil(t, err)

	fmt.Println(scheme.FT().PiTStag())
	fmt.Println(scheme.LPC().PiStag(), scheme.LPT().PiTStag())
	fmt.Println(scheme.HPC().PiStag(), scheme.HPT().PiTStag())
	fmt.Println(scheme.Payload().RPM(), scheme.LPC().MassRate())
	fmt.Println(scheme.Burner().Alpha())

	var delta = 1e-7
	assert.InDelta(
		t,
		scheme.LPC().MassRateOutput().GetState().Value().(float64),
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
		scheme.LPT().MassRateInput().GetState().Value().(float64),
		delta,
	)
	assert.InDelta(
		t,
		scheme.LPT().MassRateOutput().GetState().Value().(float64),
		scheme.FT().MassRateInput().GetState().Value().(float64),
		delta,
	)
}

func getUnit3nRegFreeTestScheme() ThreeShaftRegFreeScheme {
	return get3nRegFreeTestScheme(
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

func get3nRegFreeTestScheme(
	mpcNormEtaChar, mpcNormRpmChar,
	hpcNormEtaChar, hpcNormRpmChar c.CompressorCharFunc,
	hptNormMassRateChar, hptNormEtaChar,
	mptNormMassRateChar, mptNormEtaChar,
	lptNormMassRateChar, lptNormEtaChar c.TurbineCharFunc,
	burnerSigmaFunc func(lambda float64) float64,
	powerChar func(normRPM float64) float64,
) ThreeShaftRegFreeScheme {
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
		hpcMassRate0, fuelMassRateRel0, precision, 1, nodes.DefaultN, burnerSigmaFunc,
	)

	var hpt = c.NewSimpleParametricTurbineNode(
		hptMassRate0, hptPi0, hptEta0, t0, p0, hptDMean, precision,
		0, 0, 0,
		hptNormMassRateChar, hptNormEtaChar,
	)
	var hptPipe = c.NewPressureLossNode(hptSigma)

	var mpt = c.NewSimpleParametricTurbineNode(
		mptMassRate0, mptPi0*2, mptEta0, t0, p0, mptDMean, precision,
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

	var regenerator = c.NewParametricRegeneratorNode(
		gases.GetAir(), gases.GetAir(), mpcMassRate0, mpcMassRate0, tGas*0.75, 500,
		3e5, 3e5, 20, 20, 0.5,
		1e-3, 1.5e-3, 1e-3, 1, nodes.DefaultN,
		c.FrowardTDrop,
		c.DefaultNuFunc, c.DefaultNuFunc,
	)

	var breaker = helper.NewComplexCycleBreakNode(gases.GetAir(), 1300, 1e5, 5.5)

	return NewThreeShaftRegFreeScheme(
		gases.GetAir(), tAtm, pAtm, tGas,
		mpc, mpcPipe, mpt, mptPipe, mpEtaM,
		hpc, hpcPipe, hpt, hptPipe, hpEtaM,
		lpt, lptPipe, burner, payload,
		regenerator, breaker,
	)
}
