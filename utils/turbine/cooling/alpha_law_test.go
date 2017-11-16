package cooling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinedAlphaLaw(t *testing.T) {
	var law1 AlphaLaw = func(lengthParameter, theta float64) float64 {
		return 1
	}
	var law2 AlphaLaw = func(lengthParameter, theta float64) float64 {
		return 2
	}
	var law = JoinedAlphaLaw([]AlphaLaw{law1, law2}, []float64{0, 0.5, 1})
	assert.InDelta(t, 1, law(0.25, 100), 1e-8)
	assert.InDelta(t, 2, law(0.75, 100), 1e-8)
}
