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

func TestTwoShaftsRegeneratorScheme_GetNetwork_Smoke(t *testing.T) {
	var gasSource = source.NewComplexGasSourceNode(gases.GetAir(), 288, 1e5, 1)
	var inletPressureDrop = constructive.NewPressureLossNode(0.98)
	var turboCascade = compose.NewTurboCascadeNode(
		0.86, 6, 0.9, 0.3, func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		}, 0.99, 0.05,
	)
	var burner = constructive.NewBurnerNode(fuel.GetCH4(), 1400, 300, 0.99, 0.99, 3, 300, 0.05, 1, nodes.DefaultN)
	var compressorTurbinePipe = constructive.NewPressureLossNode(0.98)
	var freeTurbineBlock = compose.NewFreeTurbineBlock(
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
	var regenerator = constructive.NewRegeneratorNode(0.8, 0.05)

	var scheme = NewTwoShaftsRegeneratorScheme(
		gasSource, inletPressureDrop, turboCascade, burner, compressorTurbinePipe, freeTurbineBlock, regenerator,
	)
	var network, networkErr = scheme.GetNetwork()
	assert.Nil(t, networkErr)

	network.Solve(0.2, 1, 100, 0.05)
}
