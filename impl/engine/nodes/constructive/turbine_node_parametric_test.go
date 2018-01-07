package constructive

import (
	"fmt"
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

const (
	piT0  = 2
	etaT0 = 0.9
)

func TestParametricBlockedTurbineNode_Process_Smoke_Constant(t *testing.T) {
	var turbine = getTestParametricTurbine(
		func(float64, float64) float64 { return 1. },
		func(float64, float64) float64 { return 1. },
	)

	turbine.Process()

	var k = gases.KMean(gases.GetAir(), turbine.TStagOut(), turbine.TStagIn(), nodes.DefaultN)
	var expectedTt = tBlockedT * (1 - (1-math.Pow(piT0, (1-k)/k))*etaT0)
	assert.True(
		t,
		common.ApproxEqual(expectedTt, turbine.TStagOut(), 0.01),
		fmt.Sprintf("Expected T_t %f, got %f", expectedTt, turbine.TStagOut()),
	)

	assert.True(
		t,
		common.ApproxEqual(piT0, turbine.PiTStag(), 0.01),
		fmt.Sprintf("Expected pi_t %f, got %f", piT0, turbine.PiTStag()),
	)

	// labour is calculated by stagnation parameters
	var expectedPt = turbine.PStagIn() / turbine.PiTStag()
	assert.True(
		t,
		common.ApproxEqual(expectedPt, turbine.PStagOut(), 0.01),
		fmt.Sprintf("Expected p_t %f, got %f", expectedPt, turbine.PStagOut()),
	)
}

func TestParametricBlockedTurbineNode_Process_Smoke_NonLinear(t *testing.T) {
	var f = func(x float64) float64 { return 1 - 0.05*(1-x)*(1-x) }

	var turbine = getTestParametricTurbine(
		func(float64, float64) float64 { return 1. },
		func(normMassRate float64, normPiStag float64) float64 {
			return f(normMassRate)
		},
	)

	var piTVariator = variator.FromCallables(
		turbine.NormPiT, turbine.SetNormPiT,
	)

	turbine.Process()
	var piT = piT0 * 1.2
	var k = gases.KMean(gases.GetAir(), turbine.TStagOut(), turbine.TStagIn(), nodes.DefaultN)
	var expectedTt = tBlockedT * (1 - (1-math.Pow(piT, (1-k)/k))*etaT0)

	var systemCall = func() (*mat.VecDense, error) {
		turbine.Process()
		return mat.NewVecDense(1, []float64{turbine.TStagOut() - expectedTt}), nil
	}

	var solverGen = newton.NewUniformNewtonSolverGen(1e-4)
	var variatorSolver = variator.NewVariatorSolver(systemCall, []variator.Variator{piTVariator}, solverGen)

	var _, err = variatorSolver.Solve(mat.NewVecDense(1, []float64{1}), 1e-8, 100)
	assert.Nil(t, err)

	assert.True(
		t,
		common.ApproxEqual(
			piT, turbine.PiTStag(), 1e-2,
		),
		fmt.Sprintf("expected %.3f, got %.3f", piT, turbine.PiTStag()),
	)
}

func getTestParametricTurbine(normMassRateChar, normEtaChar TurbineCharFunc) ParametricTurbineNode {
	var turbine = NewSimpleParametricTurbineNode(
		massRate0, piT0, etaT0, states.StandardTemperature, states.StandardPressure,
		1, 1e-3, 0, 0, 0,
		normMassRateChar, normEtaChar,
	)
	var err = graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()), states.NewTemperaturePortState(tBlockedT),
			states.NewPressurePortState(pBlockedT),
			states.NewRPMPortState(rpm0),
		},
		[]graph.Port{
			turbine.GasInput(), turbine.TemperatureInput(), turbine.PressureInput(),
			turbine.RPMInput(),
		},
	)
	if err != nil {
		panic(err)
	}
	return turbine
}
