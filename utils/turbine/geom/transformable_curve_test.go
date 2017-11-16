package geom

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

type LineTestSuite struct {
	suite.Suite
	inletPoint  *mat.VecDense
	outletPoint *mat.VecDense
	inletAngle  float64
	outletAngle float64
	line        TransformableCurve
}

func (suite *LineTestSuite) SetupTest() {
	suite.inletPoint = mat.NewVecDense(2, []float64{0, 0})
	suite.outletPoint = mat.NewVecDense(2, []float64{1, 0})
	suite.inletAngle = math.Pi / 4
	suite.outletAngle = 3 * math.Pi / 4
	suite.line = NewTransformableCurve(NewBezier2FromOrientedPoints(
		suite.inletPoint, suite.outletPoint, suite.inletAngle, suite.outletAngle,
	))
}

func (suite *LineTestSuite) TestGetPoints() {
	var points = LinPoints(suite.line, 0, 1, 3)

	assert.True(
		suite.T(),
		mat.EqualApprox(points[0], suite.inletPoint, 1e-8),
		testMessage(0, suite.inletPoint, points[0]),
	)

	assert.True(
		suite.T(),
		mat.EqualApprox(points[2], suite.outletPoint, 1e-8),
		testMessage(0, suite.outletPoint, points[0]),
	)

	assert.Equal(
		suite.T(),
		points[1].At(0, 0),
		0.5,
		testMessage(0, 0.5, points[1].At(0, 0)),
	)
}

func TestLineTestSuite(t *testing.T) {
	suite.Run(t, new(LineTestSuite))
}
