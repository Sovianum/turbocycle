package geometry

type StageGeometryGenerator interface {
	GenerateFromRotorInlet(dMeanIn float64) StageGeometry
	GenerateFromStatorInlet(dMeanIn float64) StageGeometry
	StatorGenerator() BladingGeometryGenerator
	RotorGenerator() BladingGeometryGenerator
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
