package nodes

import (
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	pressureLossSigma = 0.95
)

func TestPressureLossNode_Process(t *testing.T) {
	var node = getTestPressureLossNode()
	var inputState = states.NewGasPortState(gases.GetAir(), tA, pA, 1)
	var outputState = states.NewGasPortState(gases.GetAir(), tA, pA, 1)

	node.GasInput().SetState(inputState)
	node.GasOutput().SetState(outputState)

	node.Process()

	var pOut = node.GasOutput().GetState().(states.GasPortState).PStag
	assert.True(
		t,
		common.ApproxEqual(pA*pressureLossSigma, pOut, 0.001),
		fmt.Sprintf("Expected p_out %f, got %f", pA*pressureLossSigma, pOut),
	)

	var pIn = node.GasInput().GetState().(states.GasPortState).PStag
	assert.True(
		t,
		common.ApproxEqual(pA/pressureLossSigma, pIn, 0.001),
		fmt.Sprintf("Expected p_in %f, got %f", pA/pressureLossSigma, pOut),
	)
}

func getTestPressureLossNode() PressureLossNode {
	return NewPressureLossNode(pressureLossSigma)
}
