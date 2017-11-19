package geometry

import "math"

func TRel(hRel float64, gen BladingGeometryGenerator) float64 {
	var factor1 = math.Pi / float64(gen.BladeNumber())
	var factor2 = gen.Elongation() / gen.LRelOut()
	var factor3 = 1 - (1 - 2 * hRel) * gen.LRelOut()
	return factor1 * factor2 * factor3
}

func LRelIn(gen BladingGeometryGenerator) float64 {
	var gammaIn, gammaOut = gen.GammaIn(), gen.GammaOut()
	var _, gammaM = GetTotalAndMeanLineAngles(gammaIn, gammaOut)

	var enom = gen.LRelOut() - gen.LRelOut()/gen.Elongation()*(math.Tan(gammaOut)-math.Tan(gammaIn))
	var denom = 2 * (1 - gen.LRelOut()/gen.Elongation()*math.Tan(gammaM))

	return enom / denom
}
