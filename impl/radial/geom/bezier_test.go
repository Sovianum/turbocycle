package geom

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

type BezierTestSuite struct {
	suite.Suite
	curve *bezierCurve
	points []*mat.VecDense
}

func (suite *BezierTestSuite) SetupTest() {
	suite.points = []*mat.VecDense{
		mat.NewVecDense(2, []float64{0, 0}),
		mat.NewVecDense(2, []float64{1, 1}),
		mat.NewVecDense(2, []float64{2, 0}),
	}
	suite.curve = NewBezierCurve(suite.points).(*bezierCurve)
}

func (suite *BezierTestSuite) TestPoint() {
	var testCases = []struct{
		t float64
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
	var testCases = []struct{
		i int
		t float64
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
	var testCases = []struct{
		i int
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
