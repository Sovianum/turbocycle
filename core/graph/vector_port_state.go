package graph

import (
	"fmt"

	"github.com/Sovianum/turbocycle/common"
	"gonum.org/v1/gonum/mat"
)

func NewVectorPortState(vec *mat.VecDense) VectorPortState {
	return VectorPortState{Vec: vec}
}

type VectorPortState struct {
	Vec *mat.VecDense
}

func (state VectorPortState) Value() interface{} {
	return state.Vec
}

func (state VectorPortState) Mix(another PortState, relaxCoef float64) (PortState, error) {
	switch v := another.(type) {
	case VectorPortState:
		var casted = another.(VectorPortState)

		if l1, l2 := state.Vec.Len(), casted.Vec.Len(); l1 != l2 {
			return nil, fmt.Errorf("lengths %d and %d do not match", l1, l2)
		}

		var result = mat.NewVecDense(state.Vec.Len(), nil)
		result.AddScaledVec(result, relaxCoef, state.Vec)
		result.AddScaledVec(result, 1-relaxCoef, casted.Vec)

		return NewVectorPortState(result), nil
	default:
		return nil, common.GetTypeError("VectorPortState", v)
	}
}

func (state VectorPortState) MaxResidual(another PortState) (float64, error) {
	switch v := another.(type) {
	case VectorPortState:
		var casted = another.(VectorPortState)

		var diff = mat.NewVecDense(state.Vec.Len(), nil)
		diff.AddScaledVec(state.Vec, -1, casted.Vec)
		res := mat.Norm(diff, 2)
		return res, nil

	default:
		return 0, common.GetTypeError("VectorPortState", v)
	}
}
