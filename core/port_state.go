package core

type PortState interface {
	Mix(another PortState, relaxCoef float64) (PortState, error)
	MaxResidual(another PortState) (float64, error)
}
