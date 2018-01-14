package constructive

import (
	"fmt"
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/common/gdf"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
)

const (
	tInFreeT  = 1200
	pInFreeT  = 3e5
	pOutFreeT = 1.2e5
)

func TestFreeTurbineNode_Process(t *testing.T) {
	var turbine = getTestFreeTurbineNode()
	assert.NotNil(t, turbine)

	var inputGas = fuel.GetCH4().GetCombustionGas(gases.GetAir(), alphaT)
	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(inputGas), states.NewTemperaturePortState(tInFreeT),
			states.NewPressurePortState(pInFreeT), states.NewMassRatePortState(1),
		},
		[]graph.Port{
			turbine.GasInput(), turbine.TemperatureInput(), turbine.PressureInput(), turbine.MassRateInput(),
		},
	)
	turbine.PressureOutput().SetState(states.NewPressurePortState(pOutFreeT))

	turbine.Process()

	var expectedPitStag = turbine.PStagIn() / turbine.PStagOut()
	assert.Equal(t, expectedPitStag, turbine.PiTStag())

	var expectedPit = expectedPitStag / gdf.Pi(lambdaOut, kAir)

	var k = gases.KMean(inputGas, turbine.TStagOut(), turbine.TStagIn(), nodes.DefaultN)
	var expectedTtStag = turbine.TStagIn() * (1 - (1-math.Pow(expectedPit, (1-k)/k))*etaT)
	assert.True(
		t,
		common.ApproxEqual(expectedTtStag, turbine.TStagOut(), 0.01),
		fmt.Sprintf("Expected T_t %f, got %f", expectedTtStag, turbine.TStagOut()),
	)

	var cp = gases.CpMean(inputGas, turbine.TStagOut(), turbine.TStagIn(), nodes.DefaultN)
	var expectedLabour = cp * (turbine.TStagIn() - turbine.TStagOut())
	assert.True(
		t,
		common.ApproxEqual(expectedLabour, turbine.LSpecific(), 0.01),
		fmt.Sprintf("Expected L_t %f, got %f", expectedLabour, turbine.LSpecific()),
	)
}

func getTestFreeTurbineNode() FreeTurbineNode {
	return NewFreeTurbineNode(
		etaT, lambdaOut, 0.05,
		func(node TurbineNode) float64 {
			return 0
		},
		func(node TurbineNode) float64 {
			return 0
		},
		func(node TurbineNode) float64 {
			return 0
		},
	)
}
