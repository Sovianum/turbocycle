package geometry

import "math"

type BladingGeometryGenerator interface {
	GenerateFromInlet(dMeanIn float64) BladingGeometry
	GenerateFromOutlet(dMeanOut float64) BladingGeometry
	LRelOut() float64
	Elongation() float64
	DeltaRel() float64
	GammaIn() float64
	GammaOut() float64
}

func NewGeneratorFromProfileAngles(lRelOut, elongation, deltaRel, gammaIn, gammaOut float64) BladingGeometryGenerator {
	return &bladingGeometryGenerator{
		lRelOut:    lRelOut,
		elongation: elongation,
		deltaRel:   deltaRel,
		gammaIn:    gammaIn,
		gammaOut:   gammaOut,
	}
}

func NewGeneratorFromTotalAndMeanAngles(lRelOut, elongation, deltaRel, totalAngle, meanAngle float64) BladingGeometryGenerator {
	var gammaIn, gammaOut = GetInnerAndOuterAngles(totalAngle, meanAngle)
	return &bladingGeometryGenerator{
		lRelOut:    lRelOut,
		elongation: elongation,
		deltaRel:   deltaRel,
		gammaIn:    gammaIn,
		gammaOut:   gammaOut,
	}
}

func GetTotalAndMeanLineAngles(gammaIn, gammaOut float64) (float64, float64) {
	var totalAngle = gammaOut - gammaIn
	var meanLineAngle = math.Atan2(0.5*(math.Tan(gammaOut)+math.Tan(gammaIn)), 1)

	return totalAngle, meanLineAngle
}

func GetInnerAndOuterAngles(totalAngle, meanLineAngle float64) (float64, float64) {
	var tanGammaM = math.Tan(meanLineAngle)
	var tanGamma = math.Tan(totalAngle)

	var tanGammaIn = tanGammaM + 1/tanGamma - math.Sqrt(1+tanGammaM*tanGammaM+1/(tanGamma*tanGamma))
	var tanGammaOut = tanGammaIn + totalAngle

	return tanGammaIn, tanGammaOut
}

type bladingGeometryGenerator struct {
	lRelOut    float64
	elongation float64
	deltaRel   float64
	gammaIn    float64
	gammaOut   float64
}

func (gen *bladingGeometryGenerator) GenerateFromInlet(dMeanIn float64) BladingGeometry {
	var _, gammaMean = GetTotalAndMeanLineAngles(gen.gammaIn, gen.gammaOut)

	var elongationRel = gen.lRelOut / gen.elongation
	var bladeWidth = elongationRel / (1 - 2*elongationRel*math.Tan(gammaMean)) * dMeanIn
	var gapWidth = bladeWidth * gen.deltaRel
	var bladingWidth = bladeWidth + gapWidth

	var dMeanOut = dMeanIn + 2*bladingWidth*math.Tan(gammaMean)
	var lOut = dMeanOut * gen.lRelOut
	var dInnerOut = dMeanOut - lOut
	var dOuterOut = dMeanOut + lOut

	var dInnerIn = dInnerOut - 2*bladingWidth*math.Tan(gen.gammaIn)
	var dOuterIn = dOuterOut - 2*bladingWidth*math.Tan(gen.gammaOut)

	var innerProfile = NewAxialProfileLine(0, dInnerIn, gen.gammaIn)
	var outerProfile = NewAxialProfileLine(0, dOuterIn, gen.gammaOut)
	var meanProfile = NewAxialProfileLine(0, dMeanIn, gammaMean)

	return NewBladingGeometry(bladeWidth, gapWidth, innerProfile, meanProfile, outerProfile)
}

func (gen *bladingGeometryGenerator) GenerateFromOutlet(dMeanOut float64) BladingGeometry {
	var _, gammaMean = GetTotalAndMeanLineAngles(gen.gammaIn, gen.gammaOut)

	var bladeWidth = gen.lRelOut / gen.elongation * dMeanOut
	var gapWidth = bladeWidth * gen.deltaRel
	var bladingWidth = bladeWidth + gapWidth

	var lOut = dMeanOut * gen.lRelOut
	var dInnerOut = dMeanOut - lOut
	var dOuterOut = dMeanOut + lOut

	var innerProfile = NewAxialProfileLine(bladingWidth, dInnerOut, gen.gammaIn)
	var outerProfile = NewAxialProfileLine(bladingWidth, dOuterOut, gen.gammaOut)
	var meanProfile = NewAxialProfileLine(bladingWidth, dMeanOut, gammaMean)

	return NewBladingGeometry(bladeWidth, gapWidth, innerProfile, meanProfile, outerProfile)
}

func (gen *bladingGeometryGenerator) LRelOut() float64 {
	return gen.lRelOut
}

func (gen *bladingGeometryGenerator) Elongation() float64 {
	return gen.elongation
}

func (gen *bladingGeometryGenerator) DeltaRel() float64 {
	return gen.deltaRel
}

func (gen *bladingGeometryGenerator) GammaIn() float64 {
	return gen.gammaIn
}

func (gen *bladingGeometryGenerator) GammaOut() float64 {
	return gen.gammaOut
}
