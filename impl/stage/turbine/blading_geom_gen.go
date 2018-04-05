package turbine

import (
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/stage/geometry"
)

const (
	MidLineFactor = 0.5
)

func TRel(hRel float64, gen BladingGeometryGenerator) float64 {
	var factor1 = math.Pi / float64(gen.BladeNumber())
	var factor2 = gen.Elongation() / gen.LRelOut()
	var factor3 = 1 - (1-2*hRel)*gen.LRelOut()
	return factor1 * factor2 * factor3
}

func LRelIn(gen BladingGeometryGenerator) float64 {
	var gammaIn, gammaOut = gen.GammaIn(), gen.GammaOut()
	var _, gammaM = geometry.GetTotalAndMeanLineAngles(gammaIn, gammaOut, MidLineFactor)

	var enom = gen.LRelOut() - gen.LRelOut()/gen.Elongation()*(math.Tan(gammaOut)-math.Tan(gammaIn))
	var denom = 2 * (1 - gen.LRelOut()/gen.Elongation()*math.Tan(gammaM))

	return enom / denom
}

func LRelOutGap(gen BladingGeometryGenerator) float64 {
	return RecalculateLRel(
		gen.LRelOut(),
		gen.DeltaRel()/gen.Elongation(),
		gen.GammaIn(), gen.GammaOut(),
	)
}

func NewBladingGeometry(bladeWidth, gapWidth float64, innerProfile, outerProfile geometry.AxialProfileLine) geometry.BladingGeometry {
	return geometry.NewBladingGeometry(bladeWidth, gapWidth, innerProfile, outerProfile, MidLineFactor)
}

func NewBladingGenerator(lRelOut, elongation, deltaRel, gammaIn, gammaOut, approxTRel float64) BladingGeometryGenerator {
	return &bladingGeometryGenerator{
		lRelOut:    lRelOut,
		elongation: elongation,
		deltaRel:   deltaRel,
		gammaIn:    gammaIn,
		gammaOut:   gammaOut,
		approxTRel: approxTRel,
	}
}

type BladingGeometryGenerator interface {
	geometry.BladingGeometryGenerator
	GenerateFromInlet(dMeanIn float64) geometry.BladingGeometry
	GenerateFromOutlet(dMeanOut float64) geometry.BladingGeometry
	LRelOut() float64
	BladeNumber() int
}

type bladingGeometryGenerator struct {
	lRelOut    float64
	elongation float64
	deltaRel   float64
	gammaIn    float64
	gammaOut   float64
	approxTRel float64
}

func (gen *bladingGeometryGenerator) GenerateFromInlet(dMeanIn float64) geometry.BladingGeometry {
	var _, gammaMean = geometry.GetTotalAndMeanLineAngles(gen.gammaIn, gen.gammaOut, MidLineFactor)

	var elongationRel = gen.lRelOut / gen.elongation
	var dMeanOut = 1 / (1 - 2*elongationRel*(1+gen.deltaRel)*math.Tan(gammaMean)) * dMeanIn
	var bladeWidth = dMeanOut * elongationRel
	var gapWidth = bladeWidth * gen.deltaRel
	var bladingWidth = bladeWidth + gapWidth

	var lOut = dMeanOut * gen.lRelOut
	var dInnerOut = dMeanOut - lOut
	var dOuterOut = dMeanOut + lOut

	var innerProfile = geometry.NewAxialProfileLine(bladingWidth, dInnerOut, gen.gammaIn)
	var outerProfile = geometry.NewAxialProfileLine(bladingWidth, dOuterOut, gen.gammaOut)

	return NewBladingGeometry(bladeWidth, gapWidth, innerProfile, outerProfile)
}

func (gen *bladingGeometryGenerator) GenerateFromOutlet(dMeanOut float64) geometry.BladingGeometry {
	var bladeWidth = gen.lRelOut / gen.elongation * dMeanOut
	var gapWidth = bladeWidth * gen.deltaRel
	var bladingWidth = bladeWidth + gapWidth

	var lOut = dMeanOut * gen.lRelOut
	var dInnerOut = dMeanOut - lOut
	var dOuterOut = dMeanOut + lOut

	var innerProfile = geometry.NewAxialProfileLine(bladingWidth, dInnerOut, gen.gammaIn)
	var outerProfile = geometry.NewAxialProfileLine(bladingWidth, dOuterOut, gen.gammaOut)

	return NewBladingGeometry(bladeWidth, gapWidth, innerProfile, outerProfile)
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

func (gen *bladingGeometryGenerator) BladeNumber() int {
	var baRel = gen.elongation
	var lRelOut = gen.lRelOut

	return common.RoundInt(math.Pi * baRel / lRelOut * 1 / gen.approxTRel)
}
