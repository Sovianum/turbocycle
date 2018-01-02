package variator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariatorFromCallables(t *testing.T) {
	var x = 1.
	var getter = func() float64 { return x }
	var setter = func(f float64) { x = f }
	var variator = VariatorFromCallables(getter, setter)

	var val = 2.
	variator.SetValue(val)

	assert.InDelta(t, val, variator.GetValue(), 1e-7)
}

func TestVariatorFromPointer(t *testing.T) {
	var x = 1.
	var variator = VariatorFromPointer(&x)

	var val = 2.
	variator.SetValue(val)

	assert.InDelta(t, val, variator.GetValue(), 1e-7)
}

func TestAdd(t *testing.T) {
	var x = 1.
	var variator = VariatorFromPointer(&x)

	Add(variator, 1)

	assert.InDelta(t, 2., x, 1e-7)
}
