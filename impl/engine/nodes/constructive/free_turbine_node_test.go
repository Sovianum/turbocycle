package constructive

import (
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/helpers/fuel"
	"github.com/Sovianum/turbocycle/helpers/gases"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"github.com/Sovianum/turbocycle/helpers/gdf"
)

const (
	tInFreeT  = 1200
	pInFreeT  = 3e5
	pOutFreeT = 1.2e5
)

func TestFreeTurbineNode_Process(t *testing.T) {
	var turbine = getTestFreeTurbineNode()
	assert.NotNil(t, turbine)

	var inputGasState = states.NewComplexGasPortState(fuel.GetCH4().GetCombustionGas(alphaT), tInFreeT, pInFreeT, 1)
	turbine.ComplexGasInput().SetState(inputGasState)
	turbine.PressureOutput().SetState(states.NewPressurePortState(pOutFreeT))

	turbine.Process()

	var expectedPitStag = turbine.PStagIn() / turbine.PStagOut()
	assert.Equal(t, expectedPitStag, turbine.PiTStag())

	var expectedPit = expectedPitStag / gdf.Pi(lambdaOut, kAir)

	var k = gases.KMean(inputGasState.Gas, turbine.TStagOut(), turbine.TStagIn(), nodes.DefaultN)
	var expectedTtStag = turbine.TStagIn() * (1 - (1-math.Pow(expectedPit, (1-k)/k))*etaT)
	assert.True(
		t,
		common.ApproxEqual(expectedTtStag, turbine.TStagOut(), 0.01),
		fmt.Sprintf("Expected T_t %f, got %f", expectedTtStag, turbine.TStagOut()),
	)

	var cp = gases.CpMean(inputGasState.Gas, turbine.TStagOut(), turbine.TStagIn(), nodes.DefaultN)
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
