package total_tests

import (
	"encoding/json"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"fmt"
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
	var gasSource1 = nodes.NewGasSource(gases.GetAir(), tAtm, pAtm)
	var gasSource2 = nodes.NewGasSource(gases.GetAir(), tAtm, pAtm)
	var compressor = nodes.NewCompressorNode(etaCompressor, piStag, precision)
	var turbine = nodes.NewBlockedTurbineNode(etaTBlocked, lambdaOut, precision, func(nodes.TurbineNode) float64 {
		return 0
	})
	var burner = nodes.NewBurnerNode(fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision)
	var freeTurbine = nodes.NewFreeTurbineNode(etaFreeT, lambdaOut, precision, func(nodes.TurbineNode) float64 {
		return 0
	})
	var powerSink1 = nodes.NewPortSinkNode()
	var powerSink2 = nodes.NewPortSinkNode()
	var pressureLossNode = nodes.NewPressureLossNode(pressureLossSigma)

	core.Link(compressor.GasInput(), gasSource1.GasOutput())
	core.Link(compressor.GasOutput(), burner.GasInput())
	core.Link(burner.GasOutput(), turbine.GasInput())
	core.Link(turbine.GasOutput(), pressureLossNode.GasInput())
	core.Link(pressureLossNode.GasOutput(), freeTurbine.GasInput())
	core.Link(freeTurbine.GasOutput(), gasSource2.GasOutput())

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
	var gasSource1 = nodes.NewGasSource(gases.GetAir(), tAtm, pAtm)
	var gasSource2 = nodes.NewGasSource(gases.GetAir(), tAtm, pAtm)
	var compressor = nodes.NewCompressorNode(etaCompressor, piStag, precision)
	var turbine = nodes.NewBlockedTurbineNode(etaTBlocked, lambdaOut, precision, func(nodes.TurbineNode) float64 {
		return 0
	})
	var burner = nodes.NewBurnerNode(fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision)
	var freeTurbine = nodes.NewFreeTurbineNode(etaFreeT, lambdaOut, precision, func(nodes.TurbineNode) float64 {
		return 0
	})
	var powerSink1 = nodes.NewPortSinkNode()
	var powerSink2 = nodes.NewPortSinkNode()
	var pressureLossNode = nodes.NewPressureLossNode(pressureLossSigma)

	core.Link(compressor.GasOutput(), burner.GasInput())
	core.Link(burner.GasOutput(), turbine.GasInput())
	core.Link(turbine.GasOutput(), pressureLossNode.GasInput())
	core.Link(pressureLossNode.GasOutput(), freeTurbine.GasInput())
	core.Link(freeTurbine.GasOutput(), gasSource2.GasOutput())

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
	var gasSource1 = nodes.NewGasSource(gases.GetAir(), tAtm, pAtm)
	var gasSource2 = nodes.NewGasSource(gases.GetAir(), tAtm, pAtm)
	var compressor = nodes.NewCompressorNode(etaCompressor, piStag, precision)
	var turbine = nodes.NewBlockedTurbineNode(etaTBlocked, lambdaOut, precision, func(nodes.TurbineNode) float64 {
		return 0
	})
	var burner = nodes.NewBurnerNode(fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision)
	var freeTurbine = nodes.NewFreeTurbineNode(etaFreeT, lambdaOut, precision, func(nodes.TurbineNode) float64 {
		return 0
	})
	var powerSink1 = nodes.NewPortSinkNode()
	var powerSink2 = nodes.NewPortSinkNode()
	var pressureLossNode = nodes.NewPressureLossNode(pressureLossSigma)

	core.Link(compressor.GasInput(), gasSource1.GasOutput())
	core.Link(compressor.GasOutput(), burner.GasInput())
	core.Link(burner.GasOutput(), turbine.GasInput())
	core.Link(turbine.GasOutput(), pressureLossNode.GasInput())
	core.Link(pressureLossNode.GasOutput(), freeTurbine.GasInput())
	core.Link(freeTurbine.GasOutput(), gasSource2.GasOutput())

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
	var gasSource1 = nodes.NewGasSource(gases.GetAir(), tAtm, pAtm)
	var gasSource2 = nodes.NewGasSource(gases.GetAir(), tAtm, pAtm)
	var compressor = nodes.NewCompressorNode(etaCompressor, piStag, precision)
	var turbine = nodes.NewBlockedTurbineNode(etaTBlocked, lambdaOut, precision, func(nodes.TurbineNode) float64 {
		return 0
	})
	var burner = nodes.NewBurnerNode(fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision)
	var freeTurbine = nodes.NewFreeTurbineNode(etaFreeT, lambdaOut, precision, func(nodes.TurbineNode) float64 {
		return 0
	})
	var powerSink1 = nodes.NewPortSinkNode()
	var powerSink2 = nodes.NewPortSinkNode()
	var pressureLossNode = nodes.NewPressureLossNode(pressureLossSigma)
	var regenerator = nodes.NewRegeneratorNode(0.9, 0.05, nodes.SigmaByColdSide)

	core.Link(gasSource1.GasOutput(), compressor.GasInput())
	core.Link(compressor.GasOutput(), regenerator.ColdInput())
	core.Link(regenerator.ColdOutput(), burner.GasInput())
	core.Link(burner.GasOutput(), turbine.GasInput())
	core.Link(turbine.GasOutput(), pressureLossNode.GasInput())
	core.Link(pressureLossNode.GasOutput(), freeTurbine.GasInput())
	core.Link(freeTurbine.GasOutput(), regenerator.HotInput())
	core.Link(regenerator.HotOutput(), gasSource2.GasOutput())

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
	var gasSource1 = nodes.NewGasSource(gases.GetAir(), tAtm, pAtm)
	var gasSource2 = nodes.NewGasSource(gases.GetAir(), tAtm, pAtm)
	var compressor = nodes.NewCompressorNode(etaCompressor, piStag, precision)
	var turbine = nodes.NewBlockedTurbineNode(etaTBlocked, lambdaOut, precision, func(nodes.TurbineNode) float64 {
		return 0
	})
	var burner = nodes.NewBurnerNode(fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision)
	var freeTurbine = nodes.NewFreeTurbineNode(etaFreeT, lambdaOut, precision, func(nodes.TurbineNode) float64 {
		return 0
	})
	var powerSink1 = nodes.NewPortSinkNode()
	var powerSink2 = nodes.NewPortSinkNode()
	var pressureLossNode = nodes.NewPressureLossNode(pressureLossSigma)
	var regenerator = nodes.NewRegeneratorNode(0.9, 0.05, nodes.SigmaByColdSide)
	var cycleBreaker = nodes.NewCycleBreakerNode(states.StandardAtmosphereState())
	var outflow = nodes.NewPressureLossNode(0.8)

	core.Link(gasSource1.GasOutput(), compressor.GasInput())
	core.Link(compressor.GasOutput(), regenerator.ColdInput())
	core.Link(regenerator.ColdOutput(), burner.GasInput())
	core.Link(burner.GasOutput(), turbine.GasInput())
	core.Link(turbine.GasOutput(), pressureLossNode.GasInput())
	core.Link(pressureLossNode.GasOutput(), freeTurbine.GasInput())
	core.Link(freeTurbine.GasOutput(), cycleBreaker.PortA())
	core.Link(cycleBreaker.PortB(), regenerator.HotInput())
	core.Link(regenerator.HotOutput(), outflow.GasInput())
	core.Link(outflow.GasOutput(), gasSource2.GasOutput())

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
