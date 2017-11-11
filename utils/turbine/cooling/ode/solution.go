package ode

type Solution interface {
	Build() (xArr, yArr []float64)
}

type eulerSolution struct {
	derivative func(x, y float64) float64
	x0         float64
	y0         float64
	xMax       float64
	step       float64
}

func (e eulerSolution) Build() (xArr, yArr []float64) {
	var x, y = e.x0, e.y0
	for x <= e.xMax {
		xArr = append(xArr, x)
		yArr = append(yArr, y)
		x, y = e.nextPoint(x, y)
	}
	return xArr, yArr
}

func (e eulerSolution) nextPoint(xCurr, yCurr float64) (xNext, yNext float64) {
	xNext = xCurr + e.step
	yNext = yCurr + e.derivative(xCurr, yCurr)*e.step
	return
}
