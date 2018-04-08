package compressor

import (
	"math"

	"github.com/Sovianum/turbocycle/impl/stage/geometry"
)

func NewIncompleteStageGeomGen(rotorGen, statorGen IncompleteBladingGeometryGenerator) IncompleteStageGeometryGenerator {
	return &incompleteStageGeometryGenerator{
		rotorGen:  rotorGen,
		statorGen: statorGen,
	}
}

type IncompleteStageGeometryGenerator interface {
	Generate(dRelIn float64) StageGeometryGenerator
}

type incompleteStageGeometryGenerator struct {
	rotorGen  IncompleteBladingGeometryGenerator
	statorGen IncompleteBladingGeometryGenerator
}

func (gen *incompleteStageGeometryGenerator) Generate(dRelIn float64) StageGeometryGenerator {
	return NewStageGeometryGenerator(dRelIn, gen.rotorGen, gen.statorGen)
}

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
	baRel := rotorIncompleteGen.Elongation()
	deltaRel := rotorIncompleteGen.DeltaRel()
	gammaIn := rotorIncompleteGen.GammaIn()
	gammaOut := rotorIncompleteGen.GammaOut()

	axialFactor := (1 + deltaRel) / baRel * (1 - dRelIn) / 2

	factor := (1 + 2*axialFactor*math.Tan(gammaIn)/dRelIn) / (1 + 2*axialFactor*math.Tan(gammaOut))
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
	stator := gen.statorGenerator.Generate(dOutStator)
	return geometry.NewStageGeometry(stator, rotor)
}

func (gen *stageGeometryGenerator) StatorGenerator() BladingGeometryGenerator {
	return gen.statorGenerator
}

func (gen *stageGeometryGenerator) RotorGenerator() BladingGeometryGenerator {
	return gen.rotorGenerator
}
