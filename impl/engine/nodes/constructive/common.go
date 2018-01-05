package constructive

import "math"

func NormMassRate(massRate, massRate0, t, t0, p, p0 float64) float64 {
	return massRate / massRate0 * math.Sqrt(t0/t) * p / p0
}

func MassRate(normMassRate, massRate0, t, t0, p, p0 float64) float64 {
	return normMassRate / (math.Sqrt(t0/t) * p / p0) * massRate0
}

func NormRpm(rpm, rpm0, t, t0 float64) float64 {
	return rpm / rpm0 * math.Sqrt(t0/t)
}

func Rpm(normRpm, rpm0, t, t0 float64) float64 {
	return normRpm / math.Sqrt(t0/t) * rpm0
}
