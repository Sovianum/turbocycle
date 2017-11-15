package cooling

import "github.com/Sovianum/turbocycle/utils/turbine/cooling/ode"

type TemperatureSystem interface {
	CpAir(theta float64) float64
	GasTemp(x float64) float64
	AlphaAir(x, theta float64) float64
	AlphaGas(x, theta float64) float64
	K(x, theta float64) float64
}

type TemperatureSolution struct {
	X               []float64
	AirTemperature  []float64
	WallTemperature []float64
}

func NewTemperatureSystem(
	solver ode.Solver,
	airMassRate float64,
	cpAir func(theta float64) float64,
	gasTemp func(x float64) float64,
	alphaAir func(x, theta float64) float64,
	alphaGas func(x, theta float64) float64,
	kFunc func(alphaAir, alphaGas float64) float64,
) TemperatureSystem {
	return temperatureSystem{
		solver: solver,

		airMassRate: airMassRate,
		cpAir:       cpAir,
		gasTemp:     gasTemp,
		alphaAir:    alphaAir,
		alphaGas:    alphaGas,
		kFunc:       kFunc,
	}
}

type temperatureSystem struct {
	solver ode.Solver

	airMassRate float64
	cpAir       func(theta float64) float64
	gasTemp     func(x float64) float64
	alphaAir    func(x, theta float64) float64
	alphaGas    func(x, theta float64) float64
	kFunc       func(alphaAir, alphaGas float64) float64
}

func (system temperatureSystem) Solve(x0, theta0, xMax, maxStep float64) TemperatureSolution {
	var xArr, tAirArr = system.solver.Solution(system.dThetaDX, x0, theta0, xMax, maxStep).Build()

	var wallTempArr = make([]float64, len(xArr))
	for i := 0; i != len(xArr); i++ {
		wallTempArr[i] = system.wallTemp(xArr[i], tAirArr[i])
	}

	return TemperatureSolution{
		X:               xArr,
		AirTemperature:  tAirArr,
		WallTemperature: wallTempArr,
	}
}

func (system temperatureSystem) CpAir(theta float64) float64 {
	return system.cpAir(theta)
}

func (system temperatureSystem) GasTemp(x float64) float64 {
	return system.gasTemp(x)
}

func (system temperatureSystem) AlphaAir(x, theta float64) float64 {
	return system.alphaAir(x, theta)
}

func (system temperatureSystem) AlphaGas(x, theta float64) float64 {
	return system.alphaGas(x, theta)
}

func (system temperatureSystem) K(x, theta float64) float64 {
	return system.kFunc(
		system.alphaAir(x, theta),
		system.alphaGas(x, theta),
	)
}

func (system temperatureSystem) dThetaDX(x, theta float64) float64 {
	var factor = 2 / (system.airMassRate * system.cpAir(theta))
	var k = system.kFunc(
		system.alphaAir(x, theta),
		system.alphaGas(x, theta),
	)
	var result = factor * k * (system.gasTemp(x) - theta)
	return result
}

func (system temperatureSystem) wallTemp(x, theta float64) float64 {
	var term1 = system.gasTemp(x)
	var kFactor = system.kFunc(system.alphaAir(x, theta), system.alphaGas(x, theta)) / system.alphaGas(x, theta)
	var tFactor = kFactor * (system.gasTemp(x) - theta)
	return term1 - kFactor*tFactor
}
