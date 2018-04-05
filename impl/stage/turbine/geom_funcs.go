package turbine

import (
	"math"

	"github.com/Sovianum/turbocycle/impl/stage/geometry"
)

// xRel = axialDistance / bladeLengthOnStart
func RecalculateLRel(
	lRel0, xRel, gammaIn, gammaOut float64,
) float64 {
	_, gammaM := geometry.GetTotalAndMeanLineAngles(gammaIn, gammaOut, MidLineFactor)
	factor := (1 + xRel*(math.Tan(gammaOut)-math.Tan(gammaIn))) / (1 + 2*xRel*lRel0*math.Tan(gammaM))
	return lRel0 * factor
}
