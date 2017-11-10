package geom

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

type Transformation interface {
	Apply(vec *mat.VecDense) *mat.VecDense
}

func ApplyAll(t Transformation, points []*mat.VecDense) []*mat.VecDense {
	var result = make([]*mat.VecDense, len(points))
	for i, point := range points {
		result[i] = t.Apply(point)
	}
	return result
}

func Unit() Transformation {
	return unit{}
}

type unit struct {}

func (unit) Apply(vec *mat.VecDense) *mat.VecDense {
	var result = newPoint()
	result.CloneVec(vec)
	return result
}

func Scale(factor float64) Transformation {
	return scale{factor: factor}
}

type scale struct {
	factor float64
}

func (t scale) Apply(vec *mat.VecDense) *mat.VecDense {
	var result = newPoint()
	result.ScaleVec(t.factor, vec)
	return result
}

func Translation(offset *mat.VecDense) Transformation {
	return translation{offset: offset}
}

type translation struct {
	offset *mat.VecDense
}

func (t translation) Apply(vec *mat.VecDense) *mat.VecDense {
	var result = newPoint()
	result.CloneVec(vec)
	result.AddVec(result, t.offset)
	return result
}

func Reflection(angle float64) Transformation {
	return reflection{
		reflectionMatrix: mat.NewDense(
			2, 2,
			[]float64{
				math.Cos(2 * angle), math.Sin(2 * angle),
				math.Sin(2 * angle), -math.Cos(2 * angle),
			},
		),
	}
}

type reflection struct {
	reflectionMatrix *mat.Dense
}

func (r reflection) Apply(vec *mat.VecDense) *mat.VecDense {
	var result = newPoint()
	result.MulVec(r.reflectionMatrix, vec)
	return (*mat.VecDense)(result)
}

func Rotation(angle float64) Transformation {
	return reflection{
		reflectionMatrix: mat.NewDense(
			2, 2,
			[]float64{
				math.Cos(angle), -math.Sin(angle),
				math.Sin(angle), -math.Cos(angle),
			},
		),
	}
}

type rotation struct {
	rotationMatrix *mat.Dense
}

func (r rotation) Apply(vec *mat.VecDense) *mat.VecDense {
	var result = newPoint()
	result.MulVec(r.rotationMatrix, vec)
	return (*mat.VecDense)(result)
}

func NewGenericTransformation(f func(*mat.VecDense) *mat.VecDense) Transformation {
	return genericTransformation{transFunc: f}
}

type genericTransformation struct {
	transFunc func(dense *mat.VecDense) *mat.VecDense
}

func (t genericTransformation) Apply(vec *mat.VecDense) *mat.VecDense {
	return t.transFunc(vec)
}
