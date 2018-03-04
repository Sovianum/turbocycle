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

const (
	hotMassRate0  = 50.
	coldMassRate0 = 50.

	hotTemperature0  = 700.
	coldTemperature0 = 500.

	hotPressure0  = 3e5
	coldPressure0 = 5e5

	velocityHot0  = 30.
	velocityCold0 = 40.

	sigma0 = 0.9

	hDiameterHot  = 2e-3
	hDiameterCold = 3e-3
)

func TestParametricRegeneratorNode_Process_Unit(t *testing.T) {
	var rn = getTestParametricRegenerator()
	var coldInput = rn.ColdInput()
	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()), states.NewTemperaturePortState(coldTemperature0),
			states.NewPressurePortState(coldPressure0), states.NewMassRatePortState(coldMassRate0),
		},
		[]graph.Port{
			coldInput.GasInput(), coldInput.TemperatureInput(),
			coldInput.PressureInput(), coldInput.MassRateInput(),
		},
	)

	var hotInput = rn.HotInput()
	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()), states.NewTemperaturePortState(hotTemperature0),
			states.NewPressurePortState(hotPressure0), states.NewMassRatePortState(hotMassRate0),
		},
		[]graph.Port{
			hotInput.GasInput(), hotInput.TemperatureInput(),
			hotInput.PressureInput(), hotInput.MassRateInput(),
		},
	)

	var err = rn.Process()
	assert.Nil(t, err)

	var tColdOut = rn.ColdOutput().TemperatureOutput().GetState().(states.TemperaturePortState).TStag
	assert.True(
		t, common.ApproxEqual(680., tColdOut, 0.01),
		fmt.Sprintf("Expected %f, got %f (precision %f)", 680., tColdOut, 0.01),
	)

	var tHotOut = rn.HotOutput().TemperatureOutput().GetState().(states.TemperaturePortState).TStag
	assert.True(t,
		common.ApproxEqual(520., tHotOut, 0.03),
		fmt.Sprintf("Expected %f, got %f (precision %f)", 520., tHotOut, 0.01),
	)
}

func getTestParametricRegenerator() RegeneratorNode {
	return NewParametricRegeneratorNode(
		gases.GetAir(), gases.GetAir(), hotMassRate0, coldMassRate0,
		hotTemperature0, coldTemperature0, hotPressure0, coldPressure0,
		velocityHot0, velocityCold0, sigma0, hDiameterHot, hDiameterCold, 1e-3,
		FrowardTDrop, DefaultNuFunc, DefaultNuFunc,
	)
}
