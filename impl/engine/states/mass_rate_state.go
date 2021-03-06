package states

import (
	"encoding/json"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
)

type MassRatePortState struct {
	MassRate float64 `json:"mass_rate"`
}

func NewMassRatePortState(massRate float64) MassRatePortState {
	return MassRatePortState{MassRate: massRate}
}

func (state MassRatePortState) Value() interface{} {
	return state.MassRate
}

func (state MassRatePortState) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		MassRateRel float64 `json:"mass_rate"`
	}{
		MassRateRel: state.MassRate,
	})
}

func (state MassRatePortState) Mix(another graph.PortState, relaxCoef float64) (graph.PortState, error) {
	switch v := another.(type) {
	case MassRatePortState:
		var casted = another.(MassRatePortState)

		return NewMassRatePortState(
			common.Lerp(state.MassRate, casted.MassRate, relaxCoef),
		), nil
	default:
		return nil, common.GetTypeError("MassRatePortState", v)
	}
}

func (state MassRatePortState) MaxResidual(another graph.PortState) (float64, error) {
	switch v := another.(type) {
	case MassRatePortState:
		return common.GetRelResidual(state.MassRate, another.(MassRatePortState).MassRate), nil
	default:
		return 0, common.GetTypeError("MassRatePortState", v)
	}
}
