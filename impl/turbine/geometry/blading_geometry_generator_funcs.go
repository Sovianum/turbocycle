package geometry

import "math"

func LRelIn(gen BladingGeometryGenerator) float64 {
	var gammaIn, gammaOut = gen.GammaIn(), gen.GammaOut()
	var _, gammaM = GetTotalAndMeanLineAngles(gammaIn, gammaOut)

	var enom = gen.LRelOut() - gen.LRelOut()/gen.Elongation()*(math.Tan(gammaOut)-math.Tan(gammaIn))
	var denom = 2 * (1 - gen.LRelOut()/gen.Elongation()*math.Tan(gammaM))

	return enom / denom
}