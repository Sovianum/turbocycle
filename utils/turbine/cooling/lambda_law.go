package cooling

import (
	"math"
)

type LambdaLaw func(xRel float64) float64

func PSLambdaLaw(lambdaIn, lambdaOut float64) LambdaLaw {
	return powerLambdaLaw(lambdaIn, lambdaOut, 2)
}

func SSLambdaLaw(lambdaIn, lambdaOut float64) LambdaLaw {
	return powerLambdaLaw(lambdaIn, lambdaOut, 0.25)
}

func powerLambdaLaw(lambdaIn, lambdaOut, power float64) LambdaLaw {
	return func(xRel float64) float64 {
		return math.Pow(
			1+(math.Pow(lambdaOut/lambdaIn, 1/power)-1)*xRel,
			power,
		) * lambdaIn
	}
}
