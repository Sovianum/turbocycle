package profile

import (
	"github.com/Sovianum/turbocycle/utils/turbine/cooling"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/ode"
	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"gonum.org/v1/gonum/mat"
)

func NewConvectiveTemperatureSystem(
	solver ode.Solver,
	airMassRate float64,
	cpAir func(theta float64) float64,
	gasTemp func(x float64) float64,
	alphaAir func(x, theta float64) float64,
	alphaGas func(x, theta float64) float64,
	wallThk func(x float64) float64,
	lambdaM func(t float64) float64,
	segment geom.Segment,
) TemperatureSystem {
	return &convectiveTemperatureSystem{
		solver: solver,

		airMassRate:     airMassRate,
		cpAir:           cpAir,
		gasTemp:         gasTemp,
		alphaAir:        alphaAir,
		alphaGas:        alphaGas,
		wallThk:         wallThk,
		lambdaM:         lambdaM,
		segment:         segment,
		lengthParameter: 0,
	}
}

type convectiveTemperatureSystem struct {
	solver ode.Solver

	airMassRate float64
	cpAir       func(theta float64) float64
	gasTemp     func(x float64) float64
	alphaAir    cooling.AlphaLaw
	alphaGas    cooling.AlphaLaw

	wallThk func(x float64) float64
	lambdaM func(t float64) float64

	segment      geom.Segment
	solutionStep float64

	lengthParameter float64
}

func (system *convectiveTemperatureSystem) Solve(t0, theta0, tMax, maxStep float64) TemperatureSolution {
	var solution = system.solver.Solution(system.dThetaDX, t0, theta0, tMax, maxStep)
	system.solutionStep = solution.Step()

	var tArr, tAirArr = solution.Build()

	var solutionLen = len(tArr)
	var tSolution = TemperatureSolution{
		ParametricCoord:  tArr,
		AirTemperature:   tAirArr,
		X:                make([]float64, solutionLen),
		Y:                make([]float64, solutionLen),
		LengthCoord:      make([]float64, solutionLen),
		AlphaAir:         make([]float64, solutionLen),
		AlphaGas:         make([]float64, solutionLen),
		WallTemperature:  make([]float64, solutionLen),
		HeatTransferCoef: make([]float64, solutionLen),
	}
	system.extendSolutionArray(&tSolution, t0, system.solutionStep, solutionLen)
	return tSolution
}

func (system *convectiveTemperatureSystem) CpAir(theta float64) float64 {
	return system.cpAir(theta)
}

func (system *convectiveTemperatureSystem) GasTemp(lengthCoord float64) float64 {
	return system.gasTemp(lengthCoord)
}

func (system *convectiveTemperatureSystem) AlphaAir(lengthCoord, theta float64) float64 {
	return system.alphaAir(lengthCoord, theta)
}

func (system *convectiveTemperatureSystem) AlphaGas(lengthCoord, theta float64) float64 {
	return system.alphaGas(lengthCoord, theta)
}

func (system *convectiveTemperatureSystem) K(x, theta float64) float64 {
	return system.heatTransferCoef(x, theta)
}

func (system *convectiveTemperatureSystem) dThetaDX(t, theta float64) float64 {
	var segmentDerivative = geom.ApproxDerivative1(system.segment, t, system.solutionStep)
	var geomFactor = mat.Norm(segmentDerivative, 2)

	var massRateFactor = 2 / (system.airMassRate * system.cpAir(theta))
	var k = system.heatTransferCoef(system.lengthParameter, theta)
	var result = geomFactor * massRateFactor * k * (system.gasTemp(system.lengthParameter) - theta)

	system.lengthParameter += geomFactor * system.solutionStep // side effect is used to prevent length calculation on each step
	return result
}

func (system *convectiveTemperatureSystem) wallTemp(lengthCoord, theta float64) float64 {
	var term1 = system.gasTemp(lengthCoord)
	var kFactor = system.heatTransferCoef(lengthCoord, theta) / system.alphaGas(lengthCoord, theta)
	var tFactor = system.gasTemp(lengthCoord) - theta
	return term1 - kFactor*tFactor
}

func (system *convectiveTemperatureSystem) heatTransferCoef(lengthCoord, theta float64) float64 {
	var alphaAir = system.alphaAir(lengthCoord, theta)
	var alphaGas = system.alphaGas(lengthCoord, theta)
	var delta = system.wallThk(lengthCoord)
	var lambdaM = system.lambdaM(theta)

	var enom = alphaAir
	var denom = 1 + alphaAir*(1/alphaGas+delta/lambdaM)
	return enom / denom
}

func (system *convectiveTemperatureSystem) extendSolutionArray(solution *TemperatureSolution, tStart, step float64, pointNum int) {
	var currLengthCoord float64 = 0
	var currT float64 = 0

	for i := 0; i != pointNum; i++ {
		var point = system.segment.Point(currT)
		solution.X[i] = point.At(0, 0)
		solution.Y[i] = point.At(1, 0)
		solution.LengthCoord[i] = currLengthCoord
		solution.AlphaAir[i] = system.AlphaAir(currLengthCoord, solution.AirTemperature[i])
		solution.AlphaGas[i] = system.alphaGas(currLengthCoord, solution.AirTemperature[i])
		solution.WallTemperature[i] = system.wallTemp(currLengthCoord, solution.AirTemperature[i])
		solution.HeatTransferCoef[i] = system.heatTransferCoef(currLengthCoord, solution.AirTemperature[i])

		var segmentDerivative = geom.ApproxDerivative1(system.segment, currT, step)
		var lengthStep = mat.Norm(segmentDerivative, 2) * step

		currT += step
		currLengthCoord += lengthStep
	}
}

