package boundary

type Solution interface {
	Build() (xArr, yArr []float64)
}

func newBoundarySolution(xArr, yArr []float64) Solution {
	return &boundarySolution{
		xArr: xArr,
		yArr: yArr,
	}
}

type boundarySolution struct {
	xArr []float64
	yArr []float64
}

func (bs *boundarySolution) Build() (xArr, yArr []float64) {
	return bs.xArr, bs.yArr
}
