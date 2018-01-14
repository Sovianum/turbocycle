package variator

import (
	"github.com/Sovianum/turbocycle/core/math"
	"gonum.org/v1/gonum/mat"
)

func NewVariatorSolver(
	systemCall SysCall,
	variators []Variator,
	solverGen math.SolverGenerator,
) math.Solver {
	return &variatorSolver{
		variators:  variators,
		systemCall: systemCall,
		solverGen:  solverGen,
	}
}

type variatorSolver struct {
	variators  []Variator
	systemCall func() (*mat.VecDense, error)
	solverGen  math.SolverGenerator
}

func (solver *variatorSolver) Solve(
	x0 *mat.VecDense,
	precision float64,
	relaxCoef float64,
	iterLimit int,
) (solution *mat.VecDense, err error) {
	var systemSolver, solverErr = solver.solverGen(solver.getResidualFunc())
	if solverErr != nil {
		return nil, solverErr
	}

	return systemSolver.Solve(x0, precision, relaxCoef, iterLimit)
}

func (solver *variatorSolver) getResidualFunc() math.EquationSystem {
	var residualFunc = func(x *mat.VecDense) (*mat.VecDense, error) {
		for i := 0; i != x.Len(); i++ {
			solver.variators[i].SetValue(x.At(i, 0))
		}
		return solver.systemCall()
	}
	return math.NewEquationSystem(residualFunc, len(solver.variators))
}
