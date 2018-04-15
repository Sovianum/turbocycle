package boundary

import (
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/ode"
	"gonum.org/v1/gonum/mat"
)

type CoefFunc func(x float64) float64

// Solver represents solver for the equation of type
// y'' + f(x)*y' + g(x)*y = h(x)
type Solver interface {
	Solve() (ode.Solution, error)
}

func NewSolver(f, g, h CoefFunc, startBC, endBC BC, x0, xMax, maxStep float64) Solver {
	return &solver{
		f:       f,
		g:       g,
		h:       h,
		startBC: startBC,
		endBC:   endBC,

		x0:      x0,
		xMax:    xMax,
		maxStep: maxStep,
	}
}

type solver struct {
	f CoefFunc
	g CoefFunc
	h CoefFunc

	startBC BC
	endBC   BC

	x0      float64
	xMax    float64
	maxStep float64

	step     float64
	pointNum int

	lhs *mat.Dense
	rhs *mat.VecDense
}

func (s *solver) Solve() (ode.Solution, error) {
	s.init()
	s.fillMatrices()

	solution := mat.NewVecDense(s.pointNum, nil)

	if err := solution.SolveVec(s.lhs, s.rhs); err != nil {
		return nil, err
	}
	return newBoundarySolution(solution, s.step), nil
}

func (s *solver) fillMatrices() {
	// first and last lines are excluded
	// (they will be filled with boundary conditions)
	for i := 1; i != s.pointNum-1; i++ {
		xi := s.step*float64(i) + s.x0
		fi := s.f(xi)
		gi := s.g(xi)
		hi := s.h(xi)

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

func (s *solver) init() {
	distance := s.xMax - s.x0
	stepNum := int(distance/s.maxStep) + 1
	s.pointNum = stepNum + 1
	s.step = distance / float64(stepNum)

	s.lhs = mat.NewDense(s.pointNum, s.pointNum, nil)
	s.rhs = mat.NewVecDense(s.pointNum, nil)
}
