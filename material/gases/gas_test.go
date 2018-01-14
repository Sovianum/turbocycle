package gases

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOxyFreeGas(t *testing.T) {
	var o2 = GetOxygen()
	var oxyFreeOx = GetOxyFreeGas(o2)
	assert.True(t, math.IsNaN(oxyFreeOx.R()))

	var air = GetAir()
	var oxyFreeAir = GetOxyFreeGas(air)
	assert.InDelta(t, 0, oxyFreeAir.OxygenMassFraction(), 1e-7)
}
