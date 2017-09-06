package states

import (
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/common"
)

type PowerPortState struct {
	LSpecific float64
}

func NewPowerPortState(lSpecific float64) PowerPortState {
	return PowerPortState{LSpecific: lSpecific}
}

func (state PowerPortState) Mix(another core.IPortState, relaxCoef float64) (core.IPortState, error) {
	switch v := another.(type) {
	case GasPortState:
		var casted = another.(PowerPortState)

		return NewPowerPortState(
			common.Lerp(state.LSpecific, casted.LSpecific, relaxCoef),
		), nil
	default:
		return nil, common.GetTypeError("PowerPortState", v)
	}
}

func (state PowerPortState) MaxResidual(another core.IPortState) (float64, error) {
	switch v := another.(type) {
	case PowerPortState:
		return common.GetRelResidual(state.LSpecific, another.(PowerPortState).LSpecific), nil

	default:
		return 0, common.GetTypeError("PowerPortState", v)
	}
}

