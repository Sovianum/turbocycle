package boundary

import (
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/ode"
)

func newBoundarySolution(xArr, yArr []float64, step float64) ode.Solution {
	return &boundarySolution{
		xArr: xArr,
		yArr: yArr,
		step: step,
	}
}

type boundarySolution struct {
	xArr []float64
	yArr []float64
	step float64
}

func (bs *boundarySolution) Build() (xArr, yArr []float64) {
	return bs.xArr, bs.yArr
}

func (bs *boundarySolution) Step() float64 {
	return bs.step
}
