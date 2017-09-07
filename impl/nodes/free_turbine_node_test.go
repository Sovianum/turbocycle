package nodes

import (
	"fmt"
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFreeTurbineNode_Process(t *testing.T) { // smoke testing
	var ftn = NewFreeTurbineNode(0.92, 0.3, 0.05, func(node TurbineNode) float64 {
		return -0.01
	})

	assert.NotNil(t, ftn)

	var inputGasState = states.NewGasPortState(fuel.GetCH4().GetCombustionGas(4), 1200, 3e5, 1)
	ftn.GasInput().SetState(inputGasState)
	var outputGasState = states.NewGasPortState(gases.GetAir(), 300, 1e5, 1)
	ftn.GasOutput().SetState(outputGasState)

	ftn.Process()

	fmt.Println(ftn.GasOutput().GetState())
	fmt.Println(ftn.PowerOutput().GetState())
}
