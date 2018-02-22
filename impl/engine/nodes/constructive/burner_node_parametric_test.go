package constructive

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

const (
	lambdaIn0        = 0.15
	pStagIn0         = 1e6
	tStagIn0         = 500.
	massRateIn0      = 40.
	fuelMassRateRel0 = 1e-2
	sigmaBurn0       = 0.95
)

func TestParametricBurnerNode_Process_Smoke_Constant(t *testing.T) {
	var bn = getTestParametricBurner(func(lambda float64) float64 {
		return sigmaBurn0
	})

	var tStagIn = tStagIn0 * 1.2
	var pStagIn = pStagIn0 * 1.2
	var massRate = massRateIn0

	var err = graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()),
			states.NewTemperaturePortState(tStagIn),
			states.NewPressurePortState(pStagIn),
			states.NewMassRatePortState(massRate),
		},
		[]graph.Port{
			bn.GasInput(), bn.TemperatureInput(), bn.PressureInput(), bn.MassRateInput(),
		},
	)
	if err != nil {
		panic(err)
	}

	bn.Process()
	assert.NotNil(t, bn)

	var expectedPOut = sigmaBurn0 * pStagIn
	assert.InDelta(
		t, expectedPOut, bn.PStagOut(), 1e-7,
	)

	var expectedMassRate = massRate * (1 + fuelMassRateRel0)
	assert.InDelta(
		t, expectedMassRate, bn.MassRateOutput().GetState().(states.MassRatePortState).MassRate, 1e-7,
	)

	bn.SetFuelRateRel(0)
	bn.Process()
	assert.InDelta(t, tStagIn, bn.TStagOut(), 1e-7)

	bn.SetFuelRateRel(1e-2)
	bn.Process()
	var tStag1 = bn.TStagOut()

	bn.SetFuelRateRel(2e-2)
	bn.Process()
	var tStag2 = bn.TStagOut()

	assert.True(t, tStag1 < tStag2)
}

func TestParametricBurnerNode_Process_Smoke_Solver(t *testing.T) {
	var bn = getTestParametricBurner(func(lambda float64) float64 {
		return sigmaBurn0
	})

	var tStagIn = tStagIn0 * 1.2
	var pStagIn = pStagIn0 * 1.2
	var massRate = massRateIn0

	var err = graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()),
			states.NewTemperaturePortState(tStagIn),
			states.NewPressurePortState(pStagIn),
			states.NewMassRatePortState(massRate),
		},
		[]graph.Port{
			bn.GasInput(), bn.TemperatureInput(), bn.PressureInput(), bn.MassRateInput(),
		},
	)
	if err != nil {
		panic(err)
	}

	var fuelRateVariator = variator.FromCallables(
		bn.FuelRateRel, bn.SetFuelRateRel,
	)

	var needTStagOut = 2000.
	var systemCall = func() (*mat.VecDense, error) {
		bn.Process()
		return mat.NewVecDense(1, []float64{bn.TStagOut() - needTStagOut}), nil
	}

	var solverGen = newton.NewUniformNewtonSolverGen(1e-4, newton.NoLog)
	var variatorSolver = variator.NewVariatorSolver(systemCall, []variator.Variator{fuelRateVariator}, solverGen)

	var _, solveErr = variatorSolver.Solve(mat.NewVecDense(1, []float64{fuelMassRateRel0}), 1e-8, 1, 100)
	assert.Nil(t, solveErr)

	assert.True(t, bn.FuelRateRel() > fuelMassRateRel0)
}

func getTestParametricBurner(sigmaFunc func(lambda float64) float64) ParametricBurnerNode {
	return NewParametricBurnerNode(
		fuel.GetCH4(), tFuel, t0, etaBurn, lambdaIn0,
		pStagIn0, tStagIn0, massRateIn0, fuelMassRateRel0, 1e-3,
		sigmaFunc,
	)
}
