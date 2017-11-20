package laws

import (
	"github.com/Sovianum/turbocycle/impl/turbine/states"
)

func NewConstantLabourLaw(law InletVelocityLaw, inletTriangle states.VelocityTriangle) OutletVelocityLaw {
	return constantLabourLaw{
		inletLaw:          law,
		inletMeanTriangle: inletTriangle,
	}
}

type constantLabourLaw struct {
	inletLaw          InletVelocityLaw
	inletMeanTriangle states.VelocityTriangle
}

func (law constantLabourLaw) OutletTriangle(triangle0 states.VelocityTriangle, hRel, lRel float64) states.VelocityTriangle {
	return states.NewOutletTriangleFromProjections(
		law.getCU(triangle0, hRel, lRel),
		law.getCA(hRel, lRel),
		law.getU(triangle0, hRel, lRel),
	)
}

func (law constantLabourLaw) getU(triangle0 states.VelocityTriangle, hRel, lRel float64) float64 {
	var rRel = getRRel(hRel, lRel)
	return triangle0.U() * rRel
}

func (law constantLabourLaw) getCA(hRel, lRel float64) float64 {
	return law.inletLaw.InletTriangle(law.inletMeanTriangle, hRel, lRel).CA()
}

func (law constantLabourLaw) getCU(triangle0 states.VelocityTriangle, hRel, lRel float64) float64 {
	var lu = law.meanLabour(triangle0)
	var u = law.inletMeanTriangle.U()
	var cu = law.inletMeanTriangle.CU()
	return lu/u - cu
}

func (law constantLabourLaw) meanLabour(triangle0 states.VelocityTriangle) float64 {
	var term1 = law.inletMeanTriangle.U() * law.inletMeanTriangle.CU()
	var term2 = triangle0.U() * triangle0.CU()
	return term1 + term2
}
