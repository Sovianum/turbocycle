package free3n

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
	midBurnerTOut = 1800.
)

func TestNewTripleShaftBurnFreeScheme_Smoke(t *testing.T) {
	t.Skip("skipping TestNewTripleShaftBurnFreeScheme_Smoke; broken")

	var scheme = getUnit3nBurnFreeTestScheme()
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
		mat.NewVecDense(10, []float64{
			1, 1, 1, 1,
			1, 1, 1,
			0.034, 1, 0.02,
		}),
		1e-8, 0.2, 1000,
	)
	assert.Nil(t, err)

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

func getUnit3nBurnFreeTestScheme() ThreeShaftBurnFreeScheme {
	return get3nBurnFreeTestScheme(
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
		func(lambda float64) float64 {
			return 1
		},
	)
}

func get3nBurnFreeTestScheme(
	mpcNormEtaChar, mpcNormRpmChar,
	hpcNormEtaChar, hpcNormRpmChar c.CompressorCharFunc,
	hptNormMassRateChar, hptNormEtaChar,
	mptNormMassRateChar, mptNormEtaChar,
	lptNormMassRateChar, lptNormEtaChar c.TurbineCharFunc,
	burnerSigmaFunc func(lambda float64) float64,
	powerChar func(normRPM float64) float64,
	midBurnerSigmaFunc func(lambda float64) float64,
) ThreeShaftBurnFreeScheme {
	var root = get3nFreeTestScheme(
		mpcNormEtaChar, mpcNormRpmChar,
		hpcNormEtaChar, hpcNormRpmChar,
		hptNormMassRateChar, hptNormEtaChar,
		mptNormMassRateChar, mptNormEtaChar,
		lptNormMassRateChar, lptNormEtaChar,
		burnerSigmaFunc, powerChar,
	)
	return NewThreeShaftBurnFreeScheme(
		gases.GetAir(), tAtm, pAtm, tGas, midBurnerTOut,
		root.LPC(), root.LPCPipe(), root.LPT(), root.LPTPipe(), mpEtaM,
		root.HPC(), root.HPCPipe(), root.HPT(), root.HPTPipe(), hpEtaM,
		root.FT(), root.FTPipe(), root.Burner(),
		root.Payload(), c.NewParametricBurnerNode(
			fuel.GetCH4(), tFuel, t0, etaBurn, lambdaIn0, p0*mpcPi0, midBurnerTOut,
			hpcMassRate0, fuelMassRateRel0, precision, 1, nodes.DefaultN, midBurnerSigmaFunc,
		),
	)
}
