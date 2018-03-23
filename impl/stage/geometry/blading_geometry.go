package geometry

import "math"

func NewBladingGeometry(bladeWidth, gapWidth float64, innerProfile, outerProfile AxialProfileLine) BladingGeometry {
	return &bladingGeometry{
		bladeWidth:   bladeWidth,
		gapWidth:     gapWidth,
		innerProfile: innerProfile,
		outerProfile: outerProfile,
		meanProfile:  MeanLine(innerProfile, outerProfile, 0.5),
	}
}

type BladingGeometry interface {
	XBladeIn() float64
	XBladeOut() float64
	XGapOut() float64
	InnerProfile() AxialProfileLine
	MeanProfile() AxialProfileLine
	OuterProfile() AxialProfileLine
}

func ExpansionAngle(geom BladingGeometry) float64 {
	return geom.OuterProfile().Angle() - geom.InnerProfile().Angle()
}

func Height(x float64, geom BladingGeometry) float64 {
	return (geom.OuterProfile().Diameter(x) - geom.InnerProfile().Diameter(x)) / 2
}

func RelativeHeight(x float64, geom BladingGeometry) float64 {
	var diameter = geom.MeanProfile().Diameter(x)
	var height = Height(x, geom)

	return height / diameter
}

func ChordProjection(geom BladingGeometry) float64 {
	return geom.XBladeOut() - geom.XBladeIn()
}

func Elongation(x float64, geom BladingGeometry) float64 {
	return Height(x, geom) / ChordProjection(geom)
}

func AxialGapProjection(geom BladingGeometry) float64 {
	return geom.XGapOut() - geom.XBladeOut()
}

func Area(x float64, geom BladingGeometry) float64 {
	return math.Pi * geom.MeanProfile().Diameter(x) * Height(x, geom)
}

func AxisDistance(x, hRel float64, geom BladingGeometry) float64 {
	var dIn = geom.InnerProfile().Diameter(x)
	var length = Height(x, geom)
	return dIn/2 + length*hRel
}

type bladingGeometry struct {
	bladeWidth   float64
	gapWidth     float64
	innerProfile AxialProfileLine
	meanProfile  AxialProfileLine
	outerProfile AxialProfileLine
}

func (geom *bladingGeometry) InnerProfile() AxialProfileLine {
	return geom.innerProfile
}

func (geom *bladingGeometry) MeanProfile() AxialProfileLine {
	return geom.meanProfile
}

func (geom *bladingGeometry) OuterProfile() AxialProfileLine {
	return geom.outerProfile
}

func (geom *bladingGeometry) XBladeIn() float64 {
	return 0
}

func (geom *bladingGeometry) XBladeOut() float64 {
	return geom.bladeWidth
}

func (geom *bladingGeometry) XGapOut() float64 {
	return geom.bladeWidth + geom.gapWidth
}
