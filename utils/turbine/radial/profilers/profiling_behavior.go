package profilers

import (
	"github.com/Sovianum/turbocycle/impl/stage/states"
)

func NewStatorProfilingBehavior() ProfilingBehavior {
	return statorProfilingBehavior{}
}

func NewRotorProfilingBehavior() ProfilingBehavior {
	return rotorProfilingBehavior{}
}

type ProfilingBehavior interface {
	ProfilingVelocity(triangle states.VelocityTriangle) float64
	ProfilingAngle(triangle states.VelocityTriangle) float64
}

type statorProfilingBehavior struct{}

func (statorProfilingBehavior) ProfilingVelocity(triangle states.VelocityTriangle) float64 {
	return triangle.C()
}

func (statorProfilingBehavior) ProfilingAngle(triangle states.VelocityTriangle) float64 {
	return triangle.Alpha()
}

type rotorProfilingBehavior struct{}

func (rotorProfilingBehavior) ProfilingVelocity(triangle states.VelocityTriangle) float64 {
	return triangle.W()
}

func (rotorProfilingBehavior) ProfilingAngle(triangle states.VelocityTriangle) float64 {
	return triangle.Beta()
}
