package geometry

func NewIncompleteGeneratorFromProfileAngles(elongation, deltaRel, gammaIn, gammaOut float64) IncompleteBladingGeometryGenerator {
	return &incompleteBladingGeometryGenerator{
		elongation: elongation,
		deltaRel:   deltaRel,
		gammaIn:    gammaIn,
		gammaOut:   gammaOut,
	}
}

func NewIncompleteGeneratorFromTotalAndMeanAngles(elongation, deltaRel, totalAngle, meanAngle float64) IncompleteBladingGeometryGenerator {
	var gammaIn, gammaOut = GetInnerAndOuterAngles(totalAngle, meanAngle)
	return &incompleteBladingGeometryGenerator{
		elongation: elongation,
		deltaRel:   deltaRel,
		gammaIn:    gammaIn,
		gammaOut:   gammaOut,
	}
}

type IncompleteBladingGeometryGenerator interface {
	GetGenerator(lRelOut float64) BladingGeometryGenerator
	Elongation() float64
	DeltaRel() float64
	GammaIn() float64
	GammaOut() float64
}

type incompleteBladingGeometryGenerator struct {
	elongation float64
	deltaRel   float64
	gammaIn    float64
	gammaOut   float64
}

func (gen *incompleteBladingGeometryGenerator) GetGenerator(lRelOut float64) BladingGeometryGenerator {
	return NewGeneratorFromProfileAngles(lRelOut, gen.elongation, gen.deltaRel, gen.gammaIn, gen.gammaOut)
}

func (gen *incompleteBladingGeometryGenerator) Elongation() float64 {
	return gen.elongation
}

func (gen *incompleteBladingGeometryGenerator) DeltaRel() float64 {
	return gen.deltaRel
}

func (gen *incompleteBladingGeometryGenerator) GammaIn() float64 {
	return gen.gammaIn
}

func (gen *incompleteBladingGeometryGenerator) GammaOut() float64 {
	return gen.gammaOut
}
