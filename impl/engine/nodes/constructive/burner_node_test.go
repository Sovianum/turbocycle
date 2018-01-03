package constructive

import (
	"fmt"
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
)

const (
	tgStag    = 1800
	tFuel     = 300
	t0        = 290
	sigmaBurn = 0.99
	etaBurn   = 0.99
	tInBurn   = 500
	pInBurn   = 6
)

func TestBurnerNode_Process(t *testing.T) {
	var bn = getTestBurner()

	var err = graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()),
			states.NewTemperaturePortState(tInBurn),
			states.NewPressurePortState(pInBurn),
			states.NewMassRateRelPortState(1),
		},
		[]graph.Port{
			bn.GasInput(), bn.TemperatureInput(), bn.PressureInput(), bn.MassRateInput(),
		},
	)
	if err != nil {
		panic(err)
	}

	bn.Process()
	assert.NotNil(t, bn)

	var expectedPOut = sigmaBurn * pInBurn
	assert.True(
		t,
		common.ApproxEqual(expectedPOut, bn.PStagOut(), 0.001),
		fmt.Sprintf("Expected p_stag_out %f, got %f", expectedPOut, bn.PStagOut()),
	)

	var cpGas = gases.CpMean(bn.GasInput().GetState().(states.GasPortState).Gas, t0, tgStag, nodes.DefaultN)
	var cpAir = gases.CpMean(gases.GetAir(), t0, tInBurn, nodes.DefaultN)
	var cpFuel = fuel.CpMean(fuel.GetCH4(), t0, tFuel, nodes.DefaultN)
	var enom = cpGas*(tgStag-t0) - cpAir*(tInBurn-t0)
	var denom = fuel.GetCH4().QLower()*etaBurn + cpFuel*(tFuel-t0) - cpGas*(tgStag-t0)
	var expectedFuelRate = enom / denom

	//todo make more precise calculations above to increase test accuracy
	assert.True(
		t,
		common.ApproxEqual(expectedFuelRate, bn.GetFuelRateRel(), 0.1),
		fmt.Sprintf("Expected g_m %f, got %f", expectedFuelRate, bn.GetFuelRateRel()),
	)

	var expectedAlpha = 1 / (expectedFuelRate * fuel.GetCH4().AirMassTheory())
	assert.True(
		t,
		common.ApproxEqual(expectedAlpha, bn.Alpha(), 0.1),
		fmt.Sprintf("Expected alpha %f, got %f", expectedAlpha, bn.Alpha()),
	)
}

func getTestBurner() BurnerNode {
	return NewBurnerNode(
		fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, 3.5, t0, 0.05,
	)
}
