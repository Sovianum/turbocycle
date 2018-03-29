package states

import "math"

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
	return triangle.u - triangle.c*math.Sin(triangle.alpha)
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
