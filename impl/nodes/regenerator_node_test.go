package nodes

import (
	"testing"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/stretchr/testify/assert"
	"github.com/Sovianum/turbocycle/common"
	"fmt"
)

func TestRegeneratorNode_Process_ColdMode(t *testing.T) {
	var rn = NewRegeneratorNode(0.9, 0.05, SigmaByColdSide)
	var coldState = states.NewGasPortState(gases.GetAir(), 300, 1e5, 1)
	var hotState = states.NewGasPortState(gases.GetAir(), 800, 1e5, 1)

	rn.ColdInput().SetState(coldState)
	rn.HotInput().SetState(hotState)

	var err = rn.Process()
	assert.Nil(t, err)

	var tColdOut = rn.ColdOutput().GetState().(states.GasPortState).TStag
	assert.True(t, common.ApproxEqual(750, tColdOut, 0.01))

	var tHotOut = rn.HotOutput().GetState().(states.GasPortState).TStag
	assert.True(t,
		common.ApproxEqual(350, tHotOut, 0.03),
		fmt.Sprintf("Expected %f, got %f (precision %f)", 350., tHotOut, 0.01),
	)
}

func TestRegeneratorNode_Process_HotMode(t *testing.T) {
	var rn = NewRegeneratorNode(0.9, 0.05, SigmaByColdSide)
	var coldState = states.NewGasPortState(gases.GetAir(), 300, 1e5, 1)
	var hotState = states.NewGasPortState(gases.GetAir(), 800, 1e5, 1)

	rn.ColdInput().SetState(coldState)
	rn.HotInput().SetState(hotState)

	var err = rn.Process()
	assert.Nil(t, err)

	var tColdOut = rn.ColdOutput().GetState().(states.GasPortState).TStag
	assert.True(t, common.ApproxEqual(750, tColdOut, 0.01))

	var tHotOut = rn.HotOutput().GetState().(states.GasPortState).TStag
	assert.True(t,
		common.ApproxEqual(350, tHotOut, 0.03),
		fmt.Sprintf("Expected %f, got %f (precision %f)", 350., tHotOut, 0.01),
	)
}


