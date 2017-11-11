package profiles

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

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
	var points = suite.profile.MeanLine().GetPoints(3)

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
	var points = suite.profile.PSLine().GetPoints(3)

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
	var points = suite.profile.SSLine().GetPoints(3)

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
	var points = suite.profile.InletEdge().GetPoints(3)

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
	var points = suite.profile.OutletEdge().GetPoints(3)

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

func TestBladeProfileTestSuite(t *testing.T) {
	suite.Run(t, new(BladeProfileTestSuite))
}
