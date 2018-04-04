package ditributions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBiParabolic(t *testing.T) {
	start := 1.
	end := 4.
	opt := 2.
	startLoss := 0.2
	endLoss := 0.1

	d := GetUnitBiParabolic(start, end, opt, startLoss, endLoss)
	assert.InDelta(t, 0.8, d(start), 1e-9)
	assert.InDelta(t, 0.9, d(end), 1e-9)
	assert.InDelta(t, 1, d(opt), 1e-9)
	assert.True(t, d(opt-0.01) < d(opt) && d(opt) > d(opt+0.01))
}
