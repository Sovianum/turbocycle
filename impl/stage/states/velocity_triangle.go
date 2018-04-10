package states

import (
	"encoding/json"
	"math"
)

type VelocityTriangle interface {
	C() float64
	W() float64
	U() float64
	Alpha() float64
	Beta() float64
	CA() float64
	WA() float64
	CU() float64
	WU() float64
}

func MarshalTriangle(triangle VelocityTriangle) ([]byte, error) {
	return json.Marshal(struct {
		C     float64 `json:"c"`
		CU    float64 `json:"cu"`
		CA    float64 `json:"ca"`
		W     float64 `json:"w"`
		WU    float64 `json:"wu"`
		WA    float64 `json:"wa"`
		U     float64 `json:"u"`
		Alpha float64 `json:"alpha"`
		Beta  float64 `json:"beta"`
	}{
		C:     triangle.C(),
		CU:    triangle.CU(),
		CA:    triangle.CA(),
		W:     triangle.W(),
		WU:    triangle.WU(),
		WA:    triangle.WA(),
		U:     triangle.U(),
		Alpha: triangle.Alpha(),
		Beta:  triangle.Beta(),
	})
}

func NewInletTriangle(u, c, alpha float64) VelocityTriangle {
	return rotorInletVelocityTriangle{
		velocityTriangle{u: u, c: c, alpha: alpha},
	}
}

func NewInletTriangleFromProjections(cu, ca, u float64) VelocityTriangle {
	var c = math.Sqrt(cu*cu + ca*ca)
	var alpha = math.Atan2(ca, cu)
	return rotorInletVelocityTriangle{
		velocityTriangle{u: u, c: c, alpha: alpha},
	}
}

func NewOutletTriangle(u, c, alpha float64) VelocityTriangle {
	return rotorOutletVelocityTriangle{
		velocityTriangle{u: u, c: c, alpha: alpha},
	}
}

func NewOutletTriangleFromProjections(cu, ca, u float64) VelocityTriangle {
	var c = math.Sqrt(cu*cu + ca*ca)
	var alpha = math.Atan2(ca, cu)
	return rotorOutletVelocityTriangle{
		velocityTriangle{u: u, c: c, alpha: alpha},
	}
}

type rotorOutletVelocityTriangle struct {
	velocityTriangle
}

func (triangle rotorOutletVelocityTriangle) MarshalJSON() ([]byte, error) {
	return MarshalTriangle(triangle)
}

func (triangle rotorOutletVelocityTriangle) WU() float64 {
	return triangle.CU() + triangle.U()
}

func (triangle rotorOutletVelocityTriangle) Beta() float64 {
	return math.Atan2(triangle.WA(), triangle.WU())
}

func (triangle rotorOutletVelocityTriangle) W() float64 {
	var term1 = triangle.c * triangle.c
	var term2 = triangle.u * triangle.u
	var term3 = 2 * triangle.c * triangle.u * math.Cos(triangle.alpha)

	return math.Sqrt(term1 + term2 + term3)
}

type rotorInletVelocityTriangle struct {
	velocityTriangle
}

func (triangle rotorInletVelocityTriangle) MarshalJSON() ([]byte, error) {
	return MarshalTriangle(triangle)
}

func (triangle rotorInletVelocityTriangle) WU() float64 {
	return triangle.CU() - triangle.U()
}

func (triangle rotorInletVelocityTriangle) Beta() float64 {
	return math.Atan2(triangle.WA(), triangle.WU())
}

func (triangle rotorInletVelocityTriangle) W() float64 {
	var term1 = triangle.c * triangle.c
	var term2 = triangle.u * triangle.u
	var term3 = -2 * triangle.c * triangle.u * math.Cos(triangle.alpha)

	return math.Sqrt(term1 + term2 + term3)
}

func NewCompressorVelocityTriangle(c, u, alpha float64) VelocityTriangle {
	return &compressorVelocityTriangle{
		velocityTriangle{
			c:     c,
			u:     u,
			alpha: alpha,
		},
	}
}

func NewCompressorVelocityTriangleFromProjections(cu, ca, u float64) VelocityTriangle {
	c := math.Sqrt(cu*cu + ca*ca)
	alpha := math.Atan2(ca, cu)
	return &compressorVelocityTriangle{
		velocityTriangle{
			c:     c,
			u:     u,
			alpha: alpha,
		},
	}
}

type compressorVelocityTriangle struct {
	velocityTriangle
}

func (triangle *compressorVelocityTriangle) MarshalJSON() ([]byte, error) {
	return MarshalTriangle(triangle)
}

func (triangle *compressorVelocityTriangle) W() float64 {
	wu := triangle.WU()
	wa := triangle.CA()
	return math.Sqrt(wu*wu + wa*wa)
}

func (triangle *compressorVelocityTriangle) Beta() float64 {
	wu := triangle.WU()
	wa := triangle.CA()
	return math.Atan2(wa, wu)
}

func (triangle *compressorVelocityTriangle) WU() float64 {
	return triangle.u - triangle.c*math.Cos(triangle.alpha)
}

type velocityTriangle struct {
	c     float64
	u     float64
	alpha float64
}

func (triangle velocityTriangle) Alpha() float64 {
	return triangle.alpha
}

func (triangle velocityTriangle) C() float64 {
	return triangle.c
}

func (triangle velocityTriangle) U() float64 {
	return triangle.u
}

func (triangle velocityTriangle) CA() float64 {
	return triangle.c * math.Sin(triangle.alpha)
}

func (triangle velocityTriangle) CU() float64 {
	return triangle.c * math.Cos(triangle.alpha)
}

func (triangle velocityTriangle) WA() float64 {
	return triangle.c * math.Sin(triangle.alpha)
}
