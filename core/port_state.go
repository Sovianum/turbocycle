package core

type IPortState interface {
	Mix(another IPortState, relaxCoef float64) (IPortState, error)
	MaxResidual(another IPortState) (float64, error)
}
