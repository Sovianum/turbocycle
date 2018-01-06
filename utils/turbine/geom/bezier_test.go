package geom

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

func TestInterpPoint(t *testing.T) {
	var testCases = []struct {
		point1        *mat.VecDense
		point2        *mat.VecDense
		interpFactor  float64
		expectedPoint *mat.VecDense
	}{
		{
			point1:        mat.NewVecDense(2, []float64{0, 0}),
			point2:        mat.NewVecDense(2, []float64{1, 1}),
			interpFactor:  0,
			expectedPoint: mat.NewVecDense(2, []float64{0, 0}),
		},
		{
			point1:        mat.NewVecDense(2, []float64{0, 0}),
			point2:        mat.NewVecDense(2, []float64{1, 1}),
			interpFactor:  1,
			expectedPoint: mat.NewVecDense(2, []float64{1, 1}),
		},
		{
			point1:        mat.NewVecDense(2, []float64{0, 0}),
			point2:        mat.NewVecDense(2, []float64{1, 1}),
			interpFactor:  0.5,
			expectedPoint: mat.NewVecDense(2, []float64{0.5, 0.5}),
		},
	}

	for i, tc := range testCases {
		var point = interpPoint(tc.point1, tc.point2, tc.interpFactor)
		assert.True(
			t,
			mat.EqualApprox(point, tc.expectedPoint, 1e-8),
			testMessage(i, tc.expectedPoint, point),
		)
	}
}

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

type BezierTestSuite struct {
	suite.Suite
	curve  *bezierCurve
	points []*mat.VecDense
}

func (suite *BezierTestSuite) SetupTest() {
	suite.points = []*mat.VecDense{
		mat.NewVecDense(2, []float64{0, 0}),
		mat.NewVecDense(2, []float64{1, 1}),
		mat.NewVecDense(2, []float64{2, 0}),
	}
	suite.curve = NewBezier(suite.points).(*bezierCurve)
}

func (suite *BezierTestSuite) TestPoint() {
	var testCases = []struct {
		t        float64
		expected *mat.VecDense
	}{
		{0, suite.points[0]},
		{1, suite.points[2]},
	}

	for i, tc := range testCases {
		assert.True(
			suite.T(),
			mat.EqualApprox(suite.curve.Point(tc.t), tc.expected, 1e-6),
			testMessage(i, suite.curve.Point(tc.t), tc.expected),
		)
	}
}

func (suite *BezierTestSuite) TestTermBasis() {
	var testCases = []struct {
		i        int
		t        float64
		expected float64
	}{
		{0, 0, 1},
		{0, 1, 0},
	}

	for i, tc := range testCases {
		assert.Equal(
			suite.T(), tc.expected, suite.curve.termBasis(tc.t, tc.i),
			testMessage(i, tc.expected, suite.curve.termBasis(tc.t, tc.i)),
		)
	}
}

func (suite *BezierTestSuite) TestTermFactor() {
	var testCases = []struct {
		i        int
		expected float64
	}{
		{0, 1},
		{1, 2},
		{2, 1},
	}

	for i, tc := range testCases {
		assert.Equal(
			suite.T(), tc.expected, suite.curve.termFactor(tc.i),
			testMessage(i, tc.expected, suite.curve.termFactor(tc.i)),
		)
	}
}

func (suite *BezierTestSuite) TestOrder() {
	assert.Equal(suite.T(), 2, suite.curve.order())
}

func TestBezierTestSuite(t *testing.T) {
	suite.Run(t, new(BezierTestSuite))
}

func testMessage(index interface{}, expected interface{}, got interface{}) string {
	return fmt.Sprintf("Failed %v: expected %v, got %v", index, expected, got)
}
