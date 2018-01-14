package gases

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOxyFreeGas(t *testing.T) {
	var gas = GetOxygen()
	var oxyFreeGas = GetOxyFreeGas(gas)
	assert.True(t, math.IsNaN(oxyFreeGas.R()))
}
