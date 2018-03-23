package laws

import (
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/stage/states"
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
		law.getCA(triangle0, hRel, lRel),
		law.getU(triangle0, hRel, lRel),
	)
}

func (law constantLabourLaw) getCA(triangle0 states.VelocityTriangle, hRel, lRel float64) float64 {
	var integralFunc = func(rRel float64) float64 {
		var rRelArray = common.LinSpace(1, rRel, 10)

		var sum float64 = 0
		for i := 0; i != len(rRelArray)-1; i++ {
			var rRelMean = (rRelArray[i] + rRelArray[i+1]) / 2
			var hRel = getHRel(rRel, lRel)
			var cu = law.getCU(triangle0, hRel, lRel)
			sum += cu * cu / rRelMean
		}

		var step = rRelArray[1] - rRelArray[0]
		return sum * step
	}

	var term1 = triangle0.CA() * triangle0.CA()
	var term2 = triangle0.CU() * triangle0.CU()

	var cu = law.getCU(triangle0, hRel, lRel)
	var term3 = -cu * cu

	var rRel = getRRel(hRel, lRel)
	var term4 = -2 * integralFunc(rRel)

	var result = math.Sqrt(term1 + term2 + term3 + term4)
	return result
}

func (law constantLabourLaw) getU(triangle0 states.VelocityTriangle, hRel, lRel float64) float64 {
	var rRel = getRRel(hRel, lRel)
	return triangle0.U() * rRel
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
