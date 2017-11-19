package geom

import "gonum.org/v1/gonum/mat"

func NewCombinator() Combinator {
	return &combinator{make([]Transformation, 0)}
}

type Combinator interface {
	Combine(t Transformation) Combinator
	Build() Transformation
}

type combinator struct {
	transformations []Transformation
}

func (c *combinator) Combine(t Transformation) Combinator {
	c.transformations = append(c.transformations, t)
	return c
}

func (c *combinator) Build() Transformation {
	var transFunc = func(vec *mat.VecDense) *mat.VecDense {
		var result = vec
		for _, trans := range c.transformations {
			result = trans.Apply(result)
		}
		return result
	}
	return NewGenericTransformation(transFunc)
}
