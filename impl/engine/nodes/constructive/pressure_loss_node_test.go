package constructive

import (
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/helpers/gases"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	pressureLossSigma = 0.95
)

func TestPressureLossNode_Process_Inflow(t *testing.T) {
	var pressureLossNode = getTestPressureLossNode()
	var compressorNode = getTestCompressor()

	core.Link(compressorNode.ComplexGasOutput(), pressureLossNode.ComplexGasInput())
	core.Link(compressorNode.ComplexGasInput(), pressureLossNode.ComplexGasOutput())

	var inputState = states.NewComplexGasPortState(gases.GetAir(), tA, pA, 1)
	pressureLossNode.ComplexGasInput().SetState(inputState)

	var err = pressureLossNode.Process()
	assert.Nil(t, err)

	var pOut = pressureLossNode.ComplexGasOutput().GetState().(states.ComplexGasPortState).PStag
	assert.True(
		t,
		common.ApproxEqual(pA*pressureLossSigma, pOut, 0.001),
		fmt.Sprintf("Expected p_out %f, got %f", pA*pressureLossSigma, pOut),
	)
}

func TestPressureLossNode_Process_Outflow(t *testing.T) {
	var pressureLossNode = getTestPressureLossNode()
	var compressorNode = getTestCompressor()

	core.Link(compressorNode.ComplexGasOutput(), pressureLossNode.ComplexGasOutput())
	core.Link(compressorNode.ComplexGasInput(), pressureLossNode.ComplexGasInput())

	var inputState = states.NewComplexGasPortState(gases.GetAir(), tA, pA, 1)
	pressureLossNode.ComplexGasOutput().SetState(inputState)

	var err = pressureLossNode.Process()
	assert.Nil(t, err)

	var pIn = pressureLossNode.ComplexGasInput().GetState().(states.ComplexGasPortState).PStag
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

	core.Link(compressorNode1.ComplexGasOutput(), pressureLossNode.ComplexGasInput())
	core.Link(compressorNode2.ComplexGasOutput(), pressureLossNode.ComplexGasOutput())

	var gasState = states.NewComplexGasPortState(gases.GetAir(), tA, pA, 1)
	pressureLossNode.ComplexGasInput().SetState(gasState)
	pressureLossNode.ComplexGasOutput().SetState(gasState)

	var err = pressureLossNode.Process()
	assert.Nil(t, err)

	var pOut = pressureLossNode.ComplexGasOutput().GetState().(states.ComplexGasPortState).PStag
	assert.True(
		t,
		common.ApproxEqual(pA*pressureLossSigma, pOut, 0.001),
		fmt.Sprintf("Expected p_out %f, got %f", pA*pressureLossSigma, pOut),
	)

	var pIn = pressureLossNode.ComplexGasInput().GetState().(states.ComplexGasPortState).PStag
	assert.True(
		t,
		common.ApproxEqual(pA/pressureLossSigma, pIn, 0.001),
		fmt.Sprintf("Expected p_in %f, got %f", pA/pressureLossSigma, pIn),
	)
}

func TestPressureLossNode_ContextDefined_True(t *testing.T) {
	var compressorNode = getTestCompressor()
	var pln1 = getTestPressureLossNode()
	var pln2 = getTestPressureLossNode()
	var pln3 = getTestPressureLossNode()

	core.Link(compressorNode.ComplexGasOutput(), pln1.ComplexGasInput())
	core.Link(pln1.ComplexGasOutput(), pln2.ComplexGasInput())
	core.Link(pln2.ComplexGasOutput(), pln3.ComplexGasInput())

	assert.True(t, pln1.ContextDefined())
	assert.True(t, pln2.ContextDefined())
	assert.True(t, pln3.ContextDefined())
}

func TestPressureLossNode_ContextDefined_False(t *testing.T) {
	var pln1 = getTestPressureLossNode()
	var pln2 = getTestPressureLossNode()
	var pln3 = getTestPressureLossNode()

	core.Link(pln1.ComplexGasOutput(), pln2.ComplexGasInput())
	core.Link(pln2.ComplexGasOutput(), pln3.ComplexGasInput())

	assert.False(t, pln1.ContextDefined())
	assert.False(t, pln2.ContextDefined())
	assert.False(t, pln3.ContextDefined())
}

func getTestPressureLossNode() PressureLossNode {
	return NewPressureLossNode(pressureLossSigma)
}
