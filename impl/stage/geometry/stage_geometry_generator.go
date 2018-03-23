package geometry

import "math"

type StageGeometryGenerator interface {
	GenerateFromRotorInlet(dMeanIn float64) StageGeometry
	GenerateFromStatorInlet(dMeanIn float64) StageGeometry
	StatorGenerator() BladingGeometryGenerator
	RotorGenerator() BladingGeometryGenerator
}

func NewStageGeometryGenerator(
	lRelOut float64,
	statorIncompleteGen,
	rotorIncompleteGen IncompleteBladingGeometryGenerator,
) StageGeometryGenerator {
	var getFactor = func(lRelOut float64) float64 {
		var gammaIn, gammaOut = statorIncompleteGen.GammaIn(), statorIncompleteGen.GammaOut()
		var _, gammaMean = GetTotalAndMeanLineAngles(gammaIn, gammaOut)

		var elongation = statorIncompleteGen.Elongation()
		var deltaRel = statorIncompleteGen.DeltaRel()

		var enom1 = elongation
		var enom2 = -(1 + deltaRel) * (math.Tan(gammaOut) - math.Tan(gammaIn))

		var denom1 = elongation
		var denom2 = -2 * (1 + deltaRel) * lRelOut * math.Tan(gammaMean)

		return (enom1 + enom2) / (denom1 + denom2)
	}

	var rotorLRelOut = lRelOut
	var statorLRelOut = rotorLRelOut * getFactor(lRelOut)

	return &stageGeometryGenerator{
		statorGenerator: statorIncompleteGen.GetGenerator(statorLRelOut),
		rotorGenerator:  rotorIncompleteGen.GetGenerator(rotorLRelOut),
	}
}

type stageGeometryGenerator struct {
	statorGenerator BladingGeometryGenerator
	rotorGenerator  BladingGeometryGenerator
}

func (gen *stageGeometryGenerator) GenerateFromRotorInlet(dMeanIn float64) StageGeometry {
	var rotor = gen.rotorGenerator.GenerateFromInlet(dMeanIn)
	var stator = gen.statorGenerator.GenerateFromOutlet(dMeanIn)
	return NewStageGeometry(stator, rotor)
}

func (gen *stageGeometryGenerator) GenerateFromStatorInlet(dMeanIn float64) StageGeometry {
	var stator = gen.statorGenerator.GenerateFromInlet(dMeanIn)
	var rotor = gen.rotorGenerator.GenerateFromInlet(stator.MeanProfile().Diameter(stator.XGapOut()))
	return NewStageGeometry(stator, rotor)
}

func (gen *stageGeometryGenerator) StatorGenerator() BladingGeometryGenerator {
	return gen.statorGenerator
}

func (gen *stageGeometryGenerator) RotorGenerator() BladingGeometryGenerator {
	return gen.rotorGenerator
}
