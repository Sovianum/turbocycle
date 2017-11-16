package geom

import (
	"testing"

	"math"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

const (
	precision = 1e8
)

type TransformationTestSuite struct {
	suite.Suite
	vec *mat.VecDense
}

func (suite *TransformationTestSuite) SetupTest() {
	suite.vec = mat.NewVecDense(2, []float64{1, 0})
}

func (suite *TransformationTestSuite) TestScale() {
	var res = Scale(2).Apply(suite.vec)
	assert.True(
		suite.T(),
		mat.EqualApprox(res, mat.NewVecDense(2, []float64{2, 0}), precision),
		testMessage(0, mat.NewVecDense(2, []float64{2, 0}), res),
	)
}

func (suite *TransformationTestSuite) TestTranslation() {
	var res = Translation(mat.NewVecDense(2, []float64{0, 1})).Apply(suite.vec)
	assert.True(
		suite.T(),
		mat.EqualApprox(res, mat.NewVecDense(2, []float64{1, 1}), precision),
		testMessage(0, mat.NewVecDense(2, []float64{1, 1}), res),
	)
}

func (suite *TransformationTestSuite) TestReflection() {
	var res = Reflection(math.Pi / 4).Apply(suite.vec)
	assert.True(
		suite.T(),
		mat.EqualApprox(res, mat.NewVecDense(2, []float64{0, 1}), precision),
		testMessage(0, mat.NewVecDense(2, []float64{0, 1}), res),
	)
}

func (suite *TransformationTestSuite) TestRotation() {
	var res = Rotation(math.Pi / 2).Apply(suite.vec)
	assert.True(
		suite.T(),
		mat.EqualApprox(res, mat.NewVecDense(2, []float64{0, 1}), precision),
		testMessage(0, mat.NewVecDense(2, []float64{0, 1}), res),
	)
}

func TestTransformationTestSuite(t *testing.T) {
	suite.Run(t, new(TransformationTestSuite))
}
