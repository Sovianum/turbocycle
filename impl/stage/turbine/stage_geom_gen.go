package turbine

import (
	"github.com/Sovianum/turbocycle/impl/stage/geometry"
)

func NewIncompleteStageGeometryGenerator(
	incompleteStatorGen,
	incompleteRotorGen IncompleteBladingGeometryGenerator,
) IncompleteStageGeometryGenerator {
	return &incompleteStageGeometryGenerator{
		incompleteStatorGen: incompleteStatorGen,
		incompleteRotorGen:  incompleteRotorGen,
	}
}

type IncompleteStageGeometryGenerator interface {
	GetGenerator(lRelIn float64) StageGeometryGenerator
}

type incompleteStageGeometryGenerator struct {
	incompleteStatorGen IncompleteBladingGeometryGenerator
	incompleteRotorGen  IncompleteBladingGeometryGenerator
}

func (gen *incompleteStageGeometryGenerator) GetGenerator(lRelIn float64) StageGeometryGenerator {
	xRelStator := (1 + gen.incompleteStatorGen.DeltaRel()) / gen.incompleteStatorGen.Elongation()
	xRelRotor := (1 + gen.incompleteRotorGen.DeltaRel()) / gen.incompleteRotorGen.Elongation()

	lRelOutStator := RecalculateLRel(lRelIn, xRelStator, gen.incompleteStatorGen.GammaIn(), gen.incompleteStatorGen.GammaOut())
	lRelOutRotor := RecalculateLRel(lRelOutStator, xRelRotor, gen.incompleteStatorGen.GammaIn(), gen.incompleteStatorGen.GammaOut())
	return NewStageGeometryGenerator(lRelOutRotor, gen.incompleteStatorGen, gen.incompleteRotorGen)
}

type StageGeometryGenerator interface {
	GenerateFromRotorInlet(dMeanIn float64) geometry.StageGeometry
	GenerateFromStatorInlet(dMeanIn float64) geometry.StageGeometry
	StatorGenerator() BladingGeometryGenerator
	RotorGenerator() BladingGeometryGenerator
}

func NewStageGeometryGenerator(
	lRelOut float64,
	statorIncompleteGen,
	rotorIncompleteGen IncompleteBladingGeometryGenerator,
) StageGeometryGenerator {
	gammaIn, gammaOut := statorIncompleteGen.GammaIn(), statorIncompleteGen.GammaOut()
	elongation := statorIncompleteGen.Elongation()
	deltaRel := statorIncompleteGen.DeltaRel()

	xRel := -(1 + deltaRel) / elongation
	lRelOutStator := RecalculateLRel(lRelOut, xRel, gammaIn, gammaOut)

	rotorLRelOut := lRelOut
	statorLRelOut := lRelOutStator

	return &stageGeometryGenerator{
		statorGenerator: statorIncompleteGen.GetGenerator(statorLRelOut),
		rotorGenerator:  rotorIncompleteGen.GetGenerator(rotorLRelOut),
	}
}

type stageGeometryGenerator struct {
	statorGenerator BladingGeometryGenerator
	rotorGenerator  BladingGeometryGenerator
}

func (gen *stageGeometryGenerator) GenerateFromRotorInlet(dMeanIn float64) geometry.StageGeometry {
	var rotor = gen.rotorGenerator.GenerateFromInlet(dMeanIn)
	var stator = gen.statorGenerator.GenerateFromOutlet(dMeanIn)
	return geometry.NewStageGeometry(stator, rotor)
}

func (gen *stageGeometryGenerator) GenerateFromStatorInlet(dMeanIn float64) geometry.StageGeometry {
	var stator = gen.statorGenerator.GenerateFromInlet(dMeanIn)
	var rotor = gen.rotorGenerator.GenerateFromInlet(stator.MeanProfile().Diameter(stator.XGapOut()))
	return geometry.NewStageGeometry(stator, rotor)
}

func (gen *stageGeometryGenerator) StatorGenerator() BladingGeometryGenerator {
	return gen.statorGenerator
}

func (gen *stageGeometryGenerator) RotorGenerator() BladingGeometryGenerator {
	return gen.rotorGenerator
}
