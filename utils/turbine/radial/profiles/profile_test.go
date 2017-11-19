package profiles

import (
	"fmt"
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

func TestRadialPoint(t *testing.T) {
	var testCases = []struct{
		startPoint *mat.VecDense
		angle float64
		direction float64
		radius float64
		expectedPoint *mat.VecDense
	}{
		{
			startPoint:mat.NewVecDense(2, []float64{0, 0}),
			angle:0,
			direction:positive,
			radius:1,
			expectedPoint:mat.NewVecDense(2, []float64{0, 1}),
		},
		{
			startPoint:mat.NewVecDense(2, []float64{0, 0}),
			angle:0,
			direction:negative,
			radius:1,
			expectedPoint:mat.NewVecDense(2, []float64{0, -1}),
		},
		{
			startPoint:mat.NewVecDense(2, []float64{1, 1}),
			angle:math.Pi / 4,
			direction:positive,
			radius:1,
			expectedPoint:mat.NewVecDense(2, []float64{1 - math.Sqrt2 / 2, 1 + math.Sqrt2 / 2}),
		},
	}

	for i, tc := range testCases {
		var point = radialPoint(tc.startPoint, tc.angle, tc.radius, tc.direction)
		assert.InDelta(t, tc.expectedPoint.At(0, 0), point.At(0, 0), 1e-8, "x at i = %d", i)
		assert.InDelta(t, tc.expectedPoint.At(1, 0), point.At(1, 0), 1e-8, "y at i = %d", i)
	}
}

type BladeProfileTestSuite struct {
	suite.Suite

	inletPSPoint  *mat.VecDense
	outletPSPoint *mat.VecDense

	inletSSPoint  *mat.VecDense
	outletSSPoint *mat.VecDense

	inletMeanPoint  *mat.VecDense
	outletMeanPoint *mat.VecDense

	inletPSAngle  float64
	outletPSAngle float64

	inletSSAngle  float64
	outletSSAngle float64

	inletMeanAngle  float64
	outletMeanAngle float64

	profile BladeProfile
}

func (suite *BladeProfileTestSuite) SetupTest() {
	suite.inletPSPoint = mat.NewVecDense(2, []float64{0, 1})
	suite.outletPSPoint = mat.NewVecDense(2, []float64{1, 1})
	suite.inletPSAngle = math.Pi / 4
	suite.outletPSAngle = 3 * math.Pi / 4

	suite.inletSSPoint = mat.NewVecDense(2, []float64{0, -1})
	suite.outletSSPoint = mat.NewVecDense(2, []float64{1, -1})
	suite.inletSSAngle = 3 * math.Pi / 4
	suite.outletSSAngle = math.Pi / 4

	suite.inletMeanPoint = mat.NewVecDense(2, []float64{0, 0})
	suite.outletMeanPoint = mat.NewVecDense(2, []float64{1, 0})
	suite.inletMeanAngle = 0.0001
	suite.outletMeanAngle = math.Pi - 0.0001

	suite.profile = NewBladeProfile(
		suite.inletPSPoint, suite.outletPSPoint,
		suite.inletSSPoint, suite.outletSSPoint,
		suite.inletMeanPoint, suite.outletMeanPoint,
		suite.inletPSAngle, suite.outletPSAngle,
		suite.inletSSAngle, suite.outletSSAngle,
		suite.inletMeanAngle, suite.outletMeanAngle,
	)
}

func (suite *BladeProfileTestSuite) TestMeanLine() {
	var points = geom.LinPoints(suite.profile.MeanLine(), 0, 1, 3)

	assert.True(
		suite.T(),
		mat.EqualApprox(points[0], suite.inletMeanPoint, 1e-8),
		testMessage(0, suite.inletPSPoint, points[0]),
	)

	assert.True(
		suite.T(),
		mat.EqualApprox(points[2], suite.outletMeanPoint, 1e-8),
		testMessage(0, suite.inletPSPoint, points[2]),
	)

	assert.InDelta(
		suite.T(),
		0.5,
		points[1].At(0, 0),
		1e-8,
		testMessage(0, 0.5, points[1].At(0, 0)),
	)
}

func (suite *BladeProfileTestSuite) TestPSLine() {
	var points = geom.LinPoints(suite.profile.PSLine(), 0, 1, 3)

	assert.True(
		suite.T(),
		mat.EqualApprox(points[0], suite.inletPSPoint, 1e-8),
		testMessage(0, suite.inletPSPoint, points[0]),
	)

	assert.True(
		suite.T(),
		mat.EqualApprox(points[2], suite.outletPSPoint, 1e-8),
		testMessage(0, suite.inletPSPoint, points[2]),
	)

	assert.InDelta(
		suite.T(),
		0.5,
		points[1].At(0, 0),
		1e-8,
		testMessage(0, 0.5, points[1].At(0, 0)),
	)
}

func (suite *BladeProfileTestSuite) TestSSLine() {
	var points = geom.LinPoints(suite.profile.SSLine(), 0, 1, 3)

	assert.True(
		suite.T(),
		mat.EqualApprox(points[0], suite.inletSSPoint, 1e-8),
		testMessage(0, suite.inletPSPoint, points[0]),
	)

	assert.True(
		suite.T(),
		mat.EqualApprox(points[2], suite.outletSSPoint, 1e-8),
		testMessage(0, suite.inletPSPoint, points[2]),
	)

	assert.InDelta(
		suite.T(),
		0.5,
		points[1].At(0, 0),
		1e-8,
		testMessage(0, 0.5, points[1].At(0, 0)),
	)
}

func (suite *BladeProfileTestSuite) TestInletEdge() {
	var points = geom.LinPoints(suite.profile.InletEdge(), 0, 1, 3)

	assert.True(
		suite.T(),
		mat.EqualApprox(points[0], suite.inletPSPoint, 1e-8),
		testMessage(0, suite.inletPSPoint, points[0]),
	)

	assert.True(
		suite.T(),
		mat.EqualApprox(points[2], suite.inletSSPoint, 1e-8),
		testMessage(0, suite.inletPSPoint, points[2]),
	)

	assert.InDelta(
		suite.T(),
		0,
		points[1].At(1, 0),
		1e-8,
		testMessage(0, 0, points[1].At(1, 0)),
	)
}

func (suite *BladeProfileTestSuite) TestOutletEdge() {
	var points = geom.LinPoints(suite.profile.OutletEdge(), 0, 1, 3)

	assert.True(
		suite.T(),
		mat.EqualApprox(points[0], suite.outletPSPoint, 1e-8),
		testMessage(0, suite.outletPSPoint, points[0]),
	)

	assert.True(
		suite.T(),
		mat.EqualApprox(points[2], suite.outletSSPoint, 1e-8),
		testMessage(0, suite.outletPSPoint, points[2]),
	)

	assert.InDelta(
		suite.T(),
		0,
		points[1].At(1, 0),
		1e-8,
		testMessage(0, 0, points[1].At(1, 0)),
	)
}

func (suite *BladeProfileTestSuite) TestPressureSegment() {
	var segment = PSSegment(suite.profile, 0.5, 0.5)
	var precision = 1e-6

	var inletEdgePoint = suite.profile.InletEdge().Point(0.5)
	var inletSegmentPoint = segment.Point(0)
	assert.InDelta(
		suite.T(),
		inletEdgePoint.At(0, 0),
		inletSegmentPoint.At(0, 0),
		precision,
	)
	assert.InDelta(
		suite.T(),
		inletEdgePoint.At(1, 0),
		inletSegmentPoint.At(1, 0),
		precision,
	)

	var outletEdgePoint = suite.profile.OutletEdge().Point(0.5)
	var outletSegmentPoint = segment.Point(1)
	assert.InDelta(
		suite.T(),
		outletEdgePoint.At(0, 0),
		outletSegmentPoint.At(0, 0),
		precision,
	)
	assert.InDelta(
		suite.T(),
		outletEdgePoint.At(1, 0),
		outletSegmentPoint.At(1, 0),
		precision,
	)
}

func (suite *BladeProfileTestSuite) TestSuctionSegment() {
	var segment = SSSegment(suite.profile, 0.5, 0.5)
	var precision = 1e-6

	var inletEdgePoint = suite.profile.InletEdge().Point(0.5)
	var inletSegmentPoint = segment.Point(0)
	assert.InDelta(
		suite.T(),
		inletEdgePoint.At(0, 0),
		inletSegmentPoint.At(0, 0),
		precision,
	)
	assert.InDelta(
		suite.T(),
		inletEdgePoint.At(1, 0),
		inletSegmentPoint.At(1, 0),
		precision,
	)

	var outletEdgePoint = suite.profile.OutletEdge().Point(0.5)
	var outletSegmentPoint = segment.Point(1)
	assert.InDelta(
		suite.T(),
		outletEdgePoint.At(0, 0),
		outletSegmentPoint.At(0, 0),
		precision,
	)
	assert.InDelta(
		suite.T(),
		outletEdgePoint.At(1, 0),
		outletSegmentPoint.At(1, 0),
		precision,
	)
}

func TestBladeProfileTestSuite(t *testing.T) {
	suite.Run(t, new(BladeProfileTestSuite))
}

func testMessage(index interface{}, expected interface{}, got interface{}) string {
	return fmt.Sprintf("Failed %v: expected %v, got %v", index, expected, got)
}
