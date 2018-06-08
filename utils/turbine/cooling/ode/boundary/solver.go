package boundary

import (
	"gonum.org/v1/gonum/mat"
)

type CoefFunc func(x float64) float64

// Solver represents solver for the equation of type
// y'' + fArr(x)*y' + gArr(x)*y = hArr(x)
type Solver interface {
	Solve() (Solution, error)
}

func NewSolverFromArrays(xArr, fArr, gArr, hArr []float64, startBC, endBC BC) Solver {
	return &solver{
		xArr:     xArr,
		fArr:     fArr,
		gArr:     gArr,
		hArr:     hArr,
		pointNum: len(xArr),
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

	startBC BC
	endBC   BC

	lhs *mat.Dense
	rhs *mat.VecDense
}

func (s *solver) Solve() (Solution, error) {
	s.fillMatrices()

	solution := mat.NewVecDense(s.pointNum, nil)

	if err := solution.SolveVec(s.lhs, s.rhs); err != nil {
		return nil, err
	}
	return newBoundarySolution(s.xArr, solution.RawVector().Data), nil
}

func (s *solver) fillMatrices() {
	s.lhs = mat.NewDense(s.pointNum, s.pointNum, nil)
	s.rhs = mat.NewVecDense(s.pointNum, nil)
	firstDerAcceptor := make([]float64, 3)
	secondDerAcceptor := make([]float64, 3)

	// first and last lines are excluded
	// (they will be filled with boundary conditions)
	for i := 1; i != s.pointNum-1; i++ {
		fi := s.fArr[i]
		gi := s.gArr[i]
		hi := s.hArr[i]

		xPrev := s.xArr[i-1]
		xCurr := s.xArr[i]
		xNext := s.xArr[i+1]

		getFirstDerivativeFactors(xPrev, xCurr, xNext, firstDerAcceptor)
		getSecondDerivativeFactors(xPrev, xCurr, xNext, secondDerAcceptor)

		s.lhs.Set(i, i-1, fi*firstDerAcceptor[0]+secondDerAcceptor[0])
		s.lhs.Set(i, i, fi*firstDerAcceptor[1]+secondDerAcceptor[1]+gi)
		s.lhs.Set(i, i+1, fi*firstDerAcceptor[2]+secondDerAcceptor[2])

		s.rhs.SetVec(i, hi)
	}

	startStep := s.xArr[1] - s.xArr[0]
	s.lhs.Set(0, 0, s.startBC.LHS0(startStep))
	s.lhs.Set(0, 1, s.startBC.LHS1(startStep))
	s.rhs.SetVec(0, s.startBC.RHS(startStep))

	endStep := s.xArr[len(s.xArr)-2] - s.xArr[len(s.xArr)-1]
	s.lhs.Set(s.pointNum-1, s.pointNum-1, s.endBC.LHS0(endStep))
	s.lhs.Set(s.pointNum-1, s.pointNum-2, s.endBC.LHS1(endStep))
	s.rhs.SetVec(s.pointNum-1, s.endBC.RHS(endStep))
}

func getFirstDerivativeFactors(xPrev, xCurr, xNext float64, acceptor []float64) {
	dxPrev := xCurr - xPrev
	dxCurr := xNext - xCurr

	acceptor[0] = -dxCurr / (dxPrev * (dxPrev + dxCurr))
	acceptor[1] = (dxPrev - dxCurr) / (dxCurr * dxPrev)
	acceptor[2] = dxPrev / (dxCurr * (dxPrev + dxCurr))
}

func getSecondDerivativeFactors(xPrev, xCurr, xNext float64, acceptor []float64) {
	dxPrev := xCurr - xPrev
	dxCurr := xNext - xCurr

	acceptor[0] = 2 / (dxPrev * (dxPrev + dxCurr))
	acceptor[1] = -2 / (dxCurr * dxPrev)
	acceptor[2] = 2 / (dxCurr * (dxPrev + dxCurr))
}
