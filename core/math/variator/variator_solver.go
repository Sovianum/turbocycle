package variator

import (
	"math/rand"

	"github.com/Sovianum/turbocycle/core/math"
	"gonum.org/v1/gonum/mat"
)

func NoiseUniformly(vec *mat.VecDense, noiseFactor float64) *mat.VecDense {
	r, _ := vec.Dims()
	result := mat.NewVecDense(r, nil)
	for i := 0; i != r; i++ {
		x := vec.At(i, 0) * (1 + rand.Float64()*noiseFactor)
		result.SetVec(i, x)
	}
	return result
}

func Noise(vec *mat.VecDense, noiseFactorVec *mat.VecDense) *mat.VecDense {
	r, _ := vec.Dims()
	result := mat.NewVecDense(r, nil)
	for i := 0; i != r; i++ {
		x := vec.At(i, 0) * (1 + rand.Float64()*noiseFactorVec.At(i, 0))
		result.SetVec(i, x)
	}
	return result
}

type VariatorSolver interface {
	math.Solver
	GetInit() *mat.VecDense
}

func NewVariatorSolver(
	systemCall SysCall,
	variators []Variator,
	solverGen math.SolverGenerator,
) VariatorSolver {
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

func (solver *variatorSolver) GetInit() *mat.VecDense {
	data := make([]float64, len(solver.variators))
	for i, v := range solver.variators {
		data[i] = v.GetValue()
	}
	return mat.NewVecDense(len(data), data)
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
