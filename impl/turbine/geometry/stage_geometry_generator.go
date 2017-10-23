package geometry

type StageGeometryGenerator interface {
	GenerateFromRotorInlet(dMeanIn float64) StageGeometry
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

func (gen *stageGeometryGenerator) StatorGenerator() BladingGeometryGenerator {
	return gen.statorGenerator
}

func (gen *stageGeometryGenerator) RotorGenerator() BladingGeometryGenerator {
	return gen.rotorGenerator
}
