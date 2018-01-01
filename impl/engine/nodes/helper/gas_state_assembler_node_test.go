package helper

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
)

const (
	pStag       = 1e8
	tStag       = 300.
	massRateRel = 8.
)

func TestGasStateAssemblerNode_ProcessAssemble(t *testing.T) {
	var gSource = source.NewGasSourceNode(gases.GetAir())
	var tSource = source.NewTemperatureSourceNode(tStag)
	var pSource = source.NewPressureSourceNode(pStag)
	var mSource = source.NewMassRateRelSourceNode(massRateRel)
	var complexSink = sink.NewComplexGasSinkNode()

	var assembler = NewGasStateAssemblerNode()
	graph.Link(assembler.GasPort(), gSource.GasOutput())
	graph.Link(assembler.TemperaturePort(), tSource.TemperatureOutput())
	graph.Link(assembler.PressurePort(), pSource.PressureOutput())
	graph.Link(assembler.MassRateRelPort(), mSource.MassRateRelOutput())
	graph.Link(assembler.ComplexGasPort(), complexSink.ComplexGasInput())

	var require = assembler.GetRequirePorts()
	assert.Equal(t, 4, len(require), len(require))

	var update = assembler.GetUpdatePorts()
	assert.Equal(t, 1, len(update), len(update))

	gSource.Process()
	tSource.Process()
	pSource.Process()
	mSource.Process()

	var err = assembler.Process()
	assert.Nil(t, err, err)

	complexSink.Process()

	assert.Equal(t, tStag, assembler.ComplexGasPort().GetState().(states.ComplexGasPortState).TStag)
	assert.Equal(t, pStag, assembler.ComplexGasPort().GetState().(states.ComplexGasPortState).PStag)
	assert.Equal(t, massRateRel, assembler.ComplexGasPort().GetState().(states.ComplexGasPortState).MassRateRel)
}

func TestGasStateAssemblerNode_ProcessDesemble(t *testing.T) {
	var gSink = sink.NewGasSinkNode()
	var tSink = sink.NewTemperatureSinkNode()
	var pSink = sink.NewPressureSinkNode()
	var mSink = sink.NewMassRateRelSinkNode()
	var complexSource = source.NewComplexGasSourceNode(gases.GetAir(), tStag, pStag)

	var assembler = NewGasStateDisassemblerNode()
	graph.Link(assembler.GasPort(), gSink.GasInput())
	graph.Link(assembler.TemperaturePort(), tSink.TemperatureInput())
	graph.Link(assembler.PressurePort(), pSink.PressureInput())
	graph.Link(assembler.MassRateRelPort(), mSink.MassRateRelInput())
	graph.Link(assembler.ComplexGasPort(), complexSource.ComplexGasOutput())

	mSink.MassRateRelInput().SetState(states.NewMassRateRelPortState(massRateRel))

	var require = assembler.GetRequirePorts()
	assert.Equal(t, 1, len(require))

	var update = assembler.GetUpdatePorts()
	assert.Equal(t, 4, len(update))

	complexSource.Process()

	var err = assembler.Process()
	assert.Nil(t, err, err)

	gSink.Process()
	tSink.Process()
	pSink.Process()
	mSink.Process()

	assert.Equal(t, tStag, assembler.TemperaturePort().GetState().(states.TemperaturePortState).TStag)
	assert.Equal(t, pStag, assembler.PressurePort().GetState().(states.PressurePortState).PStag)
	assert.Equal(t, 1., assembler.MassRateRelPort().GetState().(states.MassRateRelPortState).MassRateRel)
}
