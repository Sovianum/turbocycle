package laws

import "github.com/Sovianum/turbocycle/impl/stage/states"

type VelocityLaw interface {
	InletVelocityLaw
	OutletVelocityLaw
}

type OutletVelocityLaw interface {
	OutletTriangle(triangle0 states.VelocityTriangle, hRel, lRel float64) states.VelocityTriangle
}

type InletVelocityLaw interface {
	InletTriangle(triangle0 states.VelocityTriangle, hRel, lRel float64) states.VelocityTriangle
}
