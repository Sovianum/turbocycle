package variator

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

type SysCallTestSuite struct {
	suite.Suite
	pA float64
	pB float64

	nodeI1    *graph.TestNode
	nodeI2    *graph.TestNode
	nodeA     *graph.TestNode
	nodeB     *graph.TestNode
	nodeC     *graph.TestNode
	nodeD     *graph.TestNode
	assembler graph.VectorAssemblerNode
	out       *graph.TestNode

	network graph.Network
}

func (suite *SysCallTestSuite) SetupTest() {
	suite.pA = 1
	suite.pB = 1

	suite.nodeI1 = graph.NewTestNodeWithoutAction(0, 1, true)
	suite.nodeI1.SetName("i1")

	suite.nodeI2 = graph.NewTestNodeWithoutAction(0, 1, true)
	suite.nodeI2.SetName("i2")

	suite.nodeA = graph.NewTestNodeWithoutAction(1, 2, true)
	suite.nodeA.SetName("a")
	suite.nodeA.SetAction(func() error {
		var requiredPorts, _ = suite.nodeA.GetRequirePorts()
		var updatePorts, _ = suite.nodeA.GetUpdatePorts()

		var portI1 = requiredPorts[0]
		var portA1, portA2 = updatePorts[0], updatePorts[1]

		portA1.SetState(graph.NewNumberPortState(
			portI1.GetState().(graph.NumberPortState).Num * suite.pA,
		))
		portA2.SetState(graph.NewNumberPortState(
			portI1.GetState().(graph.NumberPortState).Num * (1 - suite.pA),
		))
		return nil
	})

	suite.nodeB = graph.NewTestNodeWithoutAction(1, 2, true)
	suite.nodeB.SetName("b")
	suite.nodeB.SetAction(func() error {
		var requiredPorts, _ = suite.nodeB.GetRequirePorts()
		var updatePorts, _ = suite.nodeB.GetUpdatePorts()

		var portI2 = requiredPorts[0]
		var portB1, portB2 = updatePorts[0], updatePorts[1]

		portB1.SetState(graph.NewNumberPortState(
			portI2.GetState().(graph.NumberPortState).Num * suite.pB,
		))
		portB2.SetState(graph.NewNumberPortState(
			portI2.GetState().(graph.NumberPortState).Num * (1 - suite.pB),
		))
		return nil
	})

	suite.nodeC = graph.NewTestNodeWithoutAction(2, 1, true)
	suite.nodeC.SetName("c")
	suite.nodeC.SetAction(func() error {
		var requiredPorts, _ = suite.nodeC.GetRequirePorts()
		var updatePorts, _ = suite.nodeC.GetUpdatePorts()

		var portC1 = updatePorts[0]
		var portA1, portB1 = requiredPorts[0], requiredPorts[1]

		portC1.SetState(graph.NewNumberPortState(
			portA1.GetState().(graph.NumberPortState).Num +
				portB1.GetState().(graph.NumberPortState).Num,
		))
		return nil
	})

	suite.nodeD = graph.NewTestNodeWithoutAction(2, 1, true)
	suite.nodeD.SetName("d")
	suite.nodeD.SetAction(func() error {
		var requiredPorts, _ = suite.nodeD.GetRequirePorts()
		var updatePorts, _ = suite.nodeD.GetUpdatePorts()

		var portD1 = updatePorts[0]
		var portA2, portB2 = requiredPorts[0], requiredPorts[1]

		portD1.SetState(graph.NewNumberPortState(
			portA2.GetState().(graph.NumberPortState).Num -
				portB2.GetState().(graph.NumberPortState).Num,
		))
		return nil
	})

	suite.assembler = graph.NewVectorAssemblerNode()
	suite.assembler.SetName("assembler")

	suite.out = graph.NewTestNodeWithoutAction(1, 0, true)
	suite.out.SetName("out")

	var i1Update, _ = suite.nodeI1.GetUpdatePorts()
	var aRequire, _ = suite.nodeA.GetRequirePorts()
	graph.Link(i1Update[0], aRequire[0])

	var i2Update, _ = suite.nodeI2.GetUpdatePorts()
	var bRequire, _ = suite.nodeB.GetRequirePorts()
	graph.Link(i2Update[0], bRequire[0])

	var aUpdate, _ = suite.nodeA.GetUpdatePorts()
	var cRequire, _ = suite.nodeC.GetRequirePorts()
	graph.Link(aUpdate[0], cRequire[0])

	var bUpdate, _ = suite.nodeB.GetUpdatePorts()
	graph.Link(bUpdate[0], cRequire[1])

	var dRequire, _ = suite.nodeD.GetRequirePorts()
	graph.Link(aUpdate[1], dRequire[0])
	graph.Link(bUpdate[1], dRequire[1])

	var outRequire, _ = suite.out.GetRequirePorts()
	graph.Link(suite.assembler.GetVectorPort(), outRequire[0])

	var cUpdate, _ = suite.nodeC.GetUpdatePorts()
	suite.assembler.AddInputPorts(cUpdate[0])

	var dUpdate, _ = suite.nodeD.GetUpdatePorts()
	suite.assembler.AddInputPorts(dUpdate[0])

	var err error = nil
	suite.network, err = graph.NewNetwork([]graph.Node{
		suite.nodeI1, suite.nodeI2, suite.nodeA, suite.nodeB, suite.nodeC,
		suite.nodeD, suite.assembler, suite.out,
	})
	if err != nil {
		panic(err)
	}
}

func (suite *SysCallTestSuite) TestSysCall_OK() {
	var i1Update, _ = suite.nodeI1.GetUpdatePorts()
	var i2Update, _ = suite.nodeI2.GetUpdatePorts()

	i1Update[0].SetState(graph.NewNumberPortState(4))
	i2Update[0].SetState(graph.NewNumberPortState(2))

	var sysCall = SysCallFromNetwork(
		suite.network, suite.assembler.GetVectorPort(), 0.5, 3, 100, 1e-8,
	)

	var paVariator = FromPointer(&suite.pA)
	var pbVariator = FromPointer(&suite.pB)

	var solverGen = newton.NewUniformNewtonSolverGen(1e-3)
	var varSolver = NewVariatorSolver(sysCall, []Variator{paVariator, pbVariator}, solverGen)

	var x0 = mat.NewVecDense(2, []float64{0, 0})
	var _, err = varSolver.Solve(x0, 1e-8, 1, 100)
	assert.Nil(suite.T(), err)

	assert.InDelta(suite.T(), 0.25, suite.pA, 1e-7)
	assert.InDelta(suite.T(), -0.5, suite.pB, 1e-7)
}

func TestSysCallTestSuite(t *testing.T) {
	suite.Run(t, new(SysCallTestSuite))
}
