package nodes

import (
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/Sovianum/turbocycle/core"
)

const (
	pressureLossSigma = 0.95
)

func TestPressureLossNode_Process_Inflow(t *testing.T) {
	var pressureLossNode = getTestPressureLossNode()
	var compressorNode = getTestCompressor()

	core.Link(compressorNode.GasOutput(), pressureLossNode.GasInput())
	core.Link(compressorNode.GasInput(), pressureLossNode.GasOutput())

	var inputState = states.NewGasPortState(gases.GetAir(), tA, pA, 1)
	pressureLossNode.GasInput().SetState(inputState)

	var err = pressureLossNode.Process()
	assert.Nil(t, err)

	var pOut = pressureLossNode.GasOutput().GetState().(states.GasPortState).PStag
	assert.True(
		t,
		common.ApproxEqual(pA*pressureLossSigma, pOut, 0.001),
		fmt.Sprintf("Expected p_out %f, got %f", pA*pressureLossSigma, pOut),
	)
}

func TestPressureLossNode_Process_Outflow(t *testing.T) {
	var pressureLossNode = getTestPressureLossNode()
	var compressorNode = getTestCompressor()

	core.Link(compressorNode.GasOutput(), pressureLossNode.GasOutput())
	core.Link(compressorNode.GasInput(), pressureLossNode.GasInput())

	var inputState = states.NewGasPortState(gases.GetAir(), tA, pA, 1)
	pressureLossNode.GasOutput().SetState(inputState)

	var err = pressureLossNode.Process()
	assert.Nil(t, err)

	var pIn = pressureLossNode.GasInput().GetState().(states.GasPortState).PStag
	assert.True(
		t,
		common.ApproxEqual(pA/pressureLossSigma, pIn, 0.001),
		fmt.Sprintf("Expected p_in %f, got %f", pA/pressureLossSigma, pIn),
	)
}

func TestPressureLossNode_Process_BiFlow(t *testing.T) {
	var pressureLossNode = getTestPressureLossNode()
	var compressorNode1 = getTestCompressor()
	var compressorNode2 = getTestCompressor()

	core.Link(compressorNode1.GasOutput(), pressureLossNode.GasInput())
	core.Link(compressorNode2.GasOutput(), pressureLossNode.GasOutput())

	var gasState = states.NewGasPortState(gases.GetAir(), tA, pA, 1)
	pressureLossNode.GasInput().SetState(gasState)
	pressureLossNode.GasOutput().SetState(gasState)

	var err = pressureLossNode.Process()
	assert.Nil(t, err)

	var pOut = pressureLossNode.GasOutput().GetState().(states.GasPortState).PStag
	assert.True(
		t,
		common.ApproxEqual(pA*pressureLossSigma, pOut, 0.001),
		fmt.Sprintf("Expected p_out %f, got %f", pA*pressureLossSigma, pOut),
	)

	var pIn = pressureLossNode.GasInput().GetState().(states.GasPortState).PStag
	assert.True(
		t,
		common.ApproxEqual(pA/pressureLossSigma, pIn, 0.001),
		fmt.Sprintf("Expected p_in %f, got %f", pA/pressureLossSigma, pIn),
	)
}

func getTestPressureLossNode() PressureLossNode {
	return NewPressureLossNode(pressureLossSigma)
}
