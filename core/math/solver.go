package math

import "gonum.org/v1/gonum/mat"

type Solver interface {
	Solve(x0 *mat.VecDense, precision float64, relaxCoef float64, iterLimit int) (solution *mat.VecDense, err error)
}

type SolverGenerator func(system EquationSystem) (Solver, error)
