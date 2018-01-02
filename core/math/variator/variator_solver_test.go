package variator

import (
	math2 "math"
	"testing"

	"github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

func TestVariatorSolver_Solve_Linear_1d(t *testing.T) {
	var x = 1.
	var xVariator = VariatorFromPointer(&x)

	var systemCall = func() (*mat.VecDense, error) {
		return mat.NewVecDense(1, []float64{x - 2}), nil
	}

	var solverGen = func(system math.EquationSystem) (math.Solver, error) {
		return newton.NewUniformNewtonSolver(system, 1e-3)
	}

	var variatorSolver = NewVariatorSolver(
		systemCall, []Variator{xVariator}, solverGen,
	)

	var _, err = variatorSolver.Solve(
		mat.NewVecDense(1, []float64{3}),
		1e-8, 100,
	)

	assert.Nil(t, err)
	assert.InDelta(t, 2., x, 1e-7)
}

func TestVariatorSolver_Solve_Linear_2d(t *testing.T) {
	var x = 1.
	var xVariator = VariatorFromPointer(&x)

	var y = 1.
	var yVariator = VariatorFromPointer(&y)

	var systemCall = func() (*mat.VecDense, error) {
		return mat.NewVecDense(2, []float64{
			x + y - 3,
			x - y - 1,
		}), nil
	}

	var solverGen = func(system math.EquationSystem) (math.Solver, error) {
		return newton.NewUniformNewtonSolver(system, 1e-3)
	}

	var variatorSolver = NewVariatorSolver(
		systemCall, []Variator{xVariator, yVariator}, solverGen,
	)

	var _, err = variatorSolver.Solve(
		mat.NewVecDense(2, []float64{0, 0}),
		1e-8, 100,
	)

	assert.Nil(t, err)
	assert.InDelta(t, 2., x, 1e-7)
	assert.InDelta(t, 1., y, 1e-7)
}

func TestVariatorSolver_Solve_NonLinear_1d(t *testing.T) {
	var x = 0.
	var xVariator = VariatorFromPointer(&x)

	var systemCall = func() (*mat.VecDense, error) {
		return mat.NewVecDense(1, []float64{math2.Pow(x, 3) - 1}), nil
	}

	var solverGen = func(system math.EquationSystem) (math.Solver, error) {
		return newton.NewUniformNewtonSolver(system, 1e-3)
	}

	var variatorSolver = NewVariatorSolver(
		systemCall, []Variator{xVariator}, solverGen,
	)

	var _, err = variatorSolver.Solve(
		mat.NewVecDense(1, []float64{3}),
		1e-8, 100,
	)

	assert.Nil(t, err)
	assert.InDelta(t, 1., x, 1e-7)
}

func TestVariatorSolver_Solve_NonLinear_2d(t *testing.T) {
	var x = 1.
	var xVariator = VariatorFromPointer(&x)

	var y = 1.
	var yVariator = VariatorFromPointer(&y)

	var systemCall = func() (*mat.VecDense, error) {
		return mat.NewVecDense(2, []float64{
			x*x - y - 3,
			x - y - 1,
		}), nil
	}

	var solverGen = func(system math.EquationSystem) (math.Solver, error) {
		return newton.NewUniformNewtonSolver(system, 1e-3)
	}

	var variatorSolver = NewVariatorSolver(
		systemCall, []Variator{xVariator, yVariator}, solverGen,
	)

	var _, err = variatorSolver.Solve(
		mat.NewVecDense(2, []float64{0.5, 0.5}),
		1e-8, 100,
	)

	assert.Nil(t, err)
	assert.InDelta(t, 2., x, 1e-7)
	assert.InDelta(t, 1., y, 1e-7)
}
