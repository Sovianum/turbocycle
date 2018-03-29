package compressor

import (
	"math"

	"github.com/Sovianum/turbocycle/impl/stage/geometry"
)

type StageGeometryGenerator interface {
	Generate(dOutIn float64) geometry.StageGeometry
	StatorGenerator() BladingGeometryGenerator
	RotorGenerator() BladingGeometryGenerator
}

func NewStageGeometryGenerator(
	dRelIn float64,
	rotorIncompleteGen,
	statorIncompleteGen IncompleteBladingGeometryGenerator,
) StageGeometryGenerator {
	lRel := rotorIncompleteGen.Elongation()
	deltaRel := rotorIncompleteGen.DeltaRel()
	gammaIn := rotorIncompleteGen.GammaIn()
	gammaOut := rotorIncompleteGen.GammaOut()

	deltaTerm := (1 + deltaRel) / lRel * (1 - dRelIn)

	factor := (1 + deltaTerm*math.Tan(gammaIn)/dRelIn) / (1 + deltaTerm*math.Tan(gammaOut))
	dRelInStator := dRelIn * factor
	return &stageGeometryGenerator{
		rotorGenerator:  rotorIncompleteGen.GetGenerator(dRelIn),
		statorGenerator: statorIncompleteGen.GetGenerator(dRelInStator),
	}
}

type stageGeometryGenerator struct {
	rotorGenerator  BladingGeometryGenerator
	statorGenerator BladingGeometryGenerator
}

func (gen *stageGeometryGenerator) Generate(dOutIn float64) geometry.StageGeometry {
	rotor := gen.rotorGenerator.Generate(dOutIn)
	dOutStator := rotor.OuterProfile().Diameter(rotor.XGapOut())
	var stator = gen.statorGenerator.Generate(dOutStator)
	return geometry.NewStageGeometry(stator, rotor)
}

func (gen *stageGeometryGenerator) StatorGenerator() BladingGeometryGenerator {
	return gen.statorGenerator
}

func (gen *stageGeometryGenerator) RotorGenerator() BladingGeometryGenerator {
	return gen.rotorGenerator
}
