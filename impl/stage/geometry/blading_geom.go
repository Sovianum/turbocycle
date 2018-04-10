package geometry

import (
	"encoding/json"
	"math"
)

func NewBladingGeometry(
	bladeWidth, gapWidth float64,
	innerProfile, outerProfile AxialProfileLine,
	midLineInterpFactor float64,
) BladingGeometry {
	return &bladingGeometry{
		bladeWidth:   bladeWidth,
		gapWidth:     gapWidth,
		innerProfile: innerProfile,
		outerProfile: outerProfile,
		meanProfile:  MeanLine(innerProfile, outerProfile, midLineInterpFactor),
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
	dOut := geom.OuterProfile().Diameter(x)
	dIn := geom.InnerProfile().Diameter(x)
	return math.Pi / 4 * (dOut*dOut - dIn*dIn)
}

func AxisDistance(x, hRel float64, geom BladingGeometry) float64 {
	var dIn = geom.InnerProfile().Diameter(x)
	var length = Height(x, geom)
	return dIn/2 + length*hRel
}

func RRel(dRel float64) float64 {
	return math.Sqrt((1 + dRel) / 2)
}

func DRel(x float64, geom BladingGeometry) float64 {
	return geom.InnerProfile().Diameter(x) / geom.OuterProfile().Diameter(x)
}

type bladingGeometry struct {
	bladeWidth   float64
	gapWidth     float64
	innerProfile AxialProfileLine
	meanProfile  AxialProfileLine
	outerProfile AxialProfileLine
}

func (geom *bladingGeometry) MarshalJSON() ([]byte, error) {
	data := struct {
		DInIn   float64 `json:"d_in_in"`
		DMeanIn float64 `json:"d_mean_in"`
		DOutIn  float64 `json:"d_out_in"`

		DInOut   float64 `json:"d_in_out"`
		DMeanOut float64 `json:"d_mean_out"`
		DOutOut  float64 `json:"d_out_out"`

		GammaIn  float64 `json:"gamma_in"`
		GammaOut float64 `json:"gamma_out"`

		Width      float64 `json:"width"`
		BladeWidth float64 `json:"blade_width"`
		GapWidth   float64 `json:"gap_width"`

		BladeHeight float64 `json:"blade_height"`
	}{
		DInIn:   geom.innerProfile.Diameter(0),
		DMeanIn: geom.meanProfile.Diameter(0),
		DOutIn:  geom.outerProfile.Diameter(0),

		DInOut:   geom.innerProfile.Diameter(geom.XGapOut()),
		DMeanOut: geom.meanProfile.Diameter(geom.XGapOut()),
		DOutOut:  geom.outerProfile.Diameter(geom.XGapOut()),

		GammaIn:  geom.innerProfile.Angle(),
		GammaOut: geom.outerProfile.Angle(),

		Width:      geom.XGapOut(),
		BladeWidth: geom.bladeWidth,
		GapWidth:   geom.gapWidth,

		BladeHeight: Height(0, geom),
	}

	return json.Marshal(data)
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
