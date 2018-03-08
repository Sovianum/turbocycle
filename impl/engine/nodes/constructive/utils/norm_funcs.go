package utils

import "math"

func NormMassRate(massRate, massRate0, t, t0, p, p0 float64) float64 {
	pRel := p / p0
	tRel := t / t0

	return massRate / massRate0 * math.Sqrt(tRel) / pRel
}

func MassRate(normMassRate, massRate0, t, t0, p, p0 float64) float64 {
	pRel := p / p0
	tRel := t / t0
	return normMassRate * pRel / math.Sqrt(tRel) * massRate0
}

func NormRpm(rpm, rpm0, t, t0 float64) float64 {
	return rpm / rpm0 * math.Sqrt(t0/t)
}

func Rpm(normRpm, rpm0, t, t0 float64) float64 {
	return normRpm / math.Sqrt(t0/t) * rpm0
}
