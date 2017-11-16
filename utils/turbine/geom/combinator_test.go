package geom

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

func TestCombinator(t *testing.T) {
	var vec = mat.NewVecDense(2, []float64{0, 1})
	var combinator = NewCombinator()

	var expected = mat.NewVecDense(2, []float64{-1, -1})
	var got = combinator.
		Combine(
			Translation(mat.NewVecDense(2, []float64{0, 1})),
		).
		Combine(
			Reflection(math.Pi / 2),
		).
		Combine(
			Rotation(-math.Pi / 2),
		).
		Build().Apply(vec)

	assert.True(
		t,
		mat.EqualApprox(got, expected, precision),
		testMessage(0, expected, got),
	)
}
