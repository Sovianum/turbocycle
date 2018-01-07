package states

import (
	"encoding/json"

	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/material/gases"
)

type GasPortState struct {
	Gas gases.Gas
}

func NewGasPortState(gas gases.Gas) GasPortState {
	return GasPortState{
		Gas: gas,
	}
}

func (state GasPortState) Value() interface{} {
	return state.Gas
}

func (state GasPortState) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct{}{})
}

func (state GasPortState) Mix(another graph.PortState, relaxCoef float64) (graph.PortState, error) {
	return state, nil // TODO add gas check
}

func (state GasPortState) MaxResidual(another graph.PortState) (float64, error) {
	return 0, nil // TODO add gas check
}
