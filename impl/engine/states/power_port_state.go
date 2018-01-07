package states

import (
	"encoding/json"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
)

type PowerPortState struct {
	LSpecific float64 `json:"l_specific"`
}

func NewPowerPortState(lSpecific float64) PowerPortState {
	return PowerPortState{LSpecific: lSpecific}
}

func (state PowerPortState) Value() interface{} {
	return state.LSpecific
}

func StandardPowerState() PowerPortState {
	return NewPowerPortState(0)
}

func (state PowerPortState) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LSpecific float64 `json:"l_specific"`
	}{
		LSpecific: state.LSpecific,
	})
}

func (state PowerPortState) Mix(another graph.PortState, relaxCoef float64) (graph.PortState, error) {
	switch v := another.(type) {
	case PowerPortState:
		var casted = another.(PowerPortState)

		return NewPowerPortState(
			common.Lerp(state.LSpecific, casted.LSpecific, relaxCoef),
		), nil
	default:
		return nil, common.GetTypeError("PowerPortState", v)
	}
}

func (state PowerPortState) MaxResidual(another graph.PortState) (float64, error) {
	switch v := another.(type) {
	case PowerPortState:
		return common.GetRelResidual(state.LSpecific, another.(PowerPortState).LSpecific), nil
	default:
		return 0, common.GetTypeError("PowerPortState", v)
	}
}
