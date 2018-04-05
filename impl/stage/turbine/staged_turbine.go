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
	Ht() float64
	SetHt(ht float64)
}

func PiStag(node StagedTurbineNode) float64 {
	result := 1.
	for _, stage := range node.Stages() {
		result *= stage.GetDataPack().Pi
	}
	return result
}

func NewStagedTurbineNode(
	rpm, alpha1, totalHeatDrop, lRelIn float64,
	phiFunc, psiFunc, reactivityFunc,
	airGapRelFunc, heatDropDistributionFunc common.Func1D,
	incompleteStageGeomGens []IncompleteStageGeometryGenerator,
	precision float64,
) StagedTurbineNode {
	result := &stagedTurbineNode{
		rpm:           rpm,
		alpha1:        alpha1,
		totalHeatDrop: totalHeatDrop,
		lRelIn:        lRelIn,

		phiFunc:                  phiFunc,
		psiFunc:                  psiFunc,
		reactivityFunc:           reactivityFunc,
		airGapRelFunc:            airGapRelFunc,
		heatDropDistributionFunc: heatDropDistributionFunc,
		incompleteStageGeomGens:  incompleteStageGeomGens,
		precision:                precision,
	}
	result.BaseStage = common.NewBaseStage(result)
	return result
}

type stagedTurbineNode struct {
	*common.BaseStage

	alpha1        float64
	rpm           float64
	totalHeatDrop float64
	lRelIn        float64

	phiFunc                  common.Func1D
	psiFunc                  common.Func1D
	reactivityFunc           common.Func1D
	airGapRelFunc            common.Func1D
	heatDropDistributionFunc common.Func1D

	incompleteStageGeomGens []IncompleteStageGeometryGenerator

	precision float64

	stages []StageNode
}

func (node *stagedTurbineNode) SetHt(ht float64) {
	node.totalHeatDrop = ht
}

func (node *stagedTurbineNode) Ht() float64 {
	return node.totalHeatDrop
}

func (node *stagedTurbineNode) GetName() string {
	return common2.EitherString(node.GetInstanceName(), "StagedTurbine")
}

func (node *stagedTurbineNode) Process() error {
	stages := make([]StageNode, node.stageNum())

	normDistrib := node.heatDropDistributionFunc.GetUnitNormalizedSamples(node.stagePositions())
	firstStage := node.createFirstStage(
		node.incompleteStageGeomGens[0].GetGenerator(node.lRelIn),
		node.alpha1, normDistrib[0]*node.totalHeatDrop,
	)
	node.initFirstStage(firstStage)
	if err := firstStage.Process(); err != nil {
		return fmt.Errorf("failed on first stage: %s", err.Error())
	}
	stages[0] = firstStage

	for i := range normDistrib[1:] {
		rotorGeom := stages[i].GetDataPack().StageGeometry.RotorGeometry()
		stage := node.createMidStage(
			i+1,
			node.incompleteStageGeomGens[i+1].GetGenerator(
				LRelOutGap(stages[i].StageGeomGen().RotorGenerator()),
			),
			rotorGeom.MeanProfile().Diameter(rotorGeom.XGapOut()),
			normDistrib[i+1]*node.totalHeatDrop,
		)
		common.LinkStages(stages[i], stage)
		common.InitFromPreviousStage(stages[i], stage)

		if err := stage.Process(); err != nil {
			return fmt.Errorf("failed on stage %d: %s", i+1, err.Error())
		}
		stages[i+1] = stage
	}
	node.stages = stages
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

func (node *stagedTurbineNode) createFirstStage(
	stageGeomGen StageGeometryGenerator, alpha1, heatDrop float64,
) StageNode {
	return NewTurbineFirstStageNode(
		alpha1, node.rpm, heatDrop,
		node.reactivityFunc(0),
		node.phiFunc(0),
		node.psiFunc(0),
		node.airGapRelFunc(0),
		node.precision, stageGeomGen,
	)
}

func (node *stagedTurbineNode) createMidStage(
	ind int,
	stageGeomGen StageGeometryGenerator,
	dMeanIn, heatDrop float64,
) StageNode {
	floatInd := float64(ind)
	return NewTurbineMidStageNode(
		dMeanIn, node.rpm, heatDrop,
		node.reactivityFunc(floatInd),
		node.phiFunc(floatInd),
		node.psiFunc(floatInd),
		node.airGapRelFunc(floatInd),
		node.precision, stageGeomGen,
	)
}

func (node *stagedTurbineNode) stagePositions() []float64 {
	result := make([]float64, node.stageNum())
	for i := range result {
		result[i] = float64(i)
	}
	return result
}

func (node *stagedTurbineNode) stageNum() int {
	return len(node.incompleteStageGeomGens)
}
