package common

import (
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

type Func1D func(x float64) float64

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
