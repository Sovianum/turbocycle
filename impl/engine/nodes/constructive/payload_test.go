package constructive

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

const (
	power0 = 1e6
)

func TestPayload_Process_Linear(t *testing.T) {
	var f = func(normRpm float64) float64 {
		return 1 - (normRpm-1)*1e5
	}

	var node = NewPayload(rpm0, power0, f)

	var normRpm = 0.5
	node.SetNormRPM(normRpm)
	var err = node.Process()
	assert.Nil(t, err)

	assert.InDelta(t, power0*f(normRpm), -node.PowerOutput().GetState().Value().(float64), 1e-7)
	assert.InDelta(t, rpm0*normRpm, node.RPMOutput().GetState().Value().(float64), 1e-7)
}

func TestPayload_Process_Linear_Solver(t *testing.T) {
	var f = func(normRpm float64) float64 {
		return normRpm
	}

	var node = NewPayload(rpm0, power0, f)

	var v = variator.FromCallables(node.NormRPM, node.SetNormRPM)
	var systemCall = func() (*mat.VecDense, error) {
		node.Process()
		return mat.NewVecDense(1, []float64{node.Power() - power0/2}), nil
	}

	var solverGen = newton.NewUniformNewtonSolverGen(1e-4, newton.NoLog)
	var variatorSolver = variator.NewVariatorSolver(systemCall, []variator.Variator{v}, solverGen)

	var _, err = variatorSolver.Solve(mat.NewVecDense(1, []float64{1}), 1e-8, 1, 100)
	assert.Nil(t, err)

	assert.InDelta(t, rpm0/2, node.RPM(), 1e-7)
}
