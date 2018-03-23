package laws

import (
	"math"

	"github.com/Sovianum/turbocycle/impl/stage/states"
)

func NewPowerVelocityLaw(powerCoef float64) VelocityLaw {
	return powerVelocityLaw{powerCoef: powerCoef}
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
	var term2 = cuMean * cuMean * (1 - math.Pow(rRel, -2)) * (1 - 1/law.powerCoef)

	return math.Sqrt(term1 + term2)
}
