package opt

import (
	"fmt"

	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"gonum.org/v1/gonum/mat"
)

type Optimizer1D interface {
	Maximize(x0 *mat.VecDense, precision float64, relaxCoef float64, iterLimit int) (*mat.VecDense, error)
	Minimize(x0 *mat.VecDense, precision float64, relaxCoef float64, iterLimit int) (*mat.VecDense, error)
}

func NewOptimizer(f func(dense *mat.VecDense) (float64, error), step float64, logFunc newton.LogFunc) Optimizer1D {
	return &optimizer1D{
		targetFunc: f,
		step:       step,
		logFunc:    logFunc,
	}
}

type optimizer1D struct {
	targetFunc func(x *mat.VecDense) (float64, error)
	step       float64
	logFunc    newton.LogFunc
}

func (opt *optimizer1D) Minimize(x0 *mat.VecDense, precision float64, relaxCoef float64, iterLimit int) (*mat.VecDense, error) {
	copy := &optimizer1D{
		targetFunc: func(x *mat.VecDense) (float64, error) {
			y, err := opt.targetFunc(x)
			if err != nil {
				return 0, err
			}
			return -y, nil
		},
		step:    opt.step,
		logFunc: opt.logFunc,
	}
	return copy.Maximize(x0, precision, relaxCoef, iterLimit)
}

func (opt *optimizer1D) Maximize(x0 *mat.VecDense, precision float64, relaxCoef float64, iterLimit int) (*mat.VecDense, error) {
	grad := mat.NewVecDense(x0.Len(), nil)
	x := x0
	var err error
	for i := 0; i != iterLimit; i++ {
		grad, err = getGradient(opt.targetFunc, x, opt.step)
		if err != nil {
			return nil, err
		}
		opt.logFunc(i, precision, grad)
		if mat.Norm(grad, 2) < precision {
			return x, nil
		}
		x.AddScaledVec(x, relaxCoef, grad)
	}
	return nil, fmt.Errorf("failed to converge grad(f) = %f", mat.Norm(grad, 2))
}

func getGradient(f func(*mat.VecDense) (float64, error), x *mat.VecDense, step float64) (*mat.VecDense, error) {
	y, err := f(x)
	if err != nil {
		return nil, err
	}
	grad := mat.NewVecDense(x.Len(), nil)
	xNew := mat.NewVecDense(x.Len(), nil)

	for i := 0; i != x.Len(); i++ {
		xNew.CloneVec(x)
		xNew.SetVec(i, xNew.At(i, 0)+step)

		yNew, err := f(xNew)
		if err != nil {
			return nil, err
		}
		grad.SetVec(i, (yNew-y)/step)
	}
	return grad, nil
}
