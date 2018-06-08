package opt

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

func TestOptimize(t *testing.T) {
	f := func(v *mat.VecDense) (float64, error) {
		x := v.At(0, 0) - 1
		return -x * x, nil
	}
	opt := NewOptimizer(f, 1e-5, newton.NoLog)
	x, err := opt.Maximize(mat.NewVecDense(1, nil), 1e-5, 0.1, 1000)
	assert.Nil(t, err)
	assert.InDelta(t, 1, x.At(0, 0), 1e-5)
}
