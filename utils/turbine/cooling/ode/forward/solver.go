package forward

import (
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/ode"
	"gonum.org/v1/gonum/floats"
)

type Solver interface {
	Solution(derivative func(x, y float64) float64, x0, y0, xMax, maxStep float64) ode.Solution
}

func NewEulerSolver() Solver {
	return eulerSolver{}
}

type eulerSolver struct{}

func (e eulerSolver) Solution(derivative func(x, y float64) float64, x0, y0, xMax, maxStep float64) ode.Solution {
	var distance = xMax - x0
	var stepNum = floats.Round(distance/maxStep, 1)
	var stepDistance = maxStep * stepNum

	var step float64
	if stepDistance == distance {
		step = maxStep
	} else {
		step = distance / (maxStep + 1)
	}

	return eulerSolution{
		derivative: derivative,
		x0:         x0,
		y0:         y0,
		xMax:       xMax,
		step:       step,
	}
}
