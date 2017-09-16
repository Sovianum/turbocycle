package nodes

import (
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
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
	var outputGasState = states.NewComplexGasPortState(fuel.GetCH4().GetCombustionGas(alphaT), 300, pOutFreeT, 1)
	turbine.ComplexGasOutput().SetState(outputGasState)

	turbine.Process()

	var expectedPit = turbine.PStagIn() / turbine.PStagOut()
	assert.Equal(t, expectedPit, turbine.PiTStag())

	var k = gases.KMean(inputGasState.Gas, turbine.TStagOut(), turbine.TStagIn(), defaultN)
	var expectedTtStag = turbine.TStagIn() * (1 - (1-math.Pow(turbine.PiTStag(), (1-k)/k))*etaT)
	assert.True(
		t,
		common.ApproxEqual(expectedTtStag, turbine.TStagOut(), 0.01),
		fmt.Sprintf("Expected T_t %f, got %f", expectedTtStag, turbine.TStagOut()),
	)

	var cp = gases.CpMean(inputGasState.Gas, turbine.TStagOut(), turbine.TStagIn(), defaultN)
	var expectedLabour = cp * (turbine.TStagIn() - turbine.TStagOut())
	assert.True(
		t,
		common.ApproxEqual(expectedLabour, turbine.LSpecific(), 0.01),
		fmt.Sprintf("Expected L_t %f, got %f", expectedLabour, turbine.LSpecific()),
	)
}

func getTestFreeTurbineNode() FreeTurbineNode {
	return NewFreeTurbineNode(etaT, lambdaOut, 0.05, func(node TurbineNode) float64 {
		return 0
	})
}
