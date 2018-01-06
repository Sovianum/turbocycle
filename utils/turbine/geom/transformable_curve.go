package geom

import "gonum.org/v1/gonum/mat"

type TransformableCurve interface {
	Curve
	Transformable
}

func NewTransformableCurve(curve Curve) TransformableCurve {
	return &transformableCurve{
		curve:      curve,
		combinator: NewCombinator(),
	}
}

type transformableCurve struct {
	curve      Curve
	combinator Combinator
}

func (tc *transformableCurve) Point(t float64) *mat.VecDense {
	return tc.combinator.Build().Apply(tc.curve.Point(t))
}

func (tc *transformableCurve) Transform(t Transformation) {
	tc.combinator.Combine(t)
}
