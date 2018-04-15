package boundary

import (
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/ode"
	"gonum.org/v1/gonum/mat"
)

func newBoundarySolution(yVec *mat.VecDense, step float64) ode.Solution {
	return &boundarySolution{
		yVec: yVec,
		step: step,
	}
}

type boundarySolution struct {
	yVec *mat.VecDense
	step float64
}

func (bs *boundarySolution) Build() (xArr, yArr []float64) {
	pointNum := bs.yVec.Len()
	x := make([]float64, pointNum)
	y := make([]float64, pointNum)

	for i := 0; i != pointNum; i++ {
		x[i] = bs.step * float64(i)
		y[i] = bs.yVec.At(i, 0)
	}

	return x, y
}

func (bs *boundarySolution) Step() float64 {
	return bs.step
}
