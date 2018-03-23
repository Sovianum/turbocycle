package geometry

type StageGeometry interface {
	StatorGeometry() BladingGeometry
	RotorGeometry() BladingGeometry
}

func NewStageGeometry(statorGeometry, rotorGeometry BladingGeometry) StageGeometry {
	return &stageGeometry{statorGeometry: statorGeometry, rotorGeometry: rotorGeometry}
}

type stageGeometry struct {
	statorGeometry BladingGeometry
	rotorGeometry  BladingGeometry
}

func (geom *stageGeometry) StatorGeometry() BladingGeometry {
	return geom.statorGeometry
}

func (geom *stageGeometry) RotorGeometry() BladingGeometry {
	return geom.rotorGeometry
}
