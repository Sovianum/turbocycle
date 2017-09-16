package states

import (
	"encoding/json"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
)

type MassRateRelPortState struct {
	MassRateRel float64 `json:"mass_rate_rel"`
}

func NewMassRateRelPortState(massRateRel float64) MassRateRelPortState {
	return MassRateRelPortState{MassRateRel: massRateRel}
}

func (state MassRateRelPortState) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		MassRateRel float64 `json:"mass_rate_rel"`
	}{
		MassRateRel: state.MassRateRel,
	})
}

func (state MassRateRelPortState) Mix(another core.PortState, relaxCoef float64) (core.PortState, error) {
	switch v := another.(type) {
	case MassRateRelPortState:
		var casted = another.(MassRateRelPortState)

		return NewMassRateRelPortState(
			common.Lerp(state.MassRateRel, casted.MassRateRel, relaxCoef),
		), nil
	default:
		return nil, common.GetTypeError("MassRateRelPortState", v)
	}
}

func (state MassRateRelPortState) MaxResidual(another core.PortState) (float64, error) {
	switch v := another.(type) {
	case MassRateRelPortState:
		return common.GetRelResidual(state.MassRateRel, another.(MassRateRelPortState).MassRateRel), nil
	default:
		return 0, common.GetTypeError("MassRateRelPortState", v)
	}
}
