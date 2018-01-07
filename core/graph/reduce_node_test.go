package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ReduceNodeTestSuite struct {
	suite.Suite
	node ReduceNode
}

func (suite *ReduceNodeTestSuite) SetupTest() {
	var multiplier = func(curr float64, port Port) (float64, error) {
		return curr * port.GetState().(NumberPortState).Num, nil
	}
	var adder = func(curr, new float64) (float64, error) {
		return curr + new, nil
	}

	suite.node = NewReduceNode(multiplier, adder, 1, 0)
}

func (suite *ReduceNodeTestSuite) TestEmpty() {
	var err = suite.node.Process()
	assert.Nil(suite.T(), err)

	var result = suite.node.OutputPort().GetState().(NumberPortState).Num
	assert.InDelta(suite.T(), 0., result, 1e-9)
}

func (suite *ReduceNodeTestSuite) TestSum() {
	var port1 = NewPort()
	var port2 = NewPort()

	port1.SetState(NewNumberPortState(1))
	port2.SetState(NewNumberPortState(2))

	suite.node.AddPortGroup(port1)
	suite.node.AddPortGroup(port2)

	var err = suite.node.Process()
	assert.Nil(suite.T(), err)

	var result = suite.node.OutputPort().GetState().(NumberPortState).Num
	assert.InDelta(suite.T(), 3., result, 1e-9)
}

func (suite *ReduceNodeTestSuite) TestProd() {
	var port1 = NewPort()
	var port2 = NewPort()

	port1.SetState(NewNumberPortState(1))
	port2.SetState(NewNumberPortState(2))

	suite.node.AddPortGroup(port1, port2)

	var err = suite.node.Process()
	assert.Nil(suite.T(), err)

	var result = suite.node.OutputPort().GetState().(NumberPortState).Num
	assert.InDelta(suite.T(), 2., result, 1e-9)
}

func (suite *ReduceNodeTestSuite) TestMultiSum() {
	var port1, port2, port3, port4, port5 = NewPort(), NewPort(), NewPort(), NewPort(), NewPort()

	port1.SetState(NewNumberPortState(1))
	port2.SetState(NewNumberPortState(2))
	port3.SetState(NewNumberPortState(3))
	port4.SetState(NewNumberPortState(4))
	port5.SetState(NewNumberPortState(5))

	suite.node.AddPortGroup(port1, port2)
	suite.node.AddPortGroup(port3, port4, port5)

	var err = suite.node.Process()
	assert.Nil(suite.T(), err)

	var result = suite.node.OutputPort().GetState().(NumberPortState).Num
	assert.InDelta(suite.T(), 62., result, 1e-9)
}

func TestReduceTestSuite(t *testing.T) {
	suite.Run(t, new(ReduceNodeTestSuite))
}
