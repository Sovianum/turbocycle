package boundary

import (
	"fmt"
	"math"
	"testing"

	"github.com/gin-gonic/gin/json"
	"github.com/stretchr/testify/assert"
)

// testing on equation y'' + y' + 1 = 0; y(0) = 0; y(L) = 0
// precise solution y = (1 - e^(-x)) / (1 - e^(-L)) * L - x
func TestSolver1Type(t *testing.T) {
	length := 10.
	s := NewSolver(
		func(x float64) float64 { return 1 },
		func(x float64) float64 { return 0 },
		func(x float64) float64 { return -1 },
		NewFirstTypeBC(0),
		NewFirstTypeBC(0),
		0, length, 0.1,
	)
	solution, err := s.Solve()
	assert.Nil(t, err)

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
	s := NewSolver(
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
	s := NewSolver(
		func(x float64) float64 { return 1 },
		func(x float64) float64 { return 0 },
		func(x float64) float64 { return -1 },
		NewThirdTypeBC(-1, 1),
		NewSecondTypeBC(0),
		0, length, 0.001,
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

	b, _ := json.MarshalIndent(struct {
		X []float64 `json:"x"`
		Y []float64 `json:"y"`
	}{x, y}, "", "    ")
	fmt.Println(string(b))
}
