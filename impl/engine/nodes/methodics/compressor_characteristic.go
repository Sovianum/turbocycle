package methodics

import (
	"math"

	"github.com/Sovianum/turbocycle/common"
	math2 "github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"gonum.org/v1/gonum/mat"
)

type CompressorCharGen interface {
	GetNormRPMChar() constructive.CompressorCharFunc
	GetNormEtaChar() constructive.CompressorCharFunc
}

func NewCompressorCharGen(
	piC0, etaStag0, massRateNorm0, precision, relaxCoef float64, iterLimit int,
) CompressorCharGen {
	return &compressorCharGen{
		piC0:          piC0,
		etaStag0:      etaStag0,
		massRateNorm0: massRateNorm0,

		rpmNormCurr: 1,
		phiCurr:     math.Pi / 4,

		precision: precision,
		relaxCoef: relaxCoef,
		iterLimit: iterLimit,
	}
}

type compressorCharGen struct {
	piC0          float64
	etaStag0      float64
	massRateNorm0 float64

	rpmNormCurr float64
	phiCurr     float64

	precision float64
	relaxCoef float64
	iterLimit int
}

func (ccg *compressorCharGen) GetNormRPMChar() constructive.CompressorCharFunc {
	return func(normMassRate, normPiStag float64) float64 {
		if err := ccg.resolveCoordinates(normMassRate, normPiStag); err != nil {
			panic(err)
		}
		return ccg.rpmNormCurr
	}
}

func (ccg *compressorCharGen) GetNormEtaChar() constructive.CompressorCharFunc {
	return func(normMassRate, normPiStag float64) float64 {
		if err := ccg.resolveCoordinates(normMassRate, normPiStag); err != nil {
			panic(err)
		}
		eta := ccg.etaStag(ccg.phiCurr, ccg.rpmNormCurr)
		return eta / ccg.etaStag0
	}
}

func (ccg *compressorCharGen) resolveCoordinates(normMassRate, normPiStag float64) error {
	sys := ccg.getEqSystem(normMassRate, normPiStag)
	solver, err := newton.NewUniformNewtonSolver(sys, 1e-3, newton.NoLog)

	solution, err := solver.Solve(mat.NewVecDense(2, []float64{
		ccg.phiCurr, ccg.rpmNormCurr,
	}), ccg.precision, ccg.relaxCoef, ccg.iterLimit)
	if err != nil {
		return err
	}
	ccg.phiCurr, ccg.rpmNormCurr = solution.At(0, 0), solution.At(1, 0)
	return nil
}

func (ccg *compressorCharGen) getEqSystem(normMassRate, normPiC float64) math2.EquationSystem {
	return math2.NewEquationSystem(func(input *mat.VecDense) (*mat.VecDense, error) {
		phi, rpmNorm := input.At(0, 0), input.At(1, 0)
		newMassRateNorm := ccg.massRateNorm(phi, rpmNorm) / ccg.massRateNorm0
		newPiCStagNorm := ccg.piCStag(phi, rpmNorm) / ccg.piC0
		return mat.NewVecDense(2, []float64{
			newMassRateNorm - normMassRate,
			newPiCStagNorm - normPiC,
		}), nil
	}, 2)
}

func (ccg *compressorCharGen) piCStag(phi, rpmNorm float64) float64 {
	return 1 + (ccg.piCOpt(rpmNorm)-1)*ccg.radius(phi, rpmNorm)*math.Sin(phi)
}

func (ccg *compressorCharGen) massRateNorm(phi, rpmNorm float64) float64 {
	return ccg.massRateNormOpt(rpmNorm) * ccg.radius(phi, rpmNorm) * math.Cos(phi)
}

func (ccg *compressorCharGen) etaStag(phi, rpmNorm float64) float64 {
	a1 := 2.4*(rpmNorm-1) + 0.2*(ccg.piC0-3) + 1.5
	if a1 < 0.2 {
		a1 = 0.2
	}

	// removed this factor cos made solution extremely unstable
	//c1 := 1.33319
	//if a1 > 0.5 {
	//	c1 = common.Sum([]float64{
	//		-7.86223, 0.203704 * ccg.piC0,
	//		34.851 * rpmNorm,
	//		1.55811e-3 * ccg.piC0 * ccg.piC0,
	//		-0.307129 * ccg.piC0 * rpmNorm,
	//		-43.3512 * rpmNorm * rpmNorm,
	//		17.8252 * rpmNorm * rpmNorm * rpmNorm,
	//	})
	//}
	//b1 := 1.
	//if phi > math.Pi/4 {
	//	b1 = 1 - math.Pow(phi-math.Pi/4, 1.2*c1)
	//}
	b1 := 1.

	factor := b1 * math.Pow(math.Sin(2*phi), a1)
	return ccg.etaStagOpt(rpmNorm) * factor
}

func (ccg *compressorCharGen) radius(phi, rpmNorm float64) float64 {
	qMNorm := 1 + 1.85/math.Exp(4.8*rpmNorm) - 2.5e-3*ccg.piC0
	b2 := 1.8 - 8.167*rpmNorm + 18.334*rpmNorm*rpmNorm
	return qMNorm/math.Cos(phi) - (qMNorm-1)*math.Sqrt2*math.Pow(math.Tan(phi), b2)
}

func (ccg *compressorCharGen) piCOpt(rpmNorm float64) float64 {
	k := 1.4
	etaFactor := ccg.etaStagOpt(rpmNorm) / ccg.etaStag0
	piFactor := math.Pow(ccg.piC0, (k-1)/k) - 1
	return math.Pow(
		1+piFactor*etaFactor*rpmNorm*rpmNorm,
		k/(k-1),
	)
}

func (ccg *compressorCharGen) etaStagOpt(rpmNorm float64) float64 {
	// a, b, c - values from original method
	//a := 1.1018 - 7.74e-2*(ccg.piC0-1)
	b := -0.5652 + 0.1486*(ccg.piC0-1)
	c := 1.366 - 7.8e-2*(ccg.piC0-1)
	// a1, b1, c1 - recalculated values for another shape of the function
	// original value of a1 was a1 = a + b + c, but it was changed so that
	// (a1+b1*(rpmNorm-1)+c1*(rpmNorm-1)*(rpmNorm-1))(rpmNorm = 0) = Pi / 2
	a1 := math.Pi / 2
	b1 := b + 2*c
	c1 := c
	return ccg.etaStag0 * math.Sin(a1+b1*(rpmNorm-1)+c1*(rpmNorm-1)*(rpmNorm-1))
}

func (ccg *compressorCharGen) massRateNormOpt(rpmNorm float64) float64 {
	a := math.Pow(ccg.piCOpt(rpmNorm)/ccg.piC0, 0.8571)
	b := ccg.massRateNorm0
	c := 1 + ccg.qRel(rpmNorm)
	return a * b * c
}

func (ccg *compressorCharGen) qRel(rpmNorm float64) float64 {
	if rpmNorm >= 1 {
		return 0
	}
	terms := []float64{
		-2.0696,
		0.319178 * ccg.piC0,
		2.92996 * rpmNorm,
		-1.53413e-2 * ccg.piC0 * ccg.piC0,
		-0.104443 * ccg.piC0 * rpmNorm,
		-1.42566 * rpmNorm * rpmNorm,
	}
	return (1 - rpmNorm) * common.Sum(terms)
}

func (ccg *compressorCharGen) piRel(rpmNorm float64) float64 {
	if rpmNorm >= 1 || ccg.piC0 <= 3 {
		return 1
	}
	a := (1 - rpmNorm) * (1 - rpmNorm)
	b := math.Pow(ccg.piC0-3, 0.5)
	c := 0.177*ccg.piC0*rpmNorm - 0.457
	return 1 - a*b*c
}
