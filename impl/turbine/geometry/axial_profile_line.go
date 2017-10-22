package geometry

import "math"

type AxialProfileLine interface {
	Diameter(x float64) float64
	Angle() float64
}

func NewAxialProfileLine(x0, d0, angle float64) AxialProfileLine {
	return &axialProfileLine{
		x0:    x0,
		d0:    d0,
		angle: angle,
	}
}

type axialProfileLine struct {
	x0    float64 // known value of axial coordinate
	d0    float64 // known value of diameter
	angle float64 // angle of this line
}

func (line *axialProfileLine) Diameter(x float64) float64 {
	return line.d0 + 2*math.Tan(line.angle)*(x-line.x0)
}

func (line *axialProfileLine) Angle() float64 {
	return line.angle
}