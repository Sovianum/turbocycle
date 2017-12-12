package cooling

import (
	"fmt"
	"math"

	"github.com/Sovianum/turbocycle/material/gases"
)

type AlphaLaw func(lengthParameter, theta float64) float64

func JoinedAlphaLaw(alphaLaws []AlphaLaw, boundaryPoints []float64) AlphaLaw {
	if len(alphaLaws) != len(boundaryPoints)-1 {
		panic(fmt.Errorf(
			"invalid length: len(alpha_law) = %d, len(boundaryPoints) = %d",
			len(alphaLaws),
			len(boundaryPoints),
		))
	}
	return func(lengthParameter, theta float64) float64 {
		for i, innerLengthParameter := range boundaryPoints[1:] {
			if lengthParameter <= innerLengthParameter {
				var localLengthParameter = lengthParameter - boundaryPoints[i]
				return alphaLaws[i](localLengthParameter, theta)
			}
		}
		return alphaLaws[len(alphaLaws)-1](lengthParameter-boundaryPoints[len(boundaryPoints)-1], theta)
		panic(fmt.Errorf(
			"out of range: t = %f, t_min = %f, t_max = %f",
			lengthParameter,
			boundaryPoints[0],
			boundaryPoints[len(boundaryPoints)-1],
		))
	}
}

func ConstantAlphaLaw(alpha float64) AlphaLaw {
	return func(lengthParameter, theta float64) float64 {
		return alpha
	}
}

func DefaultAirAlphaLaw(gas gases.Gas, bladeLength, gapWidth, massRate float64) AlphaLaw {
	return func(lengthParameter, theta float64) float64 {
		var factor1 = gas.Lambda(theta) / (2 * gapWidth)
		var factor2 = math.Pow(massRate/(bladeLength*gas.Mu(theta)), 0.8)
		var result = 0.02 * factor1 * factor2
		return result
	}
}

func CylinderAlphaLaw(gas gases.Gas, massRateIntensity float64, diameter float64) AlphaLaw {
	return func(t, theta float64) float64 {
		var factor1 = 0.74

		var factor2 = gas.Lambda(theta) / diameter
		var factor3 = math.Sqrt(massRateIntensity * diameter / gas.Mu(theta))
		return factor1 * factor2 * factor3
	}
}

func PSAlphaLaw(meanAlpha float64) AlphaLaw {
	return func(t, theta float64) float64 {
		return meanAlpha
	}
}

func InletSSAlphaLaw(meanAlpha float64) AlphaLaw {
	return func(t, theta float64) float64 {
		return 0.6 * meanAlpha
	}
}

func OutletSSAlphaLaw(meanAlpha float64) AlphaLaw {
	return func(t, theta float64) float64 {
		return 1.5 * meanAlpha
	}
}
