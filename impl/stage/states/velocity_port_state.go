package states

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
)

const (
	InletTriangleType      = "InletTriangleType"
	OutletTriangleType     = "OutletTriangleType"
	CompressorTriangleType = "CompressorTriangleType"
	VelocityInletTag       = "VelocityInletTag"
	VelocityOutletTag      = "VelocityOutletTag"
)

type VelocityPortState struct {
	Triangle     VelocityTriangle
	TriangleType string
}

func NewVelocityPortState(triangle VelocityTriangle, triangleType string) VelocityPortState {
	return VelocityPortState{
		Triangle:     triangle,
		TriangleType: triangleType,
	}
}

func (state VelocityPortState) Value() interface{} {
	return state.Triangle
}

func (state VelocityPortState) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Triangle VelocityTriangle
	}{
		Triangle: state.Triangle,
	})
}

func (state VelocityPortState) Mix(another graph.PortState, relaxCoef float64) (graph.PortState, error) {
	switch v := another.(type) {
	case VelocityPortState:
		var casted = another.(VelocityPortState)
		if casted.TriangleType != state.TriangleType {
			return nil, fmt.Errorf("incompatible triangle types %s and %s", state.TriangleType, casted.TriangleType)
		}

		var triangle1 = state.Triangle
		var triangle2 = casted.Triangle

		var u = common.Lerp(triangle1.U(), triangle2.U(), relaxCoef)
		var c = common.Lerp(triangle1.C(), triangle2.C(), relaxCoef)
		var alpha = common.Lerp(triangle1.Alpha(), triangle2.Alpha(), relaxCoef)
		var triangle VelocityTriangle

		switch state.TriangleType {
		case InletTriangleType:
			triangle = NewInletTriangle(u, c, alpha)
		case OutletTriangleType:
			triangle = NewOutletTriangle(u, c, alpha)
		default:
			return nil, fmt.Errorf("incorrect type %s", state.TriangleType)
		}

		return NewVelocityPortState(triangle, state.TriangleType), nil
	default:
		return nil, common.GetTypeError("MassRatePortState", v)
	}
}

func (state VelocityPortState) MaxResidual(another graph.PortState) (float64, error) {
	switch v := another.(type) {
	case VelocityPortState:
		var casted = another.(VelocityPortState)

		var res float64 = 0
		res = math.Max(res, common.GetRelResidual(state.Triangle.C(), casted.Triangle.C()))
		res = math.Max(res, common.GetRelResidual(state.Triangle.U(), casted.Triangle.U()))
		res = math.Max(res, common.GetRelResidual(state.Triangle.Alpha(), casted.Triangle.Alpha()))

		return res, nil
	default:
		return 0, common.GetTypeError("MassRatePortState", v)
	}
}
