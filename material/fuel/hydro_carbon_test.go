package fuel

import (
	"testing"

	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
)

func TestHydroCarbon_GetCombustionGas(t *testing.T) {
	var gas = gases.GetAir()
	var hc = NewHydroCarbon(1, 4)

	var gas1 = hc.GetCombustionGas(gas, 2)

	var gas2 = hc.GetCombustionGas(gas1, 0.5)
	assert.InDelta(t, 0, gas2.OxygenMassFraction(), 1e-7)

	assert.Panics(t, func() {
		hc.GetCombustionGas(gas2, 1)
	})
}
