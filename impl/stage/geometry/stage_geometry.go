package geometry

import "encoding/json"

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

func (geom *stageGeometry) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Stator BladingGeometry `json:"stator"`
		Rotor  BladingGeometry `json:"rotor"`
	}{
		Stator: geom.statorGeometry,
		Rotor:  geom.rotorGeometry,
	})
}

func (geom *stageGeometry) StatorGeometry() BladingGeometry {
	return geom.statorGeometry
}

func (geom *stageGeometry) RotorGeometry() BladingGeometry {
	return geom.rotorGeometry
}
