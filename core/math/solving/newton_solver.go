package solving

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func NewNewtonSolver(eqSystem EquationSystem, x0, derivativeSteps *mat.VecDense) (Solver, error) {
	if derivativeSteps.Len() != eqSystem.Order() {
		return nil, fmt.Errorf(
			"derivative step size %d does not match eqSystem order %d", derivativeSteps.Len(), eqSystem.Order(),
		)
	}

	if x0.Len() != eqSystem.Order() {
		return nil, fmt.Errorf(
			"x0 size %d does not match eqSystem order %d", x0.Len(), eqSystem.Order(),
		)
	}

	return &newtonSolver{
		eqSystem:   eqSystem,
		initValues: x0,
		steps:      derivativeSteps,
	}, nil
}

func NewUniformNewtonSolver(eqSystem EquationSystem, x0 *mat.VecDense, derivativeStep float64) (Solver, error) {
	if x0.Len() != eqSystem.Order() {
		return nil, fmt.Errorf(
			"x0 size %d does not match eqSystem order %d", x0.Len(), eqSystem.Order(),
		)
	}

	var derivativeSteps = mat.NewVecDense(eqSystem.Order(), nil)
	for i := 0; i != eqSystem.Order(); i++ {
		derivativeSteps.SetVec(i, 0)
	}

	return &newtonSolver{
		eqSystem:   eqSystem,
		initValues: x0,
		steps:      derivativeSteps,
	}, nil
}

type newtonSolver struct {
	eqSystem   EquationSystem
	initValues *mat.VecDense
	steps      *mat.VecDense // used to calculate partial derivatives
}

func (solver *newtonSolver) Solve(precision float64, iterLimit int) (solution *mat.VecDense, err error) {
	var x = solver.initValues
	var y = mat.NewVecDense(solver.order(), nil)

	var converged = false
	for i := 0; i != iterLimit; i++ {
		x, y, err = solver.getNewState(x, y)
		if mat.Norm(y, 2) <= precision {
			converged = true
			break
		}
	}

	if !converged {
		return nil, fmt.Errorf("failed to converge")
	}

	return y, nil
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
