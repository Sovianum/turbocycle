package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAverage_Const(t *testing.T) {
	var f = func(float64) float64 { return 1 }
	var ave = Average(f, 0, 100, 100)

	assert.Equal(t, ave, float64(1))
}

func TestAverage_Linear(t *testing.T) {
	var f = func(x float64) float64 { return x }
	var ave = Average(f, 0, 100, 100)

	assert.Equal(t, ave, float64(50))
}
