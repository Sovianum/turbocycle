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
	GetNormEtaChar() constructive.TurbineCharFunc
}

func NewKazandjanTurbineCharacteristic() TurbineCharacteristic {
	return new(kazandjanTurbineCharacteristic)
}

type kazandjanTurbineCharacteristic struct{}

func (tc *kazandjanTurbineCharacteristic) GetNormMassRateChar() constructive.TurbineCharFunc {
	return func(_, normPiStag float64) float64 {
		x := 1 - normPiStag
		return 1 - 0.75*x*x*x
	}
}

func (tc *kazandjanTurbineCharacteristic) GetNormEtaChar() constructive.TurbineCharFunc {
	return func(_, normPiStag float64) float64 {
		x := 1 - normPiStag
		return 1 - x*x - 0.6*x*x*x
	}
}

func NewCIAMTurbineCharacteristic(etaT0, piT0, lambdaU0, stageNum float64) TurbineCharacteristic {
	return &ciamTurbineCharacteristic{
		etaT0:    etaT0,
		piT0:     piT0,
		lambdaU0: lambdaU0,
		stageNum: stageNum,
	}
}

type ciamTurbineCharacteristic struct {
	etaT0    float64
	piT0     float64
	lambdaU0 float64
	stageNum float64
}

func (tc *ciamTurbineCharacteristic) GetNormMassRateChar() constructive.TurbineCharFunc {
	return func(_, normPiStag float64) float64 {
		return tc.massRateNorm(normPiStag * tc.piT0)
	}
}

func (tc *ciamTurbineCharacteristic) GetNormEtaChar() constructive.TurbineCharFunc {
	return func(lambdaU, normPiStag float64) float64 {
		return tc.etaT(lambdaU, normPiStag*tc.piT0)
	}
}

func (tc *ciamTurbineCharacteristic) massRateNorm(piStag float64) float64 {
	x := tc.x(piStag)
	return panicNan(tc.normMassRateRel(x))
}

func (tc *ciamTurbineCharacteristic) x(piStag float64) float64 {
	return panicNan((piStag - 1) / (tc.piT0 - 1))
}

func (tc *ciamTurbineCharacteristic) normMassRateRel(x float64) float64 {
	if x < 1 {
		return 2*math.Sqrt(x) - x
	}
	return 1
}

func (tc *ciamTurbineCharacteristic) etaT(lambdaU, piStag float64) float64 {
	etaMax := tc.etaTMax(lambdaU)
	etaTNorm := tc.etaTNorm(lambdaU, piStag)
	return panicNan(etaMax * etaTNorm)
}

func (tc *ciamTurbineCharacteristic) etaTNorm(lambdaU, piStag float64) float64 {
	yStarNorm := tc.yStarNorm(lambdaU, piStag)
	if yStarNorm <= 1 {
		return panicNan(1 - (yStarNorm-1)*(yStarNorm-1))
	}
	return panicNan(1 - 0.55*(yStarNorm-1)*(yStarNorm-1))
}

func (tc *ciamTurbineCharacteristic) yStarNorm(lambdaU, piStag float64) float64 {
	yStar0 := tc.yStar0()
	yStarOptNorm := tc.yStarOptNorm(lambdaU)
	yStarOpt := yStarOptNorm * yStar0

	//yStar := tc.yStar(lambdaU, piStag) - by methodics
	yStar := yStar0 * lambdaU / tc.lambdaU0 * tc.lambdaAd(tc.piT0) / tc.lambdaAd(piStag)

	return panicNan(yStar / yStarOpt)
}

func (tc *ciamTurbineCharacteristic) yStar0() float64 {
	// this is incorrect value but i have nothing better
	return panicNan(
		tc.lambdaU0 / math.Sqrt(1-math.Pow(tc.piT0, (1-kGas)/kGas)),
	)
}

func (tc *ciamTurbineCharacteristic) yStarOptNorm(lambdaU float64) float64 {
	lambdaUNorm := tc.lambdaUNorm(lambdaU)
	if lambdaUNorm < 0.5 {
		lambdaUNorm = 0.5
	}
	return panicNan(
		0.91 + 0.36*(lambdaUNorm-0.5)*(lambdaUNorm-0.5),
	)
}

func (tc *ciamTurbineCharacteristic) yStar(lambdaU, piTStag float64) float64 {
	return panicNan(
		math.Sqrt(tc.stageNum) * lambdaU / tc.lambdaAd(piTStag),
	)
}

func (tc *ciamTurbineCharacteristic) lambdaAd(piTStag float64) float64 {
	return panicNan(
		math.Sqrt((kGas + 1) / (kGas - 1) * (1 - math.Pow(piTStag, -kGas/(kGas-1)))),
	)
}

func (tc *ciamTurbineCharacteristic) etaTMax(lambdaU float64) float64 {
	return panicNan(
		tc.etaT0 * tc.etaTMaxNorm(tc.lambdaUNorm(lambdaU)),
	)
}

func (tc *ciamTurbineCharacteristic) etaTMaxNorm(lambdaUNorm float64) float64 {
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

func (tc *ciamTurbineCharacteristic) lambdaUNorm(lambdaU float64) float64 {
	return panicNan(lambdaU / tc.lambdaU0)
}

func panicNan(f float64) float64 {
	if math.IsNaN(f) {
		panic("nan obtained")
	}
	return f
}
