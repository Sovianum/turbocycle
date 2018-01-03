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

func TestRegeneratorNode_Process_ColdMode(t *testing.T) {
	var rn = NewRegeneratorNode(0.9, 0.05, SigmaByColdSide)
	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()), states.NewTemperaturePortState(300),
			states.NewPressurePortState(1e5), states.NewMassRateRelPortState(1),
		},
		[]graph.Port{
			rn.ColdGasInput(), rn.ColdTemperatureInput(),
			rn.ColdPressureInput(), rn.ColdMassRateInput(),
		},
	)

	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()), states.NewTemperaturePortState(800),
			states.NewPressurePortState(1e5), states.NewMassRateRelPortState(1),
		},
		[]graph.Port{
			rn.HotGasInput(), rn.HotTemperatureInput(),
			rn.HotPressureInput(), rn.HotMassRateInput(),
		},
	)

	var err = rn.Process()
	assert.Nil(t, err)

	var tColdOut = rn.ColdTemperatureOutput().GetState().(states.TemperaturePortState).TStag
	assert.True(t, common.ApproxEqual(750, tColdOut, 0.01))

	var tHotOut = rn.HotTemperatureOutput().GetState().(states.TemperaturePortState).TStag
	assert.True(t,
		common.ApproxEqual(350, tHotOut, 0.03),
		fmt.Sprintf("Expected %f, got %f (precision %f)", 350., tHotOut, 0.01),
	)
}

func TestRegeneratorNode_Process_HotMode(t *testing.T) {
	var rn = NewRegeneratorNode(0.9, 0.05, SigmaByColdSide)
	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()), states.NewTemperaturePortState(300),
			states.NewPressurePortState(1e5), states.NewMassRateRelPortState(1),
		},
		[]graph.Port{
			rn.ColdGasInput(), rn.ColdTemperatureInput(),
			rn.ColdPressureInput(), rn.ColdMassRateInput(),
		},
	)

	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()), states.NewTemperaturePortState(800),
			states.NewPressurePortState(1e5), states.NewMassRateRelPortState(1),
		},
		[]graph.Port{
			rn.HotGasInput(), rn.HotTemperatureInput(),
			rn.HotPressureInput(), rn.HotMassRateInput(),
		},
	)

	var err = rn.Process()
	assert.Nil(t, err)

	var tColdOut = rn.ColdTemperatureOutput().GetState().(states.TemperaturePortState).TStag
	assert.True(t, common.ApproxEqual(750, tColdOut, 0.01))

	var tHotOut = rn.HotTemperatureOutput().GetState().(states.TemperaturePortState).TStag
	assert.True(t,
		common.ApproxEqual(350, tHotOut, 0.03),
		fmt.Sprintf("Expected %f, got %f (precision %f)", 350., tHotOut, 0.01),
	)
}
