package states

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
)

func NewRPMPortState(rpm float64) RPMPortState {
	return RPMPortState{RPM: rpm}
}

type RPMPortState struct {
	RPM float64
}

func (state RPMPortState) Value() interface{} {
	return state.RPM
}

func (state RPMPortState) Mix(another graph.PortState, relaxCoef float64) (graph.PortState, error) {
	switch v := another.(type) {
	case PowerPortState:
		var casted = another.(RPMPortState)

		return NewPowerPortState(
			common.Lerp(state.RPM, casted.RPM, relaxCoef),
		), nil
	default:
		return nil, common.GetTypeError("RPMPortState", v)
	}
}

func (state RPMPortState) MaxResidual(another graph.PortState) (float64, error) {
	switch v := another.(type) {
	case PowerPortState:
		return common.GetRelResidual(state.RPM, another.(RPMPortState).RPM), nil
	default:
		return 0, common.GetTypeError("RPMPortState", v)
	}
}
