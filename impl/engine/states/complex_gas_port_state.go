package states

import (
	"encoding/json"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/material/gases"
	"math"
)

type ComplexGasPortState struct {
	Gas         gases.Gas
	TStag       float64
	PStag       float64
	MassRateRel float64
}

func NewComplexGasPortState(gas gases.Gas, tStag float64, pStag float64, massRateRel float64) ComplexGasPortState {
	return ComplexGasPortState{
		Gas:         gas,
		TStag:       tStag,
		PStag:       pStag,
		MassRateRel: massRateRel,
	}
}

func (state ComplexGasPortState) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PStag       float64 `json:"p_stag"`
		TStag       float64 `json:"t_stag"`
		MassRateRel float64 `json:"mass_rate_rel"`
	}{
		PStag:       state.PStag,
		TStag:       state.TStag,
		MassRateRel: state.MassRateRel,
	})
}

func StandardAtmosphereState() ComplexGasPortState {
	return NewComplexGasPortState(gases.GetAir(), 288, 1e5, 1) // TODO remove hardcoded constants
}

func (state ComplexGasPortState) Mix(another core.PortState, relaxCoef float64) (core.PortState, error) {
	switch v := another.(type) {
	case ComplexGasPortState:
		var casted = another.(ComplexGasPortState)

		return NewComplexGasPortState(
			state.Gas, // TODO add gas check
			common.Lerp(state.TStag, casted.TStag, relaxCoef),
			common.Lerp(state.PStag, casted.PStag, relaxCoef),
			common.Lerp(state.MassRateRel, casted.MassRateRel, relaxCoef),
		), nil
	default:
		return nil, common.GetTypeError("ComplexGasPortState", v)
	}
}

func (state ComplexGasPortState) MaxResidual(another core.PortState) (float64, error) {
	switch v := another.(type) {
	case ComplexGasPortState:
		var casted = another.(ComplexGasPortState)

		var res float64 = 0
		res = math.Max(res, common.GetRelResidual(state.TStag, casted.TStag))
		res = math.Max(res, common.GetRelResidual(state.PStag, casted.PStag))
		res = math.Max(res, common.GetRelResidual(state.MassRateRel, casted.MassRateRel))

		return res, nil

	default:
		return 0, common.GetTypeError("ComplexGasPortState", v)
	}
}
