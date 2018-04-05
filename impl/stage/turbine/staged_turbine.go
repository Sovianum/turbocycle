package turbine

import (
	"fmt"

	common2 "github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/stage/common"
)

type StagedTurbineNode interface {
	common.StageChannel
	Stages() []StageNode
	Stage(num int) StageNode
}

func NewStagedTurbineNode(
	rpm, totalHeatDrop float64,
	phiFunc, psiFunc, reactivityFunc,
	airGapRelFunc, heatDropDistributionFunc common.Func1D,
	stageGeomGens []StageGeometryGenerator,
	precision float64,
) StagedTurbineNode {
	result := &stagedTurbineNode{
		rpm:                      rpm,
		totalHeatDrop:            totalHeatDrop,
		phiFunc:                  phiFunc,
		psiFunc:                  psiFunc,
		reactivityFunc:           reactivityFunc,
		airGapRelFunc:            airGapRelFunc,
		heatDropDistributionFunc: heatDropDistributionFunc,
		stageGeomGens:            stageGeomGens,
		precision:                precision,
	}
	result.BaseStage = common.NewBaseStage(result)
	return result
}

type stagedTurbineNode struct {
	*common.BaseStage

	rpm           float64
	totalHeatDrop float64

	phiFunc                  common.Func1D
	psiFunc                  common.Func1D
	reactivityFunc           common.Func1D
	airGapRelFunc            common.Func1D
	heatDropDistributionFunc common.Func1D

	stageGeomGens []StageGeometryGenerator
	precision     float64

	stages []StageNode
}

func (node *stagedTurbineNode) GetName() string {
	return common2.EitherString(node.GetInstanceName(), "StagedTurbine")
}

func (node *stagedTurbineNode) Process() error {
	node.createStages()
	node.linkStages()
	node.initFirstStage(node.stages[0])
	for i, stage := range node.stages {
		if err := stage.Process(); err != nil {
			return fmt.Errorf("failed on stage %d: %s", i, err.Error())
		}
	}
	return nil
}

func (node *stagedTurbineNode) Stages() []StageNode {
	return node.stages
}

func (node *stagedTurbineNode) Stage(num int) StageNode {
	return node.stages[num]
}

func (node *stagedTurbineNode) initFirstStage(firstStage StageNode) {
	graph.CopyAll(
		[]graph.Port{
			node.GasInput(), node.TemperatureInput(),
			node.PressureInput(), node.MassRateInput(),
			node.VelocityInput(),
		},
		[]graph.Port{
			firstStage.GasInput(), firstStage.TemperatureInput(),
			firstStage.PressureInput(), firstStage.MassRateInput(),
			firstStage.VelocityInput(),
		},
	)
}

func (node *stagedTurbineNode) linkStages() {
	for i := 0; i != len(node.stages)-1; i++ {
		common.LinkStages(node.stages[i], node.stages[i+1])
	}
}

func (node *stagedTurbineNode) createStages() {
	stages := make([]StageNode, node.stageNum())
	stagePositions := node.stagePositions()
	normDistrib := node.heatDropDistributionFunc.GetUnitNormalizedSamples(stagePositions)

	for i, geomGen := range node.stageGeomGens {
		x := stagePositions[i]
		stages[i] = NewTurbineStageNode(
			node.rpm, node.totalHeatDrop*normDistrib[i],
			node.reactivityFunc(stagePositions[i]),
			node.phiFunc(x), node.psiFunc(x),
			node.airGapRelFunc(x), node.precision, geomGen,
		)
	}
	node.stages = stages
}

func (node *stagedTurbineNode) stagePositions() []float64 {
	result := make([]float64, node.stageNum())
	for i := range result {
		result[i] = float64(i)
	}
	return result
}

func (node *stagedTurbineNode) stageNum() int {
	return len(node.stageGeomGens)
}
