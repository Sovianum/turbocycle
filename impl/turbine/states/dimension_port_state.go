package states

import (
	"encoding/json"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type DimensionPortState struct {
	Dimension float64
}

func NewDiameterPortState(diameter float64) DimensionPortState {
	return DimensionPortState{Dimension: diameter}
}

func (state DimensionPortState) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Diameter float64 `json:"diameter"`
	}{Diameter: state.Dimension})
}

func (state DimensionPortState) Mix(another graph.PortState, relaxCoef float64) (graph.PortState, error) {
	switch v := another.(type) {
	case DimensionPortState:
		var casted = another.(DimensionPortState)
		return states.NewMassRatePortState(common.Lerp(state.Dimension, casted.Dimension, relaxCoef)), nil
	default:
		return nil, common.GetTypeError("DimensionPortState", v)
	}
}

func (state DimensionPortState) MaxResidual(another graph.PortState) (float64, error) {
	switch v := another.(type) {
	case DimensionPortState:
		return common.GetRelResidual(state.Dimension, another.(DimensionPortState).Dimension), nil
	default:
		return 0, common.GetTypeError("DimensionPortState", v)
	}
}
