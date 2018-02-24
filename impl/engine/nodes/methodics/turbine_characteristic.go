package methodics

import (
	"math"

	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
)

const (
	kGas = 1.33 // TODO make smth more accurate
)

type TurbineCharacteristic interface {
	GetNormMassRateChar() constructive.TurbineCharFunc
	GetNormMassRateCharConst() constructive.TurbineCharFunc
	GetNormEtaChar() constructive.TurbineCharFunc
	GetNormEtaCharConst() constructive.TurbineCharFunc
}

func NewTurbineCharacteristic(etaT0, piT0, lambdaU0, stageNum float64) TurbineCharacteristic {
	return &turbineCharacteristic{
		etaT0:    etaT0,
		piT0:     piT0,
		lambdaU0: lambdaU0,
		stageNum: stageNum,
	}
}

type turbineCharacteristic struct {
	etaT0    float64
	piT0     float64
	lambdaU0 float64
	stageNum float64
}

func (tc *turbineCharacteristic) GetNormMassRateChar() constructive.TurbineCharFunc {
	return func(lambdaU, normPiStag float64) float64 {
		return tc.massRateNorm(normPiStag * tc.piT0)
	}
}

func (tc *turbineCharacteristic) GetNormMassRateCharConst() constructive.TurbineCharFunc {
	return func(lambdaU, normPiStag float64) float64 {
		return 1
	}
}

func (tc *turbineCharacteristic) GetNormEtaChar() constructive.TurbineCharFunc {
	return func(lambdaU, normPiStag float64) float64 {
		return tc.etaT(lambdaU, normPiStag*tc.piT0)
	}
}

func (tc *turbineCharacteristic) GetNormEtaCharConst() constructive.TurbineCharFunc {
	return func(lambdaU, normPiStag float64) float64 {
		return 1
	}
}

func (tc *turbineCharacteristic) massRateNorm(piStag float64) float64 {
	return panicNan(tc.normMassRateRel(tc.x(piStag)))
}

func (tc *turbineCharacteristic) x(piStag float64) float64 {
	return panicNan((piStag - 1) / (tc.piT0 - 1))
}

func (tc *turbineCharacteristic) normMassRateRel(x float64) float64 {
	if x > 1 {
		return 1
	} else if 0.3 < x && x <= 1 {
		return panicNan(math.Pow(1-(1-x)*(1-x), 1/4))
	} else if 0.1 < x && x <= 0.3 {
		return panicNan(1.26 * math.Pow(x, 1/3))
	}
	return panicNan(1.85 * math.Sqrt(x))
}

func (tc *turbineCharacteristic) etaT(lambdaU, piStag float64) float64 {
	etaMax := tc.etaTMax(lambdaU)
	etaTNorm := tc.etaTNorm(lambdaU, piStag)
	return panicNan(etaMax * etaTNorm)
}

func (tc *turbineCharacteristic) etaTNorm(lambdaU, piStag float64) float64 {
	yStarNorm := tc.yStarNorm(lambdaU, piStag)
	if yStarNorm <= 1 {
		return panicNan(1 - (yStarNorm-1)*(yStarNorm-1))
	}
	return panicNan(1 - 0.55*(yStarNorm-1)*(yStarNorm-1))
}

func (tc *turbineCharacteristic) yStarNorm(lambdaU, piStag float64) float64 {
	yStar0 := tc.yStar0()
	yStarOptNorm := tc.yStarOptNorm(lambdaU)
	yStarOpt := yStarOptNorm * yStar0

	//yStar := tc.yStar(lambdaU, piStag) - by methodics
	yStar := yStar0 * lambdaU / tc.lambdaU0 * tc.lambdaAd(tc.piT0) / tc.lambdaAd(piStag)

	return panicNan(yStar / yStarOpt)
}

func (tc *turbineCharacteristic) yStar0() float64 {
	// this is incorrect value but i have nothing better
	return panicNan(
		tc.lambdaU0 / math.Sqrt(1-math.Pow(tc.piT0, (1-kGas)/kGas)),
	)
}

func (tc *turbineCharacteristic) yStarOptNorm(lambdaU float64) float64 {
	lambdaUNorm := tc.lambdaUNorm(lambdaU)
	if lambdaUNorm < 0.5 {
		lambdaUNorm = 0.5
	}
	return panicNan(
		0.91 + 0.36*(lambdaUNorm-0.5)*(lambdaUNorm-0.5),
	)
}

func (tc *turbineCharacteristic) yStar(lambdaU, piTStag float64) float64 {
	return panicNan(
		math.Sqrt(tc.stageNum) * lambdaU / tc.lambdaAd(piTStag),
	)
}

func (tc *turbineCharacteristic) lambdaAd(piTStag float64) float64 {
	return panicNan(
		math.Sqrt((kGas + 1) / (kGas - 1) * (1 - math.Pow(piTStag, -kGas/(kGas-1)))),
	)
}

func (tc *turbineCharacteristic) etaTMax(lambdaU float64) float64 {
	return panicNan(
		tc.etaT0 * tc.etaTMaxNorm(tc.lambdaUNorm(lambdaU)),
	)
}

func (tc *turbineCharacteristic) etaTMaxNorm(lambdaUNorm float64) float64 {
	if lambdaUNorm >= 1 {
		return panicNan(1 + 0.03*(lambdaUNorm-1))
	} else if lambdaUNorm > 0.5 {
		return panicNan(
			1 - (0.1+0.18*(math.Pow(tc.piT0, 0.4)-1.44))*(1-lambdaUNorm),
		)
	}
	return panicNan(
		1 - 0.5*(0.1+0.18*(math.Pow(tc.piT0, 0.4)-1.44)),
	)
}

func (tc *turbineCharacteristic) lambdaUNorm(lambdaU float64) float64 {
	return panicNan(lambdaU / tc.lambdaU0)
}

func panicNan(f float64) float64 {
	if math.IsNaN(f) {
		panic("nan obtained")
	}
	return f
}
