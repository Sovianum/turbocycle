package nodes

import (
	"fmt"
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockedTurbineNode_Process(t *testing.T) { // smoke testing
	var btn = NewBlockedTurbineNodeShort(0.92, func(node TurbineNode) float64 {
		return 0
	})

	assert.NotNil(t, btn)

	var gasState = states.NewGasPortState(fuel.GetCH4().GetCombustionGas(3.5), 1800, 6e5, 1)
	btn.GasInput().SetState(gasState)

	var powerState = states.NewPowerPortState(-219742)
	btn.PowerInput().SetState(powerState)

	btn.Process()
	fmt.Println(btn.PowerOutput().GetState())
	fmt.Println(btn.GasInput().GetState())
	fmt.Println(btn.GasOutput().GetState())
}
