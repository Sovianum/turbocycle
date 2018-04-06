package compressor

import (
	"fmt"
	"math"

	common2 "github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/stage/common"
	"github.com/Sovianum/turbocycle/impl/stage/geometry"
	"github.com/Sovianum/turbocycle/material/gases"
)

type dimlessFirstStage func(dRelIn float64) StageNode
type dimlessMidStage func(prevGeom geometry.StageGeometry) StageNode

type StagedCompressorNode interface {
	common.StageChannel
	Stages() []StageNode
	Stage(num int) StageNode
	GetHtLaw() common.DiscreteFunc
	SetHtLaw(htLaw common.DiscreteFunc)
	GetEtaAdLaw() common.DiscreteFunc
	SetEtaAdLaw(etaAdLaw common.DiscreteFunc)
}

func PiStag(node StagedCompressorNode) float64 {
	result := 1.
	for _, stage := range node.Stages() {
		result *= stage.GetDataPack().PiStag
	}
	return result
}

func EtaStag(node StagedCompressorNode) float64 {
	tIn := node.TemperatureInput().GetState().Value().(float64)
	tOut := node.TemperatureOutput().GetState().Value().(float64)
	pi := PiStag(node)
	gas := node.GasInput().GetState().Value().(gases.Gas)
	k := gases.KMean(gas, tIn, tOut, nodes.DefaultN)
	return tIn / (tOut - tIn) * (math.Pow(pi, (k-1)/k) - 1)
}

func NewStagedCompressorNode(
	rpm, dRelIn float64,
	geomList []StageGeometryGenerator,
	htLaw, reactivityLaw, labourCoefLaw, etaAdLaw, caCoefLaw common.DiscreteFunc,
	precision, relaxCoef, initLambda float64, iterLimit int,
) StagedCompressorNode {
	result := &stagedCompressorNode{
		rpm:      rpm,
		dRelIn:   dRelIn,
		stageNum: len(geomList),
		geomList: geomList,

		htLaw:         htLaw,
		reactivityLaw: reactivityLaw,
		labourCoefLaw: labourCoefLaw,
		etaAdLaw:      etaAdLaw,
		caCoefLaw:     caCoefLaw,

		precision:  precision,
		relaxCoef:  relaxCoef,
		initLambda: initLambda,
		iterLimit:  iterLimit,
	}
	result.BaseStage = common.NewBaseStage(result)
	return result
}

type stagedCompressorNode struct {
	*common.BaseStage

	stageNum int
	geomList []StageGeometryGenerator

	rpm    float64
	dRelIn float64

	htLaw         common.DiscreteFunc
	reactivityLaw common.DiscreteFunc
	labourCoefLaw common.DiscreteFunc
	etaAdLaw      common.DiscreteFunc
	caCoefLaw     common.DiscreteFunc

	precision  float64
	relaxCoef  float64
	initLambda float64
	iterLimit  int

	stages []StageNode
}

func (s *stagedCompressorNode) GetEtaAdLaw() common.DiscreteFunc {
	return s.etaAdLaw
}

func (s *stagedCompressorNode) SetEtaAdLaw(etaAdLaw common.DiscreteFunc) {
	s.etaAdLaw = etaAdLaw
}

func (node *stagedCompressorNode) GetHtLaw() common.DiscreteFunc {
	return node.htLaw
}

func (node *stagedCompressorNode) SetHtLaw(htLaw common.DiscreteFunc) {
	node.htLaw = htLaw
}

func (node *stagedCompressorNode) GetName() string {
	return common2.EitherString(node.GetInstanceName(), "StagedCompressorNode")
}

func (node *stagedCompressorNode) Process() error {
	preFirstStage := node.preInitFirstStage()
	preMidStages := node.preInitMidStages()
	stages, err := node.solveAll(preFirstStage, preMidStages)

	if err != nil {
		return err
	}
	node.stages = stages

	lastStage := stages[len(stages)-1]
	graph.CopyAll(
		[]graph.Port{
			lastStage.GasOutput(), lastStage.TemperatureOutput(),
			lastStage.PressureOutput(), lastStage.MassRateOutput(),
			lastStage.VelocityOutput(),
		},
		[]graph.Port{
			node.GasOutput(), node.TemperatureOutput(),
			node.PressureOutput(), node.MassRateOutput(),
			node.VelocityOutput(),
		},
	)
	return nil
}

func (node *stagedCompressorNode) Stages() []StageNode {
	return node.stages
}

func (node *stagedCompressorNode) Stage(num int) StageNode {
	return node.stages[num]
}

func (node *stagedCompressorNode) solveAll(preFirstStage dimlessFirstStage, preMidStages []dimlessMidStage) ([]StageNode, error) {
	stages := make([]StageNode, len(preMidStages)+1)
	firstStage := preFirstStage(node.dRelIn)
	stages[0] = firstStage

	node.initFirstStage(firstStage)

	if err := firstStage.Process(); err != nil {
		return nil, fmt.Errorf("failed on first stage: %s", err.Error())
	}

	for i, dimlessStage := range preMidStages {
		geom := stages[i].GetDataPack().StageGeometry

		stage := dimlessStage(geom)
		prevStage := stages[i]
		common.LinkStages(prevStage, stage)
		common.InitFromPreviousStage(prevStage, stage)

		if err := stage.Process(); err != nil {
			return nil, fmt.Errorf("failed on stage %d: %s", i+1, err.Error())
		}
		stages[i+1] = stage
	}
	return stages, nil
}

func (node *stagedCompressorNode) initFirstStage(firstStage StageNode) {
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

func (node *stagedCompressorNode) preInitMidStages() []dimlessMidStage {
	result := make([]dimlessMidStage, node.stageNum-1)
	for i := range result {
		j := i // another variable used cos variables are captured by reference
		result[j] = func(prevGeom geometry.StageGeometry) StageNode {
			return NewMidStageNode(
				prevGeom,
				node.htLaw(j), node.htLaw(j+1),
				node.reactivityLaw(j+1),
				node.labourCoefLaw(j), node.etaAdLaw(j),
				node.rpm, node.geomList[j],
				node.precision, node.relaxCoef, node.initLambda, node.iterLimit,
			)
		}
	}
	return result
}

func (node *stagedCompressorNode) preInitFirstStage() dimlessFirstStage {
	return func(dRelIn float64) StageNode {
		return NewFirstStageNode(
			dRelIn,
			node.htLaw(0), node.htLaw(1),
			node.reactivityLaw(0), node.reactivityLaw(1),
			node.labourCoefLaw(0), node.etaAdLaw(0), node.caCoefLaw(0),
			node.rpm, node.geomList[0],
			node.precision, node.relaxCoef, node.initLambda, node.iterLimit,
		)
	}
}
