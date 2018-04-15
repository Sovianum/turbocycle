package boundary

import (
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/ode"
	"gonum.org/v1/gonum/mat"
)

type CoefFunc func(x float64) float64

// Solver represents solver for the equation of type
// y'' + fArr(x)*y' + gArr(x)*y = hArr(x)
type Solver interface {
	Solve() (ode.Solution, error)
}

func NewSolverFromArrays(xArr, fArr, gArr, hArr []float64, startBC, endBC BC) Solver {
	step := xArr[1] - xArr[0]

	return &solver{
		xArr:     xArr,
		fArr:     fArr,
		gArr:     gArr,
		hArr:     hArr,
		pointNum: len(xArr),
		step:     step,
		startBC:  startBC,
		endBC:    endBC,
	}
}

func NewSolverFromFuncs(f, g, h CoefFunc, startBC, endBC BC, x0, xMax, maxStep float64) Solver {
	distance := xMax - x0
	stepNum := int(distance/maxStep) + 1
	pointNum := stepNum + 1
	step := distance / float64(stepNum)

	result := &solver{
		xArr:     make([]float64, pointNum),
		fArr:     make([]float64, pointNum),
		gArr:     make([]float64, pointNum),
		hArr:     make([]float64, pointNum),
		pointNum: pointNum,
		step:     step,
		startBC:  startBC,
		endBC:    endBC,
	}

	for i := 0; i != pointNum; i++ {
		xi := x0 + float64(i)*step
		result.xArr[i] = xi
		result.fArr[i] = f(xi)
		result.gArr[i] = g(xi)
		result.hArr[i] = h(xi)
	}
	return result
}

type solver struct {
	xArr []float64
	fArr []float64
	gArr []float64
	hArr []float64

	pointNum int
	step     float64

	startBC BC
	endBC   BC

	lhs *mat.Dense
	rhs *mat.VecDense
}

func (s *solver) Solve() (ode.Solution, error) {
	s.fillMatrices()

	solution := mat.NewVecDense(s.pointNum, nil)

	if err := solution.SolveVec(s.lhs, s.rhs); err != nil {
		return nil, err
	}
	return newBoundarySolution(s.xArr, solution.RawVector().Data, s.step), nil
}

func (s *solver) fillMatrices() {
	s.lhs = mat.NewDense(s.pointNum, s.pointNum, nil)
	s.rhs = mat.NewVecDense(s.pointNum, nil)

	// first and last lines are excluded
	// (they will be filled with boundary conditions)
	for i := 1; i != s.pointNum-1; i++ {
		fi := s.fArr[i]
		gi := s.gArr[i]
		hi := s.hArr[i]

		s.lhs.Set(i, i-1, 2-fi*s.step)
		s.lhs.Set(i, i, 2*(gi*s.step*s.step-2))
		s.lhs.Set(i, i+1, 2+fi*s.step)

		s.rhs.SetVec(i, 2*hi*s.step*s.step)
	}

	s.lhs.Set(0, 0, s.startBC.LHS0(s.step))
	s.lhs.Set(0, 1, s.startBC.LHS1(s.step))
	s.rhs.SetVec(0, s.startBC.RHS(s.step))

	s.lhs.Set(s.pointNum-1, s.pointNum-1, s.endBC.LHS0(s.step))
	s.lhs.Set(s.pointNum-1, s.pointNum-2, s.endBC.LHS1(s.step))
	s.rhs.SetVec(s.pointNum-1, s.endBC.RHS(s.step))
}
