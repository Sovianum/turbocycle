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
	pressureLossSigma = 0.95
)

func TestPressureLossNode_Process_Inflow(t *testing.T) {
	var pressureLossNode = getTestPressureLossNode()
	var compressorNode = getTestCompressor()

	graph.LinkAll(
		[]graph.Port{
			compressorNode.GasOutput(), compressorNode.TemperatureOutput(),
			compressorNode.PressureOutput(), compressorNode.MassRateOutput(),
		},
		[]graph.Port{
			pressureLossNode.GasInput(), pressureLossNode.TemperatureInput(),
			pressureLossNode.PressureInput(), pressureLossNode.MassRateInput(),
		},
	)

	graph.LinkAll(
		[]graph.Port{
			pressureLossNode.GasOutput(), pressureLossNode.TemperatureOutput(),
			pressureLossNode.PressureOutput(), pressureLossNode.MassRateOutput(),
		},
		[]graph.Port{
			compressorNode.GasInput(), compressorNode.TemperatureInput(),
			compressorNode.PressureInput(), compressorNode.MassRateInput(),
		},
	)

	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()), states.NewTemperaturePortState(tA),
			states.NewPressurePortState(pA), states.NewMassRatePortState(1),
		},
		[]graph.Port{
			pressureLossNode.GasInput(), pressureLossNode.TemperatureInput(),
			pressureLossNode.PressureInput(), pressureLossNode.MassRateInput(),
		},
	)

	var err = pressureLossNode.Process()
	assert.Nil(t, err)

	var pOut = pressureLossNode.PressureOutput().GetState().(states.PressurePortState).PStag
	assert.True(
		t,
		common.ApproxEqual(pA*pressureLossSigma, pOut, 0.001),
		fmt.Sprintf("Expected p_out %f, got %f", pA*pressureLossSigma, pOut),
	)
}

func TestPressureLossNode_Process_Outflow(t *testing.T) {
	var pressureLossNode = getTestPressureLossNode()
	var compressorNode = getTestCompressor()

	graph.LinkAll(
		[]graph.Port{
			compressorNode.GasOutput(), compressorNode.TemperatureOutput(),
			compressorNode.PressureOutput(), compressorNode.MassRateOutput(),
		},
		[]graph.Port{
			pressureLossNode.GasOutput(), pressureLossNode.TemperatureOutput(),
			pressureLossNode.PressureOutput(), pressureLossNode.MassRateOutput(),
		},
	)

	graph.LinkAll(
		[]graph.Port{
			pressureLossNode.GasInput(), pressureLossNode.TemperatureInput(),
			pressureLossNode.PressureInput(), pressureLossNode.MassRateInput(),
		},
		[]graph.Port{
			compressorNode.GasInput(), compressorNode.TemperatureInput(),
			compressorNode.PressureInput(), compressorNode.MassRateInput(),
		},
	)

	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()), states.NewTemperaturePortState(tA),
			states.NewPressurePortState(pA), states.NewMassRatePortState(1),
		},
		[]graph.Port{
			pressureLossNode.GasOutput(), pressureLossNode.TemperatureOutput(),
			pressureLossNode.PressureOutput(), pressureLossNode.MassRateOutput(),
		},
	)

	var err = pressureLossNode.Process()
	assert.Nil(t, err)

	var pIn = pressureLossNode.PressureInput().GetState().(states.PressurePortState).PStag
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

	graph.LinkAll(
		[]graph.Port{
			compressorNode.GasOutput(), compressorNode.TemperatureOutput(),
			compressorNode.PressureOutput(), compressorNode.MassRateOutput(),
		},
		[]graph.Port{
			pln1.GasInput(), pln1.TemperatureInput(),
			pln1.PressureInput(), pln1.MassRateInput(),
		},
	)
	graph.LinkAll(
		[]graph.Port{
			pln1.GasOutput(), pln1.TemperatureOutput(),
			pln1.PressureOutput(), pln1.MassRateOutput(),
		},
		[]graph.Port{
			pln2.GasInput(), pln2.TemperatureInput(),
			pln2.PressureInput(), pln2.MassRateInput(),
		},
	)
	graph.LinkAll(
		[]graph.Port{
			pln2.GasOutput(), pln2.TemperatureOutput(),
			pln2.PressureOutput(), pln2.MassRateOutput(),
		},
		[]graph.Port{
			pln3.GasInput(), pln3.TemperatureInput(),
			pln3.PressureInput(), pln3.MassRateInput(),
		},
	)

	assert.True(t, pln1.ContextDefined())
	assert.True(t, pln2.ContextDefined())
	assert.True(t, pln3.ContextDefined())
}

func TestPressureLossNode_ContextDefined_False(t *testing.T) {
	var pln1 = getTestPressureLossNode()
	var pln2 = getTestPressureLossNode()
	var pln3 = getTestPressureLossNode()

	graph.LinkAll(
		[]graph.Port{
			pln1.GasOutput(), pln1.TemperatureOutput(),
			pln1.PressureOutput(), pln1.MassRateOutput(),
		},
		[]graph.Port{
			pln2.GasInput(), pln2.TemperatureInput(),
			pln2.PressureInput(), pln2.MassRateInput(),
		},
	)
	graph.LinkAll(
		[]graph.Port{
			pln2.GasOutput(), pln2.TemperatureOutput(),
			pln2.PressureOutput(), pln2.MassRateOutput(),
		},
		[]graph.Port{
			pln3.GasInput(), pln3.TemperatureInput(),
			pln3.PressureInput(), pln3.MassRateInput(),
		},
	)

	assert.False(t, pln1.ContextDefined())
	assert.False(t, pln2.ContextDefined())
	assert.False(t, pln3.ContextDefined())
}

func getTestPressureLossNode() PressureLossNode {
	return NewPressureLossNode(pressureLossSigma)
}
