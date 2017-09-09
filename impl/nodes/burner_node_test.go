package nodes

import (
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	tgStag      = 1800
	tFuel       = 300
	sigmaBurner = 0.99
	etaBurn     = 0.99
	tInBurn     = 500
	pInBurn     = 6
)

func TestBurnerNode_Smoke(t *testing.T) { // smoke testing
	var bn = NewBurnerNodeShort(fuel.GetCH4(), 1200, 300, 0.99, 0.999)
	var inputState = states.NewGasPortState(gases.GetAir(), 500, 6e5, 1)
	bn.GasInput().SetState(inputState)

	bn.Process()

	assert.NotNil(t, bn)
}

func getTestBurner() {

}
