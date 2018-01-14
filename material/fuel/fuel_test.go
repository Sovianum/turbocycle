package fuel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCh4_GetCombustionGas(t *testing.T) {
	var fuel = GetCH4()
	var gas1 = fuel.GetCombustionGas(1)
	assert.InDelta(t, 0, gas1.OxygenMassFraction(), 1e-7)
}
