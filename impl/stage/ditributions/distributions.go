package ditributions

import "github.com/Sovianum/turbocycle/impl/stage/common"

func GetUnitConstant() common.Func1D {
	return func(x float64) float64 {
		return 1
	}
}

func GetUnitBiParabolic(start, end, opt, startLossRate, endLossRate float64) common.Func1D {
	shape := func(xNorm, y float64) float64 {
		return -(1-y)*xNorm*xNorm + 2*(1-y)*xNorm + y
	}
	return func(x float64) float64 {
		xNorm1 := (x - start) / (opt - start)
		xNorm2 := (end - x) / (end - opt)

		if x <= opt {
			return shape(xNorm1, 1-startLossRate)
		} else {
			return shape(xNorm2, 1-endLossRate)
		}
	}
}
