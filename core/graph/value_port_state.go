package graph

import (
	"math"

	"github.com/Sovianum/turbocycle/common"
)

func NewNumberPortState(num float64) NumberPortState {
	return NumberPortState{Num: num}
}

type NumberPortState struct {
	Num float64
}

func (state NumberPortState) Mix(another PortState, relaxCoef float64) (PortState, error) {
	switch v := another.(type) {
	case NumberPortState:
		var casted = another.(NumberPortState)

		return NewNumberPortState(
			common.Lerp(state.Num, casted.Num, relaxCoef),
		), nil
	default:
		return nil, common.GetTypeError("NumberPortState", v)
	}
}

func (state NumberPortState) MaxResidual(another PortState) (float64, error) {
	switch v := another.(type) {
	case NumberPortState:
		var casted = another.(NumberPortState)

		var res float64 = 0
		res = math.Max(res, common.GetRelResidual(state.Num, casted.Num))

		return res, nil

	default:
		return 0, common.GetTypeError("NumberPortState", v)
	}
}
