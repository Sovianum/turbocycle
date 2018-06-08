package boundary

import (
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/stretchr/testify/assert"
)

// testing on equation y'' + y' + 1 = 0; y(0) = 0; y(L) = 0
// precise solution y = (1 - e^(-x)) / (1 - e^(-L)) * L - x
func TestSolver1Type(t *testing.T) {
	length := 10.
	s := NewSolverFromFuncs(
		func(x float64) float64 { return 1 },
		func(x float64) float64 { return 0 },
		func(x float64) float64 { return -1 },
		NewFirstTypeBC(0),
		NewFirstTypeBC(0),
		0, length, 0.1,
	)
	solution, err := s.Solve()
	assert.Nil(t, err)
	if err != nil {
		return
	}

	precise := func(x float64) float64 {
		return (1-math.Exp(-x))/(1-math.Exp(-length))*length - x
	}

	x, y := solution.Build()

	assert.InDelta(t, 0, y[0], 1e-9)
	assert.InDelta(t, 0, y[len(y)-1], 1e-9)

	for i := range x {
		assert.InDelta(t, precise(x[i]), y[i], 1e-2, "%d", i)
	}
}

// testing on equation y'' + y' + 1 = 0; y'(0) = 0; y(L) = 0
// precise solution y = e^(-L) - e^(-x) + L - x
func TestSolver2Type(t *testing.T) {
	length := 10.
	s := NewSolverFromFuncs(
		func(x float64) float64 { return 1 },
		func(x float64) float64 { return 0 },
		func(x float64) float64 { return -1 },
		NewSecondTypeBC(0),
		NewFirstTypeBC(0),
		0, length, 0.01,
	)
	solution, err := s.Solve()
	assert.Nil(t, err)

	precise := func(x float64) float64 {
		return math.Exp(-length) - math.Exp(-x) + length - x
	}

	x, y := solution.Build()

	for i := range x {
		assert.InDelta(t, precise(x[i]), y[i], 1e-2, "%d", i)
	}
}

// testing on equation y'' + y' + 1 = 0; y'(0) = -1*(1 - y(0)); y'(L) = 0
// precise solution y = e^L*(2 - e^(-x)) - x
func TestSolver3Type(t *testing.T) {
	length := 1.
	s := NewSolverFromFuncs(
		func(x float64) float64 { return 1 },
		func(x float64) float64 { return 0 },
		func(x float64) float64 { return -1 },
		NewThirdTypeBC(-1, 1),
		NewSecondTypeBC(0),
		0, length, 0.01,
	)
	solution, err := s.Solve()
	assert.Nil(t, err)

	precise := func(x float64) float64 {
		return math.Exp(length)*(2-math.Exp(-x)) - x
	}

	x, y := solution.Build()

	for i := range x {
		assert.InDelta(t, precise(x[i]), y[i], 1e-1, "%d", i)
	}
}

// testing on equation y'' + y' + 1 = 0; y(0) = 0; y(L) = 0
// precise solution y = (1 - e^(-x)) / (1 - e^(-L)) * L - x
func TestSolverNonUniform(t *testing.T) {
	length := 1.
	precise := func(x float64) float64 {
		return (1-math.Exp(-x))/(1-math.Exp(-length))*length - x
	}
	xUniform := common.LinSpace(0, length, 1000)
	xNonUniform := applyUnit(xUniform, func(x float64) float64 { return x * x })

	fArr := make([]float64, len(xNonUniform))
	gArr := make([]float64, len(xNonUniform))
	hArr := make([]float64, len(xNonUniform))

	for i := range xNonUniform {
		fArr[i] = 1
		gArr[i] = 0
		hArr[i] = -1
	}

	s := NewSolverFromArrays(xNonUniform, fArr, gArr, hArr, NewFirstTypeBC(0), NewFirstTypeBC(0))
	solution, err := s.Solve()
	assert.Nil(t, err)
	if err != nil {
		return
	}

	x, y := solution.Build()

	assert.InDelta(t, 0, y[0], 1e-9)
	assert.InDelta(t, 0, y[len(y)-1], 1e-9)

	yPreciseArr := make([]float64, len(xNonUniform))
	for i, x := range xNonUniform {
		yPreciseArr[i] = precise(x)
	}

	//for i := range x {
	//	assert.InDelta(t, precise(x[i]), y[i], 1e-2, "%d", i)
	//}

	//data := struct {
	//	X        []float64 `json:"x"`
	//	Y        []float64 `json:"y"`
	//	YPrecise []float64 `json:"y_precise"`
	//}{
	//	X:        x,
	//	Y:        y,
	//	YPrecise: yPreciseArr,
	//}

	//b, _ := json.MarshalIndent(data, "", "    ")
	//fmt.Println(string(b))
}

func applyUnit(distribution []float64, f func(x float64) float64) []float64 {
	result := make([]float64, len(distribution))
	x0, xMax := distribution[0], distribution[len(distribution)-1]

	for i, x := range distribution {
		result[i] = f((x-x0)/(xMax-x0))*(xMax-x0) + x0
	}
	return result
}
