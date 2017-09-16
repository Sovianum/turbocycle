package total_tests

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetwork_Solve_Smoke(t *testing.T) {
	var gasSource1 = nodes.NewGasSource(gases.GetAir(), 300, 1e5)
	var gasSource2 = nodes.NewGasSource(gases.GetAir(), 300, 1e5)
	var compressor = nodes.NewCompressorNode(0.86, 6, 0.05)
	var turbine = nodes.NewBlockedTurbineNode(0.92, 0.3, 0.05, func(nodes.TurbineNode) float64 {
		return 0
	})
	var burner = nodes.NewBurnerNode(fuel.GetCH4(), 1800, 300, 0.99, 0.99, 3, 288, 0.05)
	var freeTurbine = nodes.NewFreeTurbineNode(0.92, 0.3, 0.05, func(nodes.TurbineNode) float64 {
		return 0
	})
	var powerSink1 = nodes.NewPortSinkNode()
	var powerSink2 = nodes.NewPortSinkNode()
	var pressureLossNode = nodes.NewPressureLossNode(0.98)

	core.Link(compressor.GasInput(), gasSource1.GasOutput())
	core.Link(compressor.GasOutput(), burner.GasInput())
	core.Link(burner.GasOutput(), turbine.GasInput())
	core.Link(turbine.GasOutput(), pressureLossNode.GasInput())
	core.Link(pressureLossNode.GasOutput(), freeTurbine.GasInput())
	core.Link(freeTurbine.GasOutput(), gasSource2.GasOutput())

	core.Link(compressor.PowerOutput(), turbine.PowerInput())
	core.Link(turbine.PowerOutput(), powerSink1.PowerInput())
	core.Link(freeTurbine.PowerOutput(), powerSink2.PowerInput())

	var network = core.NewNetwork([]core.Node{
		gasSource1, gasSource2, compressor, burner, turbine, pressureLossNode, freeTurbine, powerSink1, powerSink2,
	})
	var converged, err = network.Solve(1, 100, 0.05)

	assert.Nil(t, err)
	assert.True(t, converged)

	var s, _ = json.Marshal(compressor)
	fmt.Println(string(s))

	fmt.Println(freeTurbine.PiTStag())
	fmt.Println(network.GetState())
}
