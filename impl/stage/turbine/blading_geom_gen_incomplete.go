package turbine

func NewIncompleteGenerator(elongation, deltaRel, gammaIn, gammaOut, approxTRel float64) IncompleteBladingGeometryGenerator {
	return &incompleteBladingGeometryGenerator{
		elongation: elongation,
		deltaRel:   deltaRel,
		gammaIn:    gammaIn,
		gammaOut:   gammaOut,
		approxTRel: approxTRel,
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
	approxTRel float64
}

func (gen *incompleteBladingGeometryGenerator) GetGenerator(lRelOut float64) BladingGeometryGenerator {
	return NewBladingGenerator(lRelOut, gen.elongation, gen.deltaRel, gen.gammaIn, gen.gammaOut, gen.approxTRel)
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
