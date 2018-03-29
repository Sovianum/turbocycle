package geometry

import "math"

type BladingGeometryGenerator interface {
	Elongation() float64
	DeltaRel() float64
	GammaIn() float64
	GammaOut() float64
}

func GetTotalAndMeanLineAngles(gammaIn, gammaOut, midLineFactor float64) (float64, float64) {
	var totalAngle = gammaOut - gammaIn
	tan1 := math.Tan(gammaIn)
	tan2 := math.Tan(gammaOut)
	f := midLineFactor

	var meanLineAngle = math.Atan2((1-f)*tan1+f*tan2, 1)

	return totalAngle, meanLineAngle
}

func GetInnerAndOuterAngles(totalAngle, meanLineAngle, midLineFactor float64) (float64, float64) {
	x1 := math.Tan(meanLineAngle)
	x1Sqr := x1 * x1
	x2 := math.Tan(totalAngle)
	x2Sqr := x2 * x2
	f := midLineFactor

	term1 := 1 - x1*x2
	term2 := math.Sqrt(1 - (1-2*f)*x1*x2 + 4*f*(1-f)*x2Sqr + x1Sqr*x2Sqr)

	tanGammaIn := (term1 - term2) / (2 * (1 - f) * x2)

	gammaIn := math.Atan(tanGammaIn)
	gammaOut := totalAngle + gammaIn

	return gammaIn, gammaOut
}
