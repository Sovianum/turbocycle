package constructive

import (
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/common/gdf"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

const (
	etaT      = 0.92
	lambdaOut = 0.3
	alphaT    = 3.5
	tBlockedT = 1800
	pBlockedT = 6e5
	lBlockedT = 2.2e5
)

func TestBlockedTurbineNode_Process(t *testing.T) {
	var turbine = getTestBlockedTurbine()
	assert.NotNil(t, turbine)

	var gasState = states.NewComplexGasPortState(fuel.GetCH4().GetCombustionGas(alphaT), tBlockedT, pBlockedT, 1)
	turbine.ComplexGasInput().SetState(gasState)

	var powerState = states.NewPowerPortState(-lBlockedT)
	turbine.PowerInput().SetState(powerState)

	turbine.Process()

	var cpMean = gases.CpMean(gasState.Gas, turbine.TStagOut(), turbine.TStagIn(), nodes.DefaultN)
	var kMean = gases.KMean(gasState.Gas, turbine.TStagOut(), turbine.TStagIn(), nodes.DefaultN)
	var expectedTt = turbine.TStagIn() - lBlockedT/cpMean
	assert.True(
		t,
		common.ApproxEqual(expectedTt, turbine.TStagOut(), 0.01),
		fmt.Sprintf("Expected T_t %f, got %f", expectedTt, turbine.TStagOut()),
	)

	var expectedPit = math.Pow(
		1-1/etaT*(1-expectedTt/tBlockedT),
		kMean/(1-kMean),
	)
	assert.True(
		t,
		common.ApproxEqual(expectedPit, turbine.PiTStag(), 0.01),
		fmt.Sprintf("Expected pi_t %f, got %f", expectedPit, turbine.PiTStag()),
	)

	var expectedPt = turbine.PStagIn() / (turbine.PiTStag() * gdf.Pi(lambdaOut, kMean))
	assert.True(
		t,
		common.ApproxEqual(expectedPt, turbine.PStagOut(), 0.01),
		fmt.Sprintf("Expected p_t %f, got %f", expectedPt, turbine.PStagOut()),
	)
}

func getTestBlockedTurbine() BlockedTurbineNode {
	return NewBlockedTurbineNode(
		etaT, lambdaOut, 0.05,
		func(node TurbineNode) float64 {
			return 0
		},
		func(node TurbineNode) float64 {
			return 0
		},
		func(node TurbineNode) float64 {
			return 0
		},
	)
}
