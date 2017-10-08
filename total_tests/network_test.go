package total_tests

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/nodes/source"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

const (
	tAtm          = 300
	pAtm          = 1e5
	etaCompressor = 0.86
	piStag        = 8

	etaTBlocked = 0.92
	lambdaOut   = 0.3

	tgStag    = 1800
	tFuel     = 300
	sigmaBurn = 0.99
	etaBurn   = 0.99
	initAlpha = 3
	t0        = 288

	etaFreeT = 0.92

	pressureLossSigma = 0.98

	precision = 0.05
)

func TestNetwork_Solve_OK(t *testing.T) {
	var nodes = make(map[string]core.Node)

	var gasSource1 = source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm)
	nodes["gasSource1"] = gasSource1
	var gasSource2 = source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm)
	nodes["gasSource2"] = gasSource2
	var gasSink = sink.NewComplexGasSinkNode()
	nodes["gasSink"] = gasSink
	var compressor = constructive.NewCompressorNode(etaCompressor, piStag, precision)
	nodes["compressor"] = compressor
	var turbine = constructive.NewBlockedTurbineNode(etaTBlocked, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	nodes["turbine"] = turbine
	var burner = constructive.NewBurnerNode(fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision)
	nodes["burner"] = burner
	var freeTurbine = constructive.NewFreeTurbineNode(etaFreeT, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	nodes["freeTurbine"] = freeTurbine
	var powerSink1 = sink.NewPowerSinkNode()
	nodes["powerSink1"] = powerSink1
	var powerSink2 = sink.NewPowerSinkNode()
	nodes["powerSink2"] = powerSink2
	var pressureLossNode = constructive.NewPressureLossNode(pressureLossSigma)
	nodes["loss"] = pressureLossNode
	var assembler = helper.NewGasStateAssemblerNode()
	nodes["assembler"] = assembler
	var disassembler = helper.NewGasStateDisassemblerNode()
	nodes["disassembler"] = disassembler
	var tSink = sink.NewTemperatureSinkNode()
	nodes["tSink1"] = tSink
	var mSink = sink.NewMassRateRelSinkNode()
	nodes["mSink"] = mSink
	var gSink = sink.NewGasSinkNode()
	nodes["gSink"] = gSink
	var hub = helper.NewHubNode(states.NewPressurePortState(1e5))
	nodes["hub"] = hub

	core.Link(compressor.ComplexGasInput(), gasSource1.ComplexGasOutput())
	core.Link(compressor.ComplexGasOutput(), burner.ComplexGasInput())
	core.Link(burner.ComplexGasOutput(), turbine.ComplexGasInput())
	core.Link(turbine.ComplexGasOutput(), pressureLossNode.ComplexGasInput())
	core.Link(pressureLossNode.ComplexGasOutput(), freeTurbine.ComplexGasInput())

	core.Link(gasSource2.ComplexGasOutput(), disassembler.ComplexGasPort())
	core.Link(disassembler.TemperaturePort(), tSink.TemperatureInput())
	core.Link(disassembler.GasPort(), gSink.GasInput())
	core.Link(disassembler.MassRateRelPort(), mSink.MassRateRelInput())

	core.Link(disassembler.PressurePort(), hub.Inlet())
	core.Link(hub.Outlet1(), freeTurbine.PressureOutput())

	core.Link(hub.Outlet2(), assembler.PressurePort())
	core.Link(freeTurbine.TemperatureOutput(), assembler.TemperaturePort())
	core.Link(freeTurbine.MassRateRelOutput(), assembler.MassRateRelPort())
	core.Link(freeTurbine.GasOutput(), assembler.GasPort())

	core.Link(assembler.ComplexGasPort(), gasSink.ComplexGasInput())

	core.Link(compressor.PowerOutput(), turbine.PowerInput())
	core.Link(turbine.PowerOutput(), powerSink1.PowerInput())
	core.Link(freeTurbine.PowerOutput(), powerSink2.PowerInput())

	var network = core.NewNetwork(nodes)
	var converged, err = network.Solve(1, 100, 0.05)

	assert.Nil(t, err)
	assert.True(t, converged)
}

func TestNetwork_Solve_FreePorts(t *testing.T) {
	var nodes = make(map[string]core.Node)

	var gasSource1 = source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm)
	nodes["gasSource1"] = gasSource1
	var gasSource2 = source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm)
	nodes["gasSource2"] = gasSource2
	var gasSink = sink.NewComplexGasSinkNode()
	nodes["gasSink"] = gasSink
	var compressor = constructive.NewCompressorNode(etaCompressor, piStag, precision)
	nodes["compressor"] = compressor
	var turbine = constructive.NewBlockedTurbineNode(etaTBlocked, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	nodes["turbine"] = turbine
	var burner = constructive.NewBurnerNode(fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision)
	nodes["burner"] = burner
	var freeTurbine = constructive.NewFreeTurbineNode(etaFreeT, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	nodes["freeTurbine"] = freeTurbine
	var powerSink1 = sink.NewPowerSinkNode()
	nodes["powerSink1"] = powerSink1
	var powerSink2 = sink.NewPowerSinkNode()
	nodes["powerSink2"] = powerSink2
	var pressureLossNode = constructive.NewPressureLossNode(pressureLossSigma)
	nodes["loss"] = pressureLossNode
	var assembler = helper.NewGasStateAssemblerNode()
	nodes["assembler"] = assembler
	var disassembler = helper.NewGasStateDisassemblerNode()
	nodes["disassembler"] = disassembler
	var tSink = sink.NewTemperatureSinkNode()
	nodes["tSink1"] = tSink
	var mSink = sink.NewMassRateRelSinkNode()
	nodes["mSink"] = mSink
	var gSink = sink.NewGasSinkNode()
	nodes["gSink"] = gSink
	var hub = helper.NewHubNode(states.NewPressurePortState(1e5))
	nodes["hub"] = hub

	core.Link(compressor.ComplexGasInput(), gasSource1.ComplexGasOutput())
	core.Link(compressor.ComplexGasOutput(), burner.ComplexGasInput())
	core.Link(burner.ComplexGasOutput(), turbine.ComplexGasInput())
	core.Link(turbine.ComplexGasOutput(), pressureLossNode.ComplexGasInput())
	core.Link(pressureLossNode.ComplexGasOutput(), freeTurbine.ComplexGasInput())

	core.Link(gasSource2.ComplexGasOutput(), disassembler.ComplexGasPort())
	core.Link(disassembler.TemperaturePort(), tSink.TemperatureInput())
	core.Link(disassembler.GasPort(), gSink.GasInput())
	core.Link(disassembler.MassRateRelPort(), mSink.MassRateRelInput())

	core.Link(disassembler.PressurePort(), hub.Inlet())
	core.Link(hub.Outlet1(), freeTurbine.PressureOutput())

	core.Link(hub.Outlet2(), assembler.PressurePort())
	core.Link(freeTurbine.TemperatureOutput(), assembler.TemperaturePort())
	core.Link(freeTurbine.MassRateRelOutput(), assembler.MassRateRelPort())
	core.Link(freeTurbine.GasOutput(), assembler.GasPort())

	var network = core.NewNetwork(nodes)
	var _, err = network.Solve(1, 100, 0.05)

	assert.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "Found free port"))
}

func TestNetwork_Solve_NotContextDefined(t *testing.T) {
	var nodes = make(map[string]core.Node)
	var pressureLoss1 = constructive.NewPressureLossNode(pressureLossSigma)
	nodes["loss1"] = pressureLoss1
	var pressureLoss2 = constructive.NewPressureLossNode(pressureLossSigma)
	nodes["loss2"] = pressureLoss2

	core.Link(pressureLoss1.ComplexGasOutput(), pressureLoss2.ComplexGasInput())
	core.Link(pressureLoss2.ComplexGasOutput(), pressureLoss1.ComplexGasInput())

	var network = core.NewNetwork(nodes)
	var _, err = network.Solve(1, 100, 0.05)

	assert.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "Nodes"))
}

func TestNetwork_Solve_NoStart(t *testing.T) {
	var nodes = make(map[string]core.Node)
	var compressor1 = constructive.NewCompressorNode(etaCompressor, piStag, precision)
	nodes["compressor1"] = compressor1
	var compressor2 = constructive.NewCompressorNode(etaCompressor, piStag, precision)
	nodes["compressor2"] = compressor2
	var sink1 = sink.NewPowerSinkNode()
	nodes["sink1"] = sink1
	var sink2 = sink.NewPowerSinkNode()
	nodes["sink1"] = sink2

	core.Link(compressor1.ComplexGasOutput(), compressor2.ComplexGasInput())
	core.Link(compressor2.ComplexGasOutput(), compressor1.ComplexGasInput())
	core.Link(compressor1.PowerOutput(), sink1.PowerInput())
	core.Link(compressor2.PowerOutput(), sink2.PowerInput())

	var network = core.NewNetwork(nodes)
	var _, err = network.Solve(1, 100, 0.05)

	assert.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "network has no start"))
}

func TestNetwork_Solve_Cycled(t *testing.T) {
	var nodes = make(map[string]core.Node)

	var gasSource1 = source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm)
	nodes["gasSource1"] = gasSource1
	var gasSource2 = source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm)
	nodes["gasSource2"] = gasSource2
	var gasSink = sink.NewComplexGasSinkNode()
	nodes["gasSink"] = gasSink
	var compressor = constructive.NewCompressorNode(etaCompressor, piStag, precision)
	nodes["compressor"] = compressor
	var turbine = constructive.NewBlockedTurbineNode(etaTBlocked, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	nodes["turbine"] = turbine
	var burner = constructive.NewBurnerNode(fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision)
	nodes["burner"] = burner
	var freeTurbine = constructive.NewFreeTurbineNode(etaFreeT, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	nodes["freeTurbine"] = freeTurbine
	var powerSink1 = sink.NewPowerSinkNode()
	nodes["powerSink1"] = powerSink1
	var powerSink2 = sink.NewPowerSinkNode()
	nodes["powerSink2"] = powerSink2
	var pressureLossNode = constructive.NewPressureLossNode(pressureLossSigma)
	nodes["loss"] = pressureLossNode
	var assembler = helper.NewGasStateAssemblerNode()
	nodes["assembler"] = assembler
	var disassembler = helper.NewGasStateDisassemblerNode()
	nodes["disassembler"] = disassembler
	var tSink = sink.NewTemperatureSinkNode()
	nodes["tSink1"] = tSink
	var mSink = sink.NewMassRateRelSinkNode()
	nodes["mSink"] = mSink
	var gSink = sink.NewGasSinkNode()
	nodes["gSink"] = gSink
	var hub = helper.NewHubNode(states.NewPressurePortState(1e5))
	nodes["hub"] = hub
	var regenerator = constructive.NewRegeneratorNode(0.9, 0.05, constructive.SigmaByColdSide)
	nodes["regenerator"] = regenerator
	var outflow = constructive.NewPressureLossNode(pressureLossSigma)
	nodes["outflow"] = outflow

	core.Link(compressor.ComplexGasInput(), gasSource1.ComplexGasOutput())
	core.Link(compressor.ComplexGasOutput(), regenerator.ColdInput())
	core.Link(regenerator.ColdOutput(), burner.ComplexGasInput())
	core.Link(burner.ComplexGasOutput(), turbine.ComplexGasInput())
	core.Link(turbine.ComplexGasOutput(), pressureLossNode.ComplexGasInput())
	core.Link(pressureLossNode.ComplexGasOutput(), freeTurbine.ComplexGasInput())

	core.Link(gasSource2.ComplexGasOutput(), disassembler.ComplexGasPort())
	core.Link(disassembler.TemperaturePort(), tSink.TemperatureInput())
	core.Link(disassembler.GasPort(), gSink.GasInput())
	core.Link(disassembler.MassRateRelPort(), mSink.MassRateRelInput())

	core.Link(disassembler.PressurePort(), hub.Inlet())
	core.Link(hub.Outlet1(), freeTurbine.PressureOutput())

	core.Link(hub.Outlet2(), assembler.PressurePort())
	core.Link(freeTurbine.TemperatureOutput(), assembler.TemperaturePort())
	core.Link(freeTurbine.MassRateRelOutput(), assembler.MassRateRelPort())
	core.Link(freeTurbine.GasOutput(), assembler.GasPort())

	core.Link(assembler.ComplexGasPort(), regenerator.HotInput())
	core.Link(regenerator.HotOutput(), outflow.ComplexGasInput())
	core.Link(outflow.ComplexGasOutput(), gasSink.ComplexGasInput())

	core.Link(compressor.PowerOutput(), turbine.PowerInput())
	core.Link(turbine.PowerOutput(), powerSink1.PowerInput())
	core.Link(freeTurbine.PowerOutput(), powerSink2.PowerInput())

	var network = core.NewNetwork(nodes)
	var _, err = network.Solve(1, 100, 0.05)

	assert.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "Roots"))
}

func TestNetwork_Solve_NoCycle(t *testing.T) {
	var nodes = make(map[string]core.Node)

	var gasSource1 = source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm)
	nodes["gasSource1"] = gasSource1
	var gasSource2 = source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm)
	nodes["gasSource2"] = gasSource2
	var gasSink = sink.NewComplexGasSinkNode()
	nodes["gasSink"] = gasSink
	var compressor = constructive.NewCompressorNode(etaCompressor, piStag, precision)
	nodes["compressor"] = compressor
	var turbine = constructive.NewBlockedTurbineNode(etaTBlocked, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	nodes["turbine"] = turbine
	var burner = constructive.NewBurnerNode(fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision)
	nodes["burner"] = burner
	var freeTurbine = constructive.NewFreeTurbineNode(etaFreeT, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	nodes["freeTurbine"] = freeTurbine
	var powerSink1 = sink.NewPowerSinkNode()
	nodes["powerSink1"] = powerSink1
	var powerSink2 = sink.NewPowerSinkNode()
	nodes["powerSink2"] = powerSink2
	var pressureLossNode = constructive.NewPressureLossNode(pressureLossSigma)
	nodes["loss"] = pressureLossNode
	var assembler = helper.NewGasStateAssemblerNode()
	nodes["assembler"] = assembler
	var disassembler = helper.NewGasStateDisassemblerNode()
	nodes["disassembler"] = disassembler
	var tSink = sink.NewTemperatureSinkNode()
	nodes["tSink1"] = tSink
	var mSink = sink.NewMassRateRelSinkNode()
	nodes["mSink"] = mSink
	var gSink = sink.NewGasSinkNode()
	nodes["gSink"] = gSink
	var hub = helper.NewHubNode(states.NewPressurePortState(1e5))
	nodes["hub"] = hub
	var regenerator = constructive.NewRegeneratorNode(0.9, 0.05, constructive.SigmaByColdSide)
	nodes["regenerator"] = regenerator
	var outflow = constructive.NewPressureLossNode(pressureLossSigma)
	nodes["outflow"] = outflow
	var breaker = helper.NewCycleBreakerNode(states.StandardAtmosphereState())
	nodes["breaker"] = breaker

	core.Link(compressor.ComplexGasInput(), gasSource1.ComplexGasOutput())
	core.Link(compressor.ComplexGasOutput(), regenerator.ColdInput())
	core.Link(regenerator.ColdOutput(), burner.ComplexGasInput())
	core.Link(burner.ComplexGasOutput(), turbine.ComplexGasInput())
	core.Link(turbine.ComplexGasOutput(), pressureLossNode.ComplexGasInput())
	core.Link(pressureLossNode.ComplexGasOutput(), freeTurbine.ComplexGasInput())

	core.Link(gasSource2.ComplexGasOutput(), outflow.ComplexGasOutput())
	core.Link(outflow.ComplexGasInput(), disassembler.ComplexGasPort())
	core.Link(disassembler.TemperaturePort(), tSink.TemperatureInput())
	core.Link(disassembler.GasPort(), gSink.GasInput())
	core.Link(disassembler.MassRateRelPort(), mSink.MassRateRelInput())

	core.Link(disassembler.PressurePort(), hub.Inlet())
	core.Link(hub.Outlet1(), freeTurbine.PressureOutput())

	core.Link(hub.Outlet2(), assembler.PressurePort())
	core.Link(freeTurbine.TemperatureOutput(), assembler.TemperaturePort())
	core.Link(freeTurbine.MassRateRelOutput(), assembler.MassRateRelPort())
	core.Link(freeTurbine.GasOutput(), assembler.GasPort())

	core.Link(assembler.ComplexGasPort(), breaker.DataSourcePort())
	core.Link(breaker.UpdatePort(), regenerator.HotInput())
	core.Link(regenerator.HotOutput(), gasSink.ComplexGasInput())

	core.Link(compressor.PowerOutput(), turbine.PowerInput())
	core.Link(turbine.PowerOutput(), powerSink1.PowerInput())
	core.Link(freeTurbine.PowerOutput(), powerSink2.PowerInput())

	var network = core.NewNetwork(nodes)
	var converged, err = network.Solve(1, 100, 0.0)

	assert.Nil(t, err)
	assert.True(t, converged)

	var callOrder, err1 = network.GetCallOrder()
	assert.Nil(t, err1)
	fmt.Println(callOrder)
	var b, _ = json.MarshalIndent(network, "", "    ")
	os.Stdout.Write(b)
}
