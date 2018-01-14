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

func (solver *newtonSolver) Solve(
	x0 *mat.VecDense, precision float64, relaxCoef float64, iterLimit int,
) (solution *mat.VecDense, err error) {
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
		x, y, err = solver.getNewState(x, y, relaxCoef)
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

func (solver *newtonSolver) getNewState(currX, currY *mat.VecDense, relaxCoef float64) (newX, newY *mat.VecDense, err error) {
	var xOffset, offsetErr = solver.getXOffset(currX, currY)

	if offsetErr != nil {
		return nil, nil, offsetErr
	}
	xOffset.ScaleVec(relaxCoef, xOffset)

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
		if zeroErr := solver.getZeroRowColErr(jacobian); zeroErr != nil {
			return nil, fmt.Errorf("%s: %s", err.Error(), zeroErr.Error())
		}

		return nil, err
	}

	return xOffset, nil
}

func (solver *newtonSolver) getZeroRowColErr(m mat.Matrix) error {
	var zeroRows []int
	var zeroCols []int

	var r, c = m.Dims()

	for i := 0; i != r; i++ {
		if solver.isZeroRow(m, i) {
			zeroRows = append(zeroRows, i)
		}
	}
	for j := 0; j != c; j++ {
		if solver.isZeroColumn(m, j) {
			zeroCols = append(zeroCols, j)
		}
	}

	var errMsg = ""
	if zeroRows != nil {
		errMsg += fmt.Sprintf("rows %v consist of zeros (y[i] does not depend on x)", zeroRows)
	}
	if zeroCols != nil {
		errMsg += fmt.Sprintf("cols %v consist of zeros (x[i] does not affect y)", zeroCols)
	}

	if errMsg != "" {
		return fmt.Errorf(errMsg)
	}
	return nil
}

func (solver *newtonSolver) isZeroRow(m mat.Matrix, row int) bool {
	var result = true
	var _, c = m.Dims()

	for j := 0; j != c; j++ {
		if m.At(row, j) != 0 {
			result = false
			break
		}
	}
	return result
}

func (solver *newtonSolver) isZeroColumn(m mat.Matrix, col int) bool {
	var result = true
	var r, _ = m.Dims()

	for i := 0; i != r; i++ {
		if m.At(i, col) != 0 {
			result = false
			break
		}
	}
	return result
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
