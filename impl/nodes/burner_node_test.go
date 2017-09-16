package nodes

import (
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/stretchr/testify/assert"
	"testing"
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
	var inputState = states.NewComplexGasPortState(gases.GetAir(), tInBurn, pInBurn, 1)
	bn.ComplexGasInput().SetState(inputState)

	bn.Process()
	assert.NotNil(t, bn)

	var expectedPOut = sigmaBurn * pInBurn
	assert.True(
		t,
		common.ApproxEqual(expectedPOut, bn.PStagOut(), 0.001),
		fmt.Sprintf("Expected p_stag_out %f, got %f", expectedPOut, bn.PStagOut()),
	)

	var cpGas = gases.CpMean(bn.ComplexGasOutput().GetState().(states.ComplexGasPortState).Gas, t0, tgStag, defaultN)
	var cpAir = gases.CpMean(inputState.Gas, t0, tInBurn, defaultN)
	var cpFuel = fuel.CpMean(fuel.GetCH4(), t0, tFuel, defaultN)
	var enom = cpGas*(tgStag-t0) - cpAir*(tInBurn-t0)
	var denom = fuel.GetCH4().QLower()*etaBurn + cpFuel*(tFuel-t0) - cpGas*(tgStag-t0)
	var expectedFuelRate = enom / denom
	assert.True(
		t,
		common.ApproxEqual(expectedFuelRate, bn.GetFuelRateRel(), 0.01),
		fmt.Sprintf("Expected g_m %f, got %f", expectedFuelRate, bn.GetFuelRateRel()),
	)

	var expectedAlpha = 1 / (expectedFuelRate * fuel.GetCH4().AirMassTheory())
	assert.True(
		t,
		common.ApproxEqual(expectedAlpha, bn.Alpha(), 0.01),
		fmt.Sprintf("Expected alpha %f, got %f", expectedAlpha, bn.Alpha()),
	)
}

func getTestBurner() BurnerNode {
	return NewBurnerNode(
		fuel.GetCH4(), tgStag, tFuel, sigmaBurn, etaBurn, 3.5, t0, 0.05,
	)
}
