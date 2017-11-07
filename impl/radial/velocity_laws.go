package radial

import (
	"math"

	"github.com/Sovianum/turbocycle/impl/turbine/states"
)

type VelocityLaw interface {
	InletTriangle(triangle0 states.VelocityTriangle, hRel, lRel float64) states.VelocityTriangle
	OutletTriangle(triangle0 states.VelocityTriangle, hRel, lRel float64) states.VelocityTriangle
}

func NewPowerVelocityLaw(powerCoef float64) VelocityLaw {
	return powerVelocityLaw{powerCoef:powerCoef}
}

func NewConstantCirculationVelocityLaw() VelocityLaw {
	return powerVelocityLaw{1}
}

func NewSolidBodyLaw() VelocityLaw {
	return powerVelocityLaw{-1}
}

type powerVelocityLaw struct {
	powerCoef float64
}

func (law powerVelocityLaw) InletTriangle(triangle0 states.VelocityTriangle, hRel, lRel float64) states.VelocityTriangle {
	return states.NewInletTriangleFromProjections(law.getVelocities(triangle0, hRel, lRel))
}

func (law powerVelocityLaw) OutletTriangle(triangle0 states.VelocityTriangle, hRel, lRel float64) states.VelocityTriangle {
	return states.NewOutletTriangleFromProjections(law.getVelocities(triangle0, hRel, lRel))
}

func (law powerVelocityLaw) getVelocities(triangle0 states.VelocityTriangle, hRel, lRel float64) (cu, ca, u float64) {
	var rRel = getRRel(hRel, lRel)

	cu = law.getCU(triangle0, hRel, lRel)
	ca = law.getCA(triangle0, hRel, lRel)
	u = triangle0.U() * rRel
	return 
}

func (law powerVelocityLaw) getCU(triangle0 states.VelocityTriangle, hRel, lRel float64) float64 {
	var cuMean = triangle0.CU()
	var rRel = getRRel(hRel, lRel)
	return cuMean / math.Pow(rRel, law.powerCoef)
}

func (law powerVelocityLaw) getCA(triangle0 states.VelocityTriangle, hRel, lRel float64) float64 {
	var cuMean = triangle0.CU()
	var caMean = triangle0.CA()
	var rRel = getRRel(hRel, lRel)

	var term1 = caMean * caMean
	var term2 = cuMean * cuMean * (1 - math.Pow(rRel, -2)) * (1 - 1 / law.powerCoef)

	return math.Sqrt(term1 + term2)
}

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
	var cuMean = triangle0.CU()		// here you ALWAYS need absolute velocity projection
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

func getRRel(hRel, lRel float64) float64 {
	return hRel * (1 + lRel) + (1 - hRel) * (1 - lRel)
}
