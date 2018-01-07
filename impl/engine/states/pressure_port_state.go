package states

import (
	"encoding/json"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
)

type PressurePortState struct {
	PStag float64
}

func NewPressurePortState(pStag float64) PressurePortState {
	return PressurePortState{PStag: pStag}
}

func (state PressurePortState) Value() interface{} {
	return state.PStag
}

func (state PressurePortState) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PStag float64 `json:"p_stag"`
	}{
		PStag: state.PStag,
	})
}

func (state PressurePortState) Mix(another graph.PortState, relaxCoef float64) (graph.PortState, error) {
	switch v := another.(type) {
	case PressurePortState:
		var casted = another.(PressurePortState)

		return NewPowerPortState(
			common.Lerp(state.PStag, casted.PStag, relaxCoef),
		), nil
	default:
		return nil, common.GetTypeError("PressurePortState", v)
	}
}

func (state PressurePortState) MaxResidual(another graph.PortState) (float64, error) {
	switch v := another.(type) {
	case PressurePortState:
		return common.GetRelResidual(state.PStag, another.(PressurePortState).PStag), nil
	default:
		return 0, common.GetTypeError("PressurePortState", v)
	}
}
