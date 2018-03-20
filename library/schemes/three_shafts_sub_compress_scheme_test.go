package schemes

import (
	"testing"

	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/compose"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
)

func TestThreeShaftsSubCompressScheme_GetNetwork_Smoke(t *testing.T) {
	gasSource := source.NewComplexGasSourceNode(gases.GetAir(), 288, 1e5, 1)
	inletPressureDrop := constructive.NewPressureLossNode(0.98)
	middlePressureCascade := compose.NewTurboCascadeNode(
		0.86, 5,
		0.92, 0.3, func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		}, 0.99, 0.05,
	)
	gasGenerator := compose.NewGasGeneratorNode(
		0.86, 6, fuel.GetCH4(),
		1400, 300, 0.99, 0.99, 3, 300,
		0.9, 0.3, func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		0.99, 0.05, 1, nodes.DefaultN,
	)
	middlePressureCompressorPipe := constructive.NewPressureLossNode(0.98)
	highPressureTurbinePipe := constructive.NewPressureLossNode(0.98)
	middlePressureTurbinePipe := constructive.NewPressureLossNode(0.98)
	freeTurbineBlock := compose.NewFreeTurbineBlock(
		1e5,
		0.92, 0.3, 0.05, func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		}, 0.9,
	)

	scheme := NewThreeShaftsSubCompressScheme(
		gasSource, inletPressureDrop, middlePressureCascade, gasGenerator, middlePressureCompressorPipe,
		highPressureTurbinePipe, middlePressureTurbinePipe, freeTurbineBlock,
		constructive.NewGasSplitter(0.1),
		constructive.NewGasCombiner(1e-6, 1, 100),
		constructive.NewCompressorNode(0.76, 2, 1e-5),
		constructive.NewCoolerNode(400, 0.98),
	)

	network, networkErr := scheme.GetNetwork()
	assert.Nil(t, networkErr)
	assert.Nil(t, network.Solve(0.2, 1, 100, 0.05))
}
