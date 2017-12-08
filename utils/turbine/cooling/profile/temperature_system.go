package profile

type TemperatureSystem interface {
	Solve(t0, theta0, tMax, maxStep float64) TemperatureSolution
	CpAir(theta float64) float64
	GasTemp(x float64) float64
	AlphaAir(x, theta float64) float64
	AlphaGas(x, theta float64) float64
	K(x, theta float64) float64
}