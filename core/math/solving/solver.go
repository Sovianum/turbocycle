package solving

import "gonum.org/v1/gonum/mat"

type Solver interface {
	Solve(precision float64, iterLimit int) (solution *mat.VecDense, err error)
}
