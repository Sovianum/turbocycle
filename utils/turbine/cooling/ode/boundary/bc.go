package boundary

// BC is a boundary condition for equation of type
// yArr'' + fArr(x, yArr)*yArr' + gArr(x, yArr) = 0
type BC interface {
	LHS0(step float64) float64
	LHS1(step float64) float64
	RHS(step float64) float64
}

func NewThirdTypeBC(alpha, y0 float64) BC {
	return &thirdTypeBC{
		alpha: alpha,
		y0:    y0,
	}
}

// bc of type yArr'(x_0) = alpha * (y_0 - y(x_0))
type thirdTypeBC struct {
	alpha float64
	y0    float64
}

func (bc *thirdTypeBC) LHS0(step float64) float64 {
	return 1 - step*bc.alpha
}

func (bc *thirdTypeBC) LHS1(step float64) float64 {
	return -1
}

func (bc *thirdTypeBC) RHS(step float64) float64 {
	return -step * bc.alpha * bc.y0
}

func NewSecondTypeBC(derivative float64) BC {
	return &secondTypeBC{derivative: derivative}
}

// bc of type yArr'(x_0) = yArr'_0
type secondTypeBC struct {
	derivative float64
}

func (bc *secondTypeBC) LHS0(step float64) float64 {
	return 1
}

func (bc *secondTypeBC) LHS1(step float64) float64 {
	return -1
}

func (bc *secondTypeBC) RHS(step float64) float64 {
	return -step * bc.derivative
}

func NewFirstTypeBC(value float64) BC {
	return &firstTypeBC{value: value}
}

// bc of type yArr(x_0) = y_0
type firstTypeBC struct {
	value float64
}

func (bc *firstTypeBC) LHS0(step float64) float64 { return 1 }

func (bc *firstTypeBC) LHS1(step float64) float64 { return 0 }

func (bc *firstTypeBC) RHS(step float64) float64 { return bc.value }
