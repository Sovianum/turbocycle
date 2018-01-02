package graph

import (
	"github.com/Sovianum/turbocycle/common"
)

func NewTestPortState(val float64) TestPortState {
	return TestPortState{val: val}
}

type TestPortState struct {
	val float64
}

func (state TestPortState) Mix(another PortState, relaxCoef float64) (PortState, error) {
	return TestPortState{
		val: common.Lerp(state.val, another.(TestPortState).val, relaxCoef),
	}, nil
}

func (state TestPortState) MaxResidual(another PortState) (float64, error) {
	return common.GetRelResidual(state.val, another.(TestPortState).val), nil
}
