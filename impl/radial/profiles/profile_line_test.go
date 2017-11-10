package profiles

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

func TestIntersectionPoint(t *testing.T) {
	var testCases = []struct {
		point1   *mat.VecDense
		point2   *mat.VecDense
		angle1   float64
		angle2   float64
		expected *mat.VecDense
	}{
		{
			point1:   mat.NewVecDense(2, []float64{0, 0}),
			point2:   mat.NewVecDense(2, []float64{1, 0}),
			angle1:   math.Pi / 4,
			angle2:   3 * math.Pi / 4,
			expected: mat.NewVecDense(2, []float64{0.5, 0.5}),
		},
	}

	for i, tc := range testCases {
		var point = intersectionPoint(
			tc.point1, tc.point2,
			tc.angle1, tc.angle2,
		)
		assert.True(
			t,
			mat.EqualApprox(point, tc.expected, 1e-8),
			testMessage(i, tc.expected, point),
		)
	}
}

type LineTestSuite struct {
	suite.Suite
	inletPoint  *mat.VecDense
	outletPoint *mat.VecDense
	inletAngle  float64
	outletAngle float64
	line Line
}

func (suite *LineTestSuite) SetupTest() {
	suite.inletPoint = mat.NewVecDense(2, []float64{0, 0})
	suite.outletPoint = mat.NewVecDense(2, []float64{1, 0})
	suite.inletAngle = math.Pi / 4
	suite.outletAngle = 3 * math.Pi / 4
	suite.line = NewLine(
		suite.inletPoint, suite.outletPoint, suite.inletAngle, suite.outletAngle,
	)
}

func (suite *LineTestSuite) TestGetPoints() {
	var points = suite.line.GetPoints(3)

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

func testMessage(index interface{}, expected interface{}, got interface{}) string {
	return fmt.Sprintf("Failed %v: expected %v, got %v", index, expected, got)
}
