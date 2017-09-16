package core

import "encoding/json"

type PortState interface {
	json.Marshaler
	Mix(another PortState, relaxCoef float64) (PortState, error)
	MaxResidual(another PortState) (float64, error)
}
