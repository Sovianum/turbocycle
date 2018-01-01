package solving

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

type EquationSystem interface {
	GetResiduals(x *mat.VecDense) (*mat.VecDense, error)
	Order() int
}

func NewEquationSystem(residualFunc func(*mat.VecDense) (*mat.VecDense, error), order int) EquationSystem {
	return &equationSystem{
		order:        order,
		residualFunc: residualFunc,
	}
}

type equationSystem struct {
	order        int
	residualFunc func(*mat.VecDense) (*mat.VecDense, error)
}

func (system *equationSystem) GetResiduals(x *mat.VecDense) (*mat.VecDense, error) {
	if xDim := x.Len(); xDim != system.order {
		return nil, fmt.Errorf("system order %d and argument dimension %d do not match", system.order, xDim)
	}

	var residuals, resErr = system.residualFunc(x)
	if resErr != nil {
		return nil, resErr
	}

	if resDim := residuals.Len(); resDim != system.order {
		return nil, fmt.Errorf("system order %d and result dimension %d do not match", system.order, resDim)
	}

	return residuals, nil
}

func (system *equationSystem) Order() int {
	return system.order
}
