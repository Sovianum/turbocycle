package geometry

import "math"

type BladingGeometryGenerator interface {
	Generate(dInnerIn, dOuterIn float64) BladingGeometry
}

func NewGeneratorFromProfileAngles(lRelOut, elongation, deltaRel, gammaIn, gammaOut float64) BladingGeometryGenerator {
	return &bladingGeometryGenerator{
		lRelOut:lRelOut,
		elongation:elongation,
		deltaRel:deltaRel,
		gammaIn:gammaIn,
		gammaOut:gammaOut,
	}
}

func NewGeneratorFromTotalAndMeanAngles(lRelOut, elongation, deltaRel, totalAngle, meanAngle float64) BladingGeometryGenerator {
	var gammaIn, gammaOut = getInnerAndOuterAngles(totalAngle, meanAngle)
	return &bladingGeometryGenerator{
		lRelOut:lRelOut,
		elongation:elongation,
		deltaRel:deltaRel,
		gammaIn:gammaIn,
		gammaOut:gammaOut,
	}
}

type bladingGeometryGenerator struct {
	lRelOut    float64
	elongation float64
	deltaRel   float64
	gammaIn    float64
	gammaOut   float64
}

func (gen *bladingGeometryGenerator) Generate(dInnerIn, dOuterIn float64) BladingGeometry {
	var _, gammaMean = getTotalAndMeanLineAngles(gen.gammaIn, gen.gammaOut)
	var dMeanIn = (dInnerIn + dOuterIn) / 2

	var innerProfile = NewAxialProfileLine(0, dInnerIn, gen.gammaIn)
	var outerProfile = NewAxialProfileLine(0, dOuterIn, gen.gammaOut)
	var	meanProfile = NewAxialProfileLine(0, dMeanIn, gammaMean)

	var elongationRel = gen.lRelOut / gen.elongation
	var bladeWidth = elongationRel / (1 - 2 * elongationRel * math.Tan(gammaMean)) * dMeanIn
	var gapWidth = bladeWidth * gen.deltaRel

	return NewBladingGeometry(bladeWidth, gapWidth, innerProfile, meanProfile, outerProfile)
}

func getTotalAndMeanLineAngles(gammaIn, gammaOut float64) (float64, float64) {
	var totalAngle = gammaOut - gammaIn
	var meanLineAngle = math.Atan2(0.5*(math.Tan(gammaOut)+math.Tan(gammaIn)), 1)

	return totalAngle, meanLineAngle
}

func getInnerAndOuterAngles(totalAngle, meanLineAngle float64) (float64, float64) {
	var tanGammaM = math.Tan(meanLineAngle)
	var tanGamma = math.Tan(totalAngle)

	var tanGammaIn = tanGammaM + 1/tanGamma - math.Sqrt(1+tanGammaM*tanGammaM+1/(tanGamma*tanGamma))
	var tanGammaOut = tanGammaIn + totalAngle

	return tanGammaIn, tanGammaOut
}
