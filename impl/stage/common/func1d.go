package common

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"gonum.org/v1/gonum/mat"
)

func GetEqSys1D(sysCall variator.SysCall, funcSetter func(Func1D), fg FuncGen1D) math.EquationSystem {
	return math.NewEquationSystem(func(vec *mat.VecDense) (*mat.VecDense, error) {
		x := vec.At(0, 0)
		funcSetter(fg(x))
		return sysCall()
	}, 1)
}

type FuncGen1D func(p float64) Func1D

func NewDistributionVariator(gen FuncGen1D, funcSetter func(Func1D)) variator.Variator {
	var localF float64
	setter := func(f float64) {
		localF = f
		funcSetter(gen(f))
	}
	getter := func() float64 {
		return localF
	}
	return variator.FromCallables(getter, setter)
}

func Scaler(base Func1D) FuncGen1D {
	return func(p float64) Func1D {
		return base.Scale(p)
	}
}

func DisplacerX(base Func1D) FuncGen1D {
	return func(p float64) Func1D {
		return base.MoveX(p)
	}
}

func DisplacerY(base Func1D) FuncGen1D {
	return func(p float64) Func1D {
		return base.MoveY(p)
	}
}

func Func1DFromConst(c float64) Func1D {
	return func(x float64) float64 {
		return c
	}
}

type Func1D func(x float64) float64

func (f Func1D) GetUnitNormalizedSamples(xs []float64) []float64 {
	return f.GetNormalizedSamples(xs, 1)
}

func (f Func1D) GetNormalizedSamples(xs []float64, targetSum float64) []float64 {
	samples := f.GetSamples(xs)
	NormalizeSamples(samples, targetSum)
	return samples
}

func (f Func1D) GetSamples(xs []float64) []float64 {
	result := make([]float64, len(xs))
	for i, x := range xs {
		result[i] = f(x)
	}
	return result
}

func NormalizeSamples(samples []float64, targetSum float64) {
	sum := common.Sum(samples)
	factor := targetSum / sum
	for i := range samples {
		samples[i] *= factor
	}
}

func (f Func1D) Then(another Func1D) Func1D {
	return func(x float64) float64 {
		return another(f(x))
	}
}

func (f Func1D) Scale(factor float64) Func1D {
	return func(x float64) float64 {
		return factor * f(x)
	}
}

func (f Func1D) MoveX(offset float64) Func1D {
	return func(x float64) float64 {
		return f(x - offset)
	}
}

func (f Func1D) MoveY(offset float64) Func1D {
	return func(x float64) float64 {
		return f(x) + offset
	}
}
