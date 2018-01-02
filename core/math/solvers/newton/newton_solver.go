package newton

import (
	"fmt"

	"github.com/Sovianum/turbocycle/core/math"
	"gonum.org/v1/gonum/mat"
)

func NewUniformNewtonSolverGen(derivativeStep float64) math.SolverGenerator {
	return func(system math.EquationSystem) (math.Solver, error) {
		return NewUniformNewtonSolver(system, derivativeStep)
	}
}

func NewNewtonSolver(eqSystem math.EquationSystem, derivativeSteps *mat.VecDense) (math.Solver, error) {
	if derivativeSteps.Len() != eqSystem.Order() {
		return nil, fmt.Errorf(
			"derivative step size %d does not match eqSystem order %d", derivativeSteps.Len(), eqSystem.Order(),
		)
	}

	return &newtonSolver{
		eqSystem: eqSystem,
		steps:    derivativeSteps,
	}, nil
}

func NewUniformNewtonSolver(eqSystem math.EquationSystem, derivativeStep float64) (math.Solver, error) {
	var derivativeSteps = mat.NewVecDense(eqSystem.Order(), nil)
	for i := 0; i != eqSystem.Order(); i++ {
		derivativeSteps.SetVec(i, derivativeStep)
	}

	return &newtonSolver{
		eqSystem: eqSystem,
		steps:    derivativeSteps,
	}, nil
}

type newtonSolver struct {
	eqSystem math.EquationSystem
	steps    *mat.VecDense // used to calculate partial derivatives
}

func (solver *newtonSolver) Solve(x0 *mat.VecDense, precision float64, iterLimit int) (solution *mat.VecDense, err error) {
	if x0.Len() != solver.order() {
		return nil, fmt.Errorf(
			"x0 size %d does not match eqSystem order %d", x0.Len(), solver.order(),
		)
	}

	var x = x0
	var y, yErr = solver.eqSystem.GetResiduals(x)
	if yErr != nil {
		return nil, yErr
	}

	var converged = false
	for i := 0; i != iterLimit; i++ {
		x, y, err = solver.getNewState(x, y)
		if err != nil {
			return nil, err
		}
		if mat.Norm(y, 2) <= precision {
			converged = true
			break
		}
	}

	if !converged {
		return nil, fmt.Errorf("failed to converge")
	}

	return x, nil
}

func (solver *newtonSolver) getNewState(currX, currY *mat.VecDense) (newX, newY *mat.VecDense, err error) {
	var xOffset, offsetErr = solver.getXOffset(currX, currY)

	if offsetErr != nil {
		return nil, nil, offsetErr
	}

	newX = mat.NewVecDense(solver.order(), nil)
	newX.AddVec(currX, xOffset)

	newY, err = solver.eqSystem.GetResiduals(newX)

	if err != nil {
		return nil, nil, err
	}

	return newX, newY, nil
}

func (solver *newtonSolver) getXOffset(currX, currY *mat.VecDense) (*mat.VecDense, error) {
	var jacobian, jacobianErr = solver.getJacobian(currX, currY)
	if jacobianErr != nil {
		return nil, jacobianErr
	}

	var rhs = mat.NewVecDense(solver.order(), nil)
	rhs.ScaleVec(-1, currY)

	var xOffset = mat.NewVecDense(solver.order(), nil)
	var err = xOffset.SolveVec(jacobian, rhs)

	if err != nil {
		return nil, err
	}

	return xOffset, nil
}

func (solver *newtonSolver) getJacobian(currX, currY *mat.VecDense) (*mat.Dense, error) {
	var derMatrix = mat.NewDense(solver.order(), solver.order(), nil)
	for j := 0; j != solver.order(); j++ {
		var derVector, vecErr = solver.getPartialDerivatives(currX, currY, j)
		if vecErr != nil {
			return nil, vecErr
		}

		derMatrix.SetCol(j, derVector.RawVector().Data)
	}
	return derMatrix, nil
}

func (solver *newtonSolver) getPartialDerivatives(currX, currY *mat.VecDense, i int) (*mat.VecDense, error) {
	var newX = variate(currX, i, solver.steps.At(i, 0))
	var newY, err = solver.eqSystem.GetResiduals(newX)

	if err != nil {
		return nil, err
	}

	var derVec = mat.NewVecDense(solver.order(), nil)
	derVec.SubVec(newY, currY)
	derVec.ScaleVec(1/solver.steps.At(i, 0), derVec)

	return derVec, nil
}

func (solver *newtonSolver) order() int {
	return solver.eqSystem.Order()
}

func variate(xVec *mat.VecDense, i int, offset float64) *mat.VecDense {
	var result = mat.NewVecDense(xVec.Len(), nil)
	result.CloneVec(xVec)
	result.SetVec(i, result.At(i, 0)+offset)
	return result
}
