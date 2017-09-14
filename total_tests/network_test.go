package total_tests

import (
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetwork_Solve_Smoke(t *testing.T) {
	var compressor = nodes.NewCompressorNode(0.86, 6, 0.05)
	var gasSource = nodes.NewGasSource(gases.GetAir(), 300, 1e5)
	var gasSink = nodes.NewGasSinkNode()
	var powerSink = nodes.NewPortSinkNode()
	var pressureLossNode = nodes.NewPressureLossNode(0.98)

	core.Link(compressor.GasInput(), gasSource.GasOutput())
	core.Link(compressor.GasOutput(), pressureLossNode.GasInput())
	core.Link(pressureLossNode.GasOutput(), gasSink.GasInput())
	core.Link(compressor.PowerOutput(), powerSink.PowerInput())

	var network = core.NewNetwork([]core.Node{compressor, gasSource, gasSink, powerSink, pressureLossNode})
	var converged, err = network.Solve(1, 100, 0.05)

	assert.Nil(t, err)
	assert.True(t, converged)

	fmt.Println(gasSink.GasInput().GetState())
}
