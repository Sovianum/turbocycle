package states

import (
	"encoding/json"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
)

type TemperaturePortState struct {
	TStag float64
}

func NewTemperaturePortState(tStag float64) TemperaturePortState {
	return TemperaturePortState{TStag: tStag}
}

func (state TemperaturePortState) Value() interface{} {
	return state.TStag
}

func (state TemperaturePortState) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TStag float64 `json:"t_stag"`
	}{
		TStag: state.TStag,
	})
}

func (state TemperaturePortState) Mix(another graph.PortState, relaxCoef float64) (graph.PortState, error) {
	switch v := another.(type) {
	case TemperaturePortState:
		var casted = another.(TemperaturePortState)

		return NewPowerPortState(
			common.Lerp(state.TStag, casted.TStag, relaxCoef),
		), nil
	default:
		return nil, common.GetTypeError("TemperaturePortState", v)
	}
}

func (state TemperaturePortState) MaxResidual(another graph.PortState) (float64, error) {
	switch v := another.(type) {
	case TemperaturePortState:
		return common.GetRelResidual(state.TStag, another.(TemperaturePortState).TStag), nil
	default:
		return 0, common.GetTypeError("TemperaturePortState", v)
	}
}
