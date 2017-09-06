package nodes

import (
	"testing"
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/stretchr/testify/assert"
	"fmt"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/Sovianum/turbocycle/gases"
)

func TestBurnerNode_Smoke(t *testing.T) {	// smoke testing
	var bn = NewBurnerNodeShort(fuel.GetCH4(), 1200, 300, 0.99, 0.999)
	var inputState = states.NewGasPortState(gases.GetAir(), 500, 6e5, 1)
	bn.GasInput().SetState(inputState)

	bn.Process()

	assert.NotNil(t, bn)
	fmt.Println(bn.GasOutput().GetState(), bn.alpha)
}
