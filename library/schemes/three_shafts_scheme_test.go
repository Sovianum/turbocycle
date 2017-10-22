package schemes

import (
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/compose"
	"github.com/Sovianum/turbocycle/helpers/fuel"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/helpers/gases"
	"fmt"
	"encoding/json"
	"os"
)

func TestThreeShaftsScheme_GetNetwork_Smoke(t *testing.T) {
	var gasSource = source.NewComplexGasSourceNode(gases.GetAir(), 288, 1e5)
	var inletPressureDrop = constructive.NewPressureLossNode(0.98)
	var middlePressureCascade = compose.NewTurboCascadeNode(
		0.86, 5,
		0.92, 0.3, func(node constructive.TurbineNode) float64 {
			return 0
		}, 0.99, 0.05,
	)
	var gasGenerator = compose.NewGasGeneratorNode(
		0.86, 6, fuel.GetCH4(),
		1400, 300, 0.99, 0.99, 3, 300,
		0.9, 0.3, func(node constructive.TurbineNode) float64 {
			return 0
		},
		0.99, 0.05,
	)
	var middlePressureCompressorPipe = constructive.NewPressureLossNode(0.98)
	var highPressureTurbinePipe = constructive.NewPressureLossNode(0.98)
	var middlePressureTurbinePipe = constructive.NewPressureLossNode(0.98)
	var freeTurbineBlock = compose.NewFreeTurbineBlock(
		1e5,
		0.92, 0.3, 0.05, func(node constructive.TurbineNode) float64 {
			return 0
		}, 0.9,
	)

	var scheme = NewThreeShaftsScheme(
		gasSource, inletPressureDrop, middlePressureCascade, gasGenerator, middlePressureCompressorPipe,
		highPressureTurbinePipe, middlePressureTurbinePipe, freeTurbineBlock,
	)
	var network = scheme.GetNetwork()
	var callOrder, err1 = network.GetCallOrder()
	assert.Nil(t, err1)
	fmt.Println(callOrder)
	network.Solve(0.2, 100, 0.05)
	var b, _ = json.MarshalIndent(network, "", "    ")
	os.Stdout.Write(b)
}
