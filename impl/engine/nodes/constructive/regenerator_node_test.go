package constructive

import (
	"fmt"
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
)

func TestRegeneratorNode_Process(t *testing.T) {
	rn := getTestRegenerator()
	assert.Nil(t, rn.Process())

	var tColdOut = rn.ColdOutput().TemperatureOutput().GetState().(states.TemperaturePortState).TStag
	assert.True(t, common.ApproxEqual(750, tColdOut, 0.01))

	var tHotOut = rn.HotOutput().TemperatureOutput().GetState().(states.TemperaturePortState).TStag
	assert.True(t,
		common.ApproxEqual(350, tHotOut, 0.03),
		fmt.Sprintf("Expected %f, got %f (precision %f)", 350., tHotOut, 0.01),
	)
}

func getTestRegenerator() RegeneratorNode {
	rn := NewRegeneratorNode(0.9, 0.05)
	coldInput := rn.ColdInput()
	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()), states.NewTemperaturePortState(300),
			states.NewPressurePortState(1e5), states.NewMassRatePortState(1),
		},
		[]graph.Port{
			coldInput.GasInput(), coldInput.TemperatureInput(),
			coldInput.PressureInput(), coldInput.MassRateInput(),
		},
	)

	hotInput := rn.HotInput()
	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()), states.NewTemperaturePortState(800),
			states.NewPressurePortState(1e5), states.NewMassRatePortState(1),
		},
		[]graph.Port{
			hotInput.GasInput(), hotInput.TemperatureInput(),
			hotInput.PressureInput(), hotInput.MassRateInput(),
		},
	)
	return rn
}
