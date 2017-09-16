package total_tests

import (
	"encoding/json"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/nodes/source"
	"github.com/Sovianum/turbocycle/impl/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"fmt"
	"github.com/Sovianum/turbocycle/impl/nodes/helper"
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
	var gasSource1 = source.NewComplexGasSource(gases.GetAir(), tAtm, pAtm)
	var gasSource2 = source.NewComplexGasSource(gases.GetAir(), tAtm, pAtm)
	var compressor = constructive.NewCompressorNode(etaCompressor, piStag, precision)
	var turbine = constructive.NewBlockedTurbineNode(etaTBlocked, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	var burner = constructive.NewBurnerNode(fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision)
	var freeTurbine = constructive.NewFreeTurbineNode(etaFreeT, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	var powerSink1 = sink.NewPowerSinkNode()
	var powerSink2 = sink.NewPowerSinkNode()
	var pressureLossNode = constructive.NewPressureLossNode(pressureLossSigma)

	core.Link(compressor.ComplexGasInput(), gasSource1.ComplexGasOutput())
	core.Link(compressor.ComplexGasOutput(), burner.ComplexGasInput())
	core.Link(burner.ComplexGasOutput(), turbine.ComplexGasInput())
	core.Link(turbine.ComplexGasOutput(), pressureLossNode.ComplexGasInput())
	core.Link(pressureLossNode.ComplexGasOutput(), freeTurbine.ComplexGasInput())
	core.Link(freeTurbine.ComplexGasOutput(), gasSource2.ComplexGasOutput())

	core.Link(compressor.PowerOutput(), turbine.PowerInput())
	core.Link(turbine.PowerOutput(), powerSink1.PowerInput())
	core.Link(freeTurbine.PowerOutput(), powerSink2.PowerInput())

	var network = core.NewNetwork(map[string]core.Node{
		"gasSource1":   gasSource1,
		"gasSource2":   gasSource2,
		"compressor":   compressor,
		"burner":       burner,
		"turbine":      turbine,
		"pressureLoss": pressureLossNode,
		"freeTurbine":  freeTurbine,
		"powerSink1":   powerSink1,
		"powerSink2":   powerSink2,
	})
	var converged, err = network.Solve(1, 100, 0.05)

	assert.Nil(t, err)
	assert.True(t, converged)
}

func TestNetwork_Solve_FreePorts(t *testing.T) {
	var gasSource1 = source.NewComplexGasSource(gases.GetAir(), tAtm, pAtm)
	var gasSource2 = source.NewComplexGasSource(gases.GetAir(), tAtm, pAtm)
	var compressor = constructive.NewCompressorNode(etaCompressor, piStag, precision)
	var turbine = constructive.NewBlockedTurbineNode(etaTBlocked, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	var burner = constructive.NewBurnerNode(fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision)
	var freeTurbine = constructive.NewFreeTurbineNode(etaFreeT, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	var powerSink1 = sink.NewPowerSinkNode()
	var powerSink2 = sink.NewPowerSinkNode()
	var pressureLossNode = constructive.NewPressureLossNode(pressureLossSigma)

	core.Link(compressor.ComplexGasOutput(), burner.ComplexGasInput())
	core.Link(burner.ComplexGasOutput(), turbine.ComplexGasInput())
	core.Link(turbine.ComplexGasOutput(), pressureLossNode.ComplexGasInput())
	core.Link(pressureLossNode.ComplexGasOutput(), freeTurbine.ComplexGasInput())
	core.Link(freeTurbine.ComplexGasOutput(), gasSource2.ComplexGasOutput())

	core.Link(compressor.PowerOutput(), turbine.PowerInput())
	core.Link(turbine.PowerOutput(), powerSink1.PowerInput())
	core.Link(freeTurbine.PowerOutput(), powerSink2.PowerInput())

	var network = core.NewNetwork(map[string]core.Node{
		"gasSource1":   gasSource1,
		"gasSource2":   gasSource2,
		"compressor":   compressor,
		"burner":       burner,
		"turbine":      turbine,
		"pressureLoss": pressureLossNode,
		"freeTurbine":  freeTurbine,
		"powerSink1":   powerSink1,
		"powerSink2":   powerSink2,
	})
	var _, err = network.Solve(1, 100, 0.05)

	assert.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "Found free port"))
}

func TestNetwork_Solve_CanNotCall(t *testing.T) {
	var gasSource1 = source.NewComplexGasSource(gases.GetAir(), tAtm, pAtm)
	var gasSource2 = source.NewComplexGasSource(gases.GetAir(), tAtm, pAtm)
	var compressor = constructive.NewCompressorNode(etaCompressor, piStag, precision)
	var turbine = constructive.NewBlockedTurbineNode(etaTBlocked, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	var burner = constructive.NewBurnerNode(fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision)
	var freeTurbine = constructive.NewFreeTurbineNode(etaFreeT, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	var powerSink1 = sink.NewPowerSinkNode()
	var powerSink2 = sink.NewPowerSinkNode()
	var pressureLossNode = constructive.NewPressureLossNode(pressureLossSigma)

	core.Link(compressor.ComplexGasInput(), gasSource1.ComplexGasOutput())
	core.Link(compressor.ComplexGasOutput(), burner.ComplexGasInput())
	core.Link(burner.ComplexGasOutput(), turbine.ComplexGasInput())
	core.Link(turbine.ComplexGasOutput(), pressureLossNode.ComplexGasInput())
	core.Link(pressureLossNode.ComplexGasOutput(), freeTurbine.ComplexGasInput())
	core.Link(freeTurbine.ComplexGasOutput(), gasSource2.ComplexGasOutput())

	core.Link(compressor.PowerOutput(), turbine.PowerInput())
	core.Link(turbine.PowerOutput(), powerSink1.PowerInput())
	core.Link(freeTurbine.PowerOutput(), powerSink2.PowerInput())

	var network = core.NewNetwork(map[string]core.Node{
		"gasSource1":   gasSource1,
		"compressor":   compressor,
		"burner":       burner,
		"turbine":      turbine,
		"pressureLoss": pressureLossNode,
		"freeTurbine":  freeTurbine,
		"powerSink1":   powerSink1,
		"powerSink2":   powerSink2,
	})
	var _, err = network.Solve(1, 100, 0.05)

	assert.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "Roots"))
}

func TestNetwork_Solve_Cycled(t *testing.T) {
	var gasSource1 = source.NewComplexGasSource(gases.GetAir(), tAtm, pAtm)
	var gasSource2 = source.NewComplexGasSource(gases.GetAir(), tAtm, pAtm)
	var compressor = constructive.NewCompressorNode(etaCompressor, piStag, precision)
	var turbine = constructive.NewBlockedTurbineNode(etaTBlocked, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	var burner = constructive.NewBurnerNode(fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision)
	var freeTurbine = constructive.NewFreeTurbineNode(etaFreeT, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	var powerSink1 = sink.NewPowerSinkNode()
	var powerSink2 = sink.NewPowerSinkNode()
	var pressureLossNode = constructive.NewPressureLossNode(pressureLossSigma)
	var regenerator = constructive.NewRegeneratorNode(0.9, 0.05, constructive.SigmaByColdSide)

	core.Link(gasSource1.ComplexGasOutput(), compressor.ComplexGasInput())
	core.Link(compressor.ComplexGasOutput(), regenerator.ColdInput())
	core.Link(regenerator.ColdOutput(), burner.ComplexGasInput())
	core.Link(burner.ComplexGasOutput(), turbine.ComplexGasInput())
	core.Link(turbine.ComplexGasOutput(), pressureLossNode.ComplexGasInput())
	core.Link(pressureLossNode.ComplexGasOutput(), freeTurbine.ComplexGasInput())
	core.Link(freeTurbine.ComplexGasOutput(), regenerator.HotInput())
	core.Link(regenerator.HotOutput(), gasSource2.ComplexGasOutput())

	core.Link(compressor.PowerOutput(), turbine.PowerInput())
	core.Link(turbine.PowerOutput(), powerSink1.PowerInput())
	core.Link(freeTurbine.PowerOutput(), powerSink2.PowerInput())

	var network = core.NewNetwork(map[string]core.Node{
		"gasSource1":   gasSource1,
		"gasSource2":   gasSource2,
		"compressor":   compressor,
		"burner":       burner,
		"turbine":      turbine,
		"pressureLoss": pressureLossNode,
		"freeTurbine":  freeTurbine,
		"powerSink1":   powerSink1,
		"powerSink2":   powerSink2,
		"regenerator":  regenerator,
	})
	var _, err = network.Solve(1, 100, 0.05)

	assert.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "Roots"))
}

func TestNetwork_Solve_NoCycle(t *testing.T) {
	var gasSource1 = source.NewComplexGasSource(gases.GetAir(), tAtm, pAtm)
	var gasSource2 = source.NewComplexGasSource(gases.GetAir(), tAtm, pAtm)
	var compressor = constructive.NewCompressorNode(etaCompressor, piStag, precision)
	var turbine = constructive.NewBlockedTurbineNode(etaTBlocked, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	var burner = constructive.NewBurnerNode(fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision)
	var freeTurbine = constructive.NewFreeTurbineNode(etaFreeT, lambdaOut, precision, func(constructive.TurbineNode) float64 {
		return 0
	})
	var powerSink1 = sink.NewPowerSinkNode()
	var powerSink2 = sink.NewPowerSinkNode()
	var pressureLossNode = constructive.NewPressureLossNode(pressureLossSigma)
	var regenerator = constructive.NewRegeneratorNode(0.9, 0.05, constructive.SigmaByColdSide)
	var cycleBreaker = helper.NewCycleBreakerNode(states.StandardAtmosphereState())
	var outflow = constructive.NewPressureLossNode(0.8)

	core.Link(gasSource1.ComplexGasOutput(), compressor.ComplexGasInput())
	core.Link(compressor.ComplexGasOutput(), regenerator.ColdInput())
	core.Link(regenerator.ColdOutput(), burner.ComplexGasInput())
	core.Link(burner.ComplexGasOutput(), turbine.ComplexGasInput())
	core.Link(turbine.ComplexGasOutput(), pressureLossNode.ComplexGasInput())
	core.Link(pressureLossNode.ComplexGasOutput(), freeTurbine.ComplexGasInput())
	core.Link(freeTurbine.ComplexGasOutput(), cycleBreaker.PortA())
	core.Link(cycleBreaker.PortB(), regenerator.HotInput())
	core.Link(regenerator.HotOutput(), outflow.ComplexGasInput())
	core.Link(outflow.ComplexGasOutput(), gasSource2.ComplexGasOutput())

	core.Link(compressor.PowerOutput(), turbine.PowerInput())
	core.Link(turbine.PowerOutput(), powerSink1.PowerInput())
	core.Link(freeTurbine.PowerOutput(), powerSink2.PowerInput())

	var network = core.NewNetwork(map[string]core.Node{
		"0_gasSource1":    gasSource1,
		"1_gasSource2":    gasSource2,
		"2_compressor":    compressor,
		"3_regenerator":   regenerator,
		"4_burner":        burner,
		"5_turbine":       turbine,
		"6_pressureLoss":  pressureLossNode,
		"7_freeTurbine":   freeTurbine,
		"8_powerSink1":    powerSink1,
		"9_powerSink2":    powerSink2,
		"a_cycle_breaker": cycleBreaker,
		"b_outflow":       outflow,
	})
	var converged, err = network.Solve(1, 100, 0.05)

	assert.Nil(t, err)
	assert.True(t, converged)

	fmt.Println(outflow.GetRequirePortTags())

	var callOrder, _ = network.GetCallOrder()
	fmt.Println(callOrder)
	var b, _ = json.MarshalIndent(network, "", "    ")
	os.Stdout.Write(b)
}
