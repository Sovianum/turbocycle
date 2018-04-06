package ditributions

import "github.com/Sovianum/turbocycle/impl/stage/common"

func GetLinear(x0, y0, x1, y1 float64) common.Func1D {
	k := (y1 - y0) / (x1 - x0)
	b := (y0*x1 - x0*y1) / (x1 - x0)
	return func(x float64) float64 {
		return k*x + b
	}
}

func GetUnitConstant() common.Func1D {
	return GetConstant(1)
}

func GetConstant(c float64) common.Func1D {
	return func(x float64) float64 {
		return c
	}
}

func GetUnitBiParabolic(start, end, opt, startLossRate, endLossRate float64) common.Func1D {
	shape := func(xNorm, y float64) float64 {
		return -(1-y)*xNorm*xNorm + 2*(1-y)*xNorm + y
	}
	return func(x float64) float64 {
		xNorm1 := (x - start) / (opt - start)
		if opt == start {
			xNorm1 = 1
		}

		xNorm2 := (end - x) / (end - opt)
		if end == opt {
			xNorm2 = 1
		}

		if x <= opt {
			return shape(xNorm1, 1-startLossRate)
		} else {
			return shape(xNorm2, 1-endLossRate)
		}
	}
}
