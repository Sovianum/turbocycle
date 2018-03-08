package schemes

import (
	"io/ioutil"
	"testing"

	"encoding/json"

	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/compose"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
)

func TestThreeShaftsCoolingRegeneratorScheme_GetNetwork_Smoke(t *testing.T) {
	var gasSource = source.NewComplexGasSourceNode(gases.GetAir(), 288, 1e5, 1)
	var inletPressureDrop = constructive.NewPressureLossNode(0.98)
	var middlePressureCascade = compose.NewTurboCascadeNode(
		0.86, 5,
		0.92, 0.3,
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		0.99, 0.05,
	)

	var regenerativeGasGenerator = compose.NewRegenerativeGasGeneratorNode(
		0.86, 6, fuel.GetCH4(),
		1400, 300, 0.99, 0.99, 3, 300,
		0.9, 0.3,
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		0.8, 0.99, 0.99, 0.05, 1, nodes.DefaultN,
	)
	var middlePressureCompressorPipe = constructive.NewPressureLossNode(0.98)
	var highPressureTurbinePipe = constructive.NewPressureLossNode(0.98)
	var cooler = constructive.NewCoolerNode(300, 0.98)
	var middlePressureTurbinePipe = constructive.NewPressureLossNode(0.98)
	var freeTurbineBlock = compose.NewFreeTurbineBlock(
		1e5,
		0.92, 0.3, 0.05,
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		0.9,
	)
	freeTurbineBlock.FreeTurbine().TemperatureOutput().SetState(states.NewTemperaturePortState(900))

	var scheme = NewThreeShaftsCoolingRegeneratorScheme(
		gasSource, inletPressureDrop, middlePressureCascade, cooler, regenerativeGasGenerator, middlePressureCompressorPipe,
		highPressureTurbinePipe, middlePressureTurbinePipe, freeTurbineBlock,
	)
	scheme.InitGasGeneratorCompressor(gases.GetAir(), 500, 5e5, 1)
	scheme.InitGasGeneratorHeatExchanger(gases.GetAir(), 900, 1e5, 1)

	var network, networkErr = scheme.GetNetwork()
	assert.Nil(t, networkErr)

	var solveErr = network.Solve(0.2, 1, 100, 0.05)
	assert.Nil(t, solveErr)

	var b, _ = json.MarshalIndent(network, "", "    ")
	ioutil.Discard.Write(b)
}
