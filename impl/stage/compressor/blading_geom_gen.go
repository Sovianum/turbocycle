package compressor

import (
	"math"

	"github.com/Sovianum/turbocycle/impl/stage/geometry"
)

func NewBladingGeometry(bladeWidth, gapWidth float64, innerProfile, outerProfile geometry.AxialProfileLine) geometry.BladingGeometry {
	dRel := outerProfile.Diameter(0) / innerProfile.Diameter(0)
	midLineFactor := math.Sqrt((1 + dRel*dRel) / 2)
	return geometry.NewBladingGeometry(bladeWidth, gapWidth, innerProfile, outerProfile, midLineFactor)
}

func NewBladingGenerator(dRelIn, elongation, deltaRel, gammaIn, gammaOut float64) BladingGeometryGenerator {
	return &bladingGeometryGenerator{
		dRelIn:     dRelIn,
		elongation: elongation,
		deltaRel:   deltaRel,
		gammaIn:    gammaIn,
		gammaOut:   gammaOut,
	}
}

type BladingGeometryGenerator interface {
	geometry.BladingGeometryGenerator
	Generate(dOutIn float64) geometry.BladingGeometry
	DRelIn() float64
}

type bladingGeometryGenerator struct {
	dRelIn     float64
	elongation float64
	deltaRel   float64
	gammaIn    float64
	gammaOut   float64
	approxTRel float64
}

func (gen *bladingGeometryGenerator) Generate(dOutIn float64) geometry.BladingGeometry {
	dInIn := dOutIn * gen.dRelIn
	innerProfile := geometry.NewAxialProfileLine(0, dInIn, gen.gammaIn)
	outerProfile := geometry.NewAxialProfileLine(0, dOutIn, gen.gammaOut)

	bladeWidth := dOutIn * (1 - gen.dRelIn) / 2 / gen.elongation
	gapWidth := bladeWidth * gen.deltaRel

	return NewBladingGeometry(bladeWidth, gapWidth, innerProfile, outerProfile)
}

func (gen *bladingGeometryGenerator) DRelIn() float64 {
	return gen.dRelIn
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
