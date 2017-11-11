package gdf

import "math"

func Epsilon(lambda float64, k float64) float64 {
	var x = 1 - (k-1)/(k+1)*lambda*lambda
	var y = 1 / (k - 1)
	return math.Pow(x, y)
}

func Pi(lambda float64, k float64) float64 {
	var x = 1 - (k-1)/(k+1)*lambda*lambda
	var y = k / (k - 1)
	return math.Pow(x, y)
}

func Tau(lambda float64, k float64) float64 {
	return 1 - (k-1)/(k+1)*lambda*lambda
}

func ACrit(k float64, R float64, TStag float64) float64 {
	return math.Sqrt(2 * k / (k + 1) * R * TStag)
}

//returns q(lambda, k) * m
func Q(lambda float64, k float64, R float64) float64 {
	var coef = math.Sqrt(2 * k / ((k + 1) * R))
	var x = 1 - (k-1)/(k+1)*lambda*lambda
	var y = 1 / (k - 1)
	return coef * lambda * math.Pow(x, y)
}

func Mach(lambda float64, k float64) float64 {
	var num = 2 / (k + 1) * lambda * lambda
	var denom = 1 - (k-1)/(k+1)*lambda*lambda

	return math.Sqrt(num / denom)
}

func Lambda(mach float64, k float64) float64 {
	var num = (k + 1) / 2 * mach * mach
	var denom = 1 + 2*mach*mach

	return math.Sqrt(num / denom)
}
