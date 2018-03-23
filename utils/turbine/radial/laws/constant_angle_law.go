package laws

import (
	"math"

	"github.com/Sovianum/turbocycle/impl/stage/states"
)

func NewConstantAbsoluteAngleLaw() VelocityLaw {
	return constantAngleLaw{func(triangle states.VelocityTriangle) float64 {
		return triangle.Alpha()
	}}
}

func NewConstantRelativeAngleLaw() VelocityLaw {
	return constantAngleLaw{func(triangle states.VelocityTriangle) float64 {
		return triangle.Beta()
	}}
}

type constantAngleLaw struct {
	characteristicAngle func(triangle states.VelocityTriangle) float64
}

func (law constantAngleLaw) InletTriangle(triangle0 states.VelocityTriangle, hRel, lRel float64) states.VelocityTriangle {
	return states.NewInletTriangleFromProjections(law.getVelocities(triangle0, hRel, lRel))
}

func (law constantAngleLaw) OutletTriangle(triangle0 states.VelocityTriangle, hRel, lRel float64) states.VelocityTriangle {
	return states.NewOutletTriangleFromProjections(law.getVelocities(triangle0, hRel, lRel))
}

func (law constantAngleLaw) getVelocities(triangle0 states.VelocityTriangle, hRel, lRel float64) (cu, ca, u float64) {
	var rRel = getRRel(hRel, lRel)

	cu = law.getCU(triangle0, hRel, lRel)
	ca = law.getCA(triangle0, hRel, lRel)
	u = triangle0.U() * rRel
	return
}

func (law constantAngleLaw) getCU(triangle0 states.VelocityTriangle, hRel, lRel float64) float64 {
	var cuMean = triangle0.CU() // here you ALWAYS need absolute velocity projection
	var alpha = triangle0.Alpha()
	var rRel = getRRel(hRel, lRel)
	return cuMean / math.Pow(rRel, math.Pow(math.Cos(alpha), 2))
}

func (law constantAngleLaw) getCA(triangle0 states.VelocityTriangle, hRel, lRel float64) float64 {
	var velocityAMean = triangle0.CA()
	var angle = law.characteristicAngle(triangle0)
	var rRel = getRRel(hRel, lRel)

	return velocityAMean / math.Pow(rRel, math.Pow(math.Cos(angle), 2))
}
