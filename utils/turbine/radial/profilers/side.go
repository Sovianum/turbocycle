package profilers

import "github.com/Sovianum/turbocycle/impl/turbine/states"

func NewInletSide(profiler Profiler) Side {
	return inletSide{profiler: profiler}
}

func NewOutletSide(profiler Profiler) Side {
	return outletSide{profiler: profiler}
}

type Side interface {
	Triangle(hRel float64) states.VelocityTriangle
}

type inletSide struct {
	profiler Profiler
}

func (side inletSide) Triangle(hRel float64) states.VelocityTriangle {
	return side.profiler.InletTriangle(hRel)
}

type outletSide struct {
	profiler Profiler
}

func (side outletSide) Triangle(hRel float64) states.VelocityTriangle {
	return side.profiler.OutletTriangle(hRel)
}
