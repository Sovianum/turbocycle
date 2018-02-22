package newton

import (
	"math"
	"testing"

	math2 "github.com/Sovianum/turbocycle/core/math"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

func TestNewtonSolver_Solve_Linear_1d(t *testing.T) {
	var eqSystem = math2.NewEquationSystem(func(x *mat.VecDense) (*mat.VecDense, error) {
		return mat.NewVecDense(1, []float64{x.At(0, 0) - 1}), nil
	}, 1)

	var solver, solverErr = NewUniformNewtonSolver(eqSystem, 1e-3, NoLog)
	assert.Nil(t, solverErr)

	var solution, solutionErr = solver.Solve(mat.NewVecDense(1, []float64{0}), 1e-8, 1, 100)
	assert.Nil(t, solutionErr)

	assert.InDelta(t, 1, solution.At(0, 0), 1e-7)
}

func TestNewtonSolver_Solve_Linear_2d(t *testing.T) {
	var eqSystem = math2.NewEquationSystem(func(x *mat.VecDense) (*mat.VecDense, error) {
		return mat.NewVecDense(2, []float64{
			x.At(0, 0) + x.At(1, 0) - 3,
			x.At(0, 0) - x.At(1, 0) - 1,
		}), nil
	}, 2)

	var solver, solverErr = NewUniformNewtonSolver(eqSystem, 1e-3, NoLog)
	assert.Nil(t, solverErr)

	var solution, solutionErr = solver.Solve(mat.NewVecDense(2, []float64{0, 0}), 1e-8, 1, 100)
	assert.Nil(t, solutionErr)

	assert.InDelta(t, 2., solution.At(0, 0), 1e-7)
	assert.InDelta(t, 1., solution.At(1, 0), 1e-7)
}

func TestNewtonSolver_Solve_NonLinear_1d(t *testing.T) {
	var eqSystem = math2.NewEquationSystem(func(x *mat.VecDense) (*mat.VecDense, error) {
		return mat.NewVecDense(1, []float64{
			math.Pow(x.At(0, 0), 3) - 1,
		}), nil
	}, 1)

	var solver, solverErr = NewUniformNewtonSolver(eqSystem, 1e-3, NoLog)
	assert.Nil(t, solverErr)

	var solution, solutionErr = solver.Solve(mat.NewVecDense(1, []float64{0}), 1e-8, 1, 100)
	assert.Nil(t, solutionErr)

	assert.InDelta(t, 1, solution.At(0, 0), 1e-7)
}

func TestNewtonSolver_Solve_NonLinear_2d(t *testing.T) {
	var eqSystem = math2.NewEquationSystem(func(x *mat.VecDense) (*mat.VecDense, error) {
		var x1 = x.At(0, 0)
		var x2 = x.At(1, 0)
		return mat.NewVecDense(2, []float64{
			x1*x1 - x2 - 3,
			x1 - x2 - 1,
		}), nil
	}, 2)

	var solver, solverErr = NewUniformNewtonSolver(eqSystem, 1e-3, NoLog)
	assert.Nil(t, solverErr)

	var solution, solutionErr = solver.Solve(mat.NewVecDense(2, []float64{-0.5, -0.5}), 1e-8, 1, 100)
	assert.Nil(t, solutionErr)

	assert.InDelta(t, -1, solution.At(0, 0), 1e-7)
	assert.InDelta(t, -2, solution.At(1, 0), 1e-7)
}
