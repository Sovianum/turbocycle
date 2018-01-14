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
	midBurnerTOut = 1800.
)

func TestNewTripleShaftBurnFreeScheme_Smoke(t *testing.T) {
	var scheme = getUnit3nBurnFreeTestScheme()
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
		mat.NewVecDense(10, []float64{
			1, 1, 1, 1,
			1, 1, 1,
			0.034, 1, 0.02,
		}),
		1e-8, 0.2, 1000,
	)
	assert.Nil(t, err)

	//fmt.Println()
	//fmt.Println(scheme.LPT().PiTStag())
	//fmt.Println(scheme.MPC().PiStag(), scheme.MPT().PiTStag())
	//fmt.Println(scheme.HPC().PiStag(), scheme.HPT().PiTStag())
	//fmt.Println(scheme.Payload().RPM(), scheme.MPC().MassRate())
	//fmt.Println(scheme.Burner().Alpha(), scheme.Burner().FuelRateRel())
	//fmt.Println(scheme.MidBurner().Alpha(), scheme.MidBurner().FuelRateRel())

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
		root.MPC(), root.MPCPipe(), root.MPT(), root.MPTPipe(), mpEtaM,
		root.HPC(), root.HPCPipe(), root.HPT(), root.HPTPipe(), hpEtaM,
		root.LPT(), root.LPTPipe(), root.Burner(),
		root.Payload(), c.NewParametricBurnerNode(
			fuel.GetCH4(), tFuel, t0, etaBurn, lambdaIn0, p0*mpcPi0, midBurnerTOut,
			hpcMassRate0, fuelMassRateRel0, precision, midBurnerSigmaFunc,
		),
	)
}
