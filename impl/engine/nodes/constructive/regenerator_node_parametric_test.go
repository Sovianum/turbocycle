package constructive

import (
	"fmt"
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
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
	rn := getTestParametricRegenerator()
	assert.Nil(t, rn.Process())

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

func TestParametricRegeneratorNode_Consistency(t *testing.T) {
	rn := getTestParametricRegenerator()
	assert.Nil(t, rn.Process())
	sigma := rn.Sigma()

	assert.InDelta(t, sigma0, sigma, 1e-6)
}

func TestNewParametricRegeneratorNodeFromProto_Consistency(t *testing.T) {
	proto := getTestRegenerator()
	assert.Nil(t, proto.Process())

	rn := getTestParametricRegeneratorFromProto()
	assert.Nil(t, rn.Process())

	assert.InDelta(
		t,
		proto.ColdOutput().PressureOutput().GetState().Value().(float64),
		rn.ColdOutput().PressureOutput().GetState().Value().(float64),
		1e-6,
	)
	assert.InDelta(
		t,
		proto.ColdOutput().TemperatureOutput().GetState().Value().(float64),
		rn.ColdOutput().TemperatureOutput().GetState().Value().(float64),
		1e-6,
	)
	assert.InDelta(
		t,
		proto.HotOutput().PressureOutput().GetState().Value().(float64),
		rn.HotOutput().PressureOutput().GetState().Value().(float64),
		1e-6,
	)
	assert.InDelta(
		t,
		proto.HotOutput().TemperatureOutput().GetState().Value().(float64),
		rn.HotOutput().TemperatureOutput().GetState().Value().(float64),
		1e-1,
	)
	assert.InDelta(t, proto.Sigma(), rn.Sigma(), 1e-6)
}

func getTestParametricRegeneratorFromProto() RegeneratorNode {
	rn := getTestRegenerator()
	hotMassRate := hotMassRate0 * rn.HotInput().MassRateInput().GetState().Value().(float64)
	coldMassRate := hotMassRate0 * rn.ColdInput().MassRateInput().GetState().Value().(float64)

	result := NewParametricRegeneratorNodeFromProto(
		rn,
		hotMassRate, coldMassRate,
		velocityHot0, velocityCold0,
		hDiameterHot, hDiameterCold,
		1e-3, 1, nodes.DefaultN,
		FrowardTDrop, DefaultNuFunc, DefaultNuFunc,
	)
	result.HotInput().MassRateInput().SetState(states.NewMassRatePortState(hotMassRate))
	result.ColdInput().MassRateInput().SetState(states.NewMassRatePortState(coldMassRate))
	return result
}

func getTestParametricRegenerator() RegeneratorNode {
	result := NewParametricRegeneratorNode(
		gases.GetAir(), gases.GetAir(), hotMassRate0, coldMassRate0,
		hotTemperature0, coldTemperature0, hotPressure0, coldPressure0,
		velocityHot0, velocityCold0, sigma0, hDiameterHot, hDiameterCold,
		1e-3, 1, nodes.DefaultN,
		FrowardTDrop, DefaultNuFunc, DefaultNuFunc,
	)

	coldInput := result.ColdInput()
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

	hotInput := result.HotInput()
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
	return result
}
