package compressor

import (
	"fmt"
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/common/gdf"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	states2 "github.com/Sovianum/turbocycle/impl/engine/states"
	common2 "github.com/Sovianum/turbocycle/impl/stage/common"
	"github.com/Sovianum/turbocycle/impl/stage/geometry"
	"github.com/Sovianum/turbocycle/impl/stage/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type StageNode interface {
	common2.StageChannel
	GetDataPack() *DataPack
}

type DataPack struct {
	Err error

	HTCoef          float64 `json:"heat_drop_coef"`
	HT              float64 `json:"heat_drop"`
	LabourCoef      float64 `json:"labour_coef"`
	Labour          float64 `json:"labour"`
	AdiabaticLabour float64 `json:"adiabatic_labour"`
	EtaAd           float64 `json:"eta_ad"`
	T1Stag          float64 `json:"t_in"`
	TemperatureDrop float64 `json:"temperature_drop"`
	T3Stag          float64 `json:"t_out"`
	P1Stag          float64 `json:"p_in"`
	P3Stag          float64 `json:"p_out"`
	PiStag          float64 `json:"pi_stag"`

	StageGeometry geometry.StageGeometry `json:"stage_geometry"`
	Area1         float64                `json:"area_1"`
	Area3         float64                `json:"area_3"`

	UOut           float64                 `json:"u_out"`
	InletTriangle  states.VelocityTriangle `json:"inlet_triangle"`
	OutletTriangle states.VelocityTriangle `json:"outlet_triangle"`
	MidTriangle    states.VelocityTriangle `json:"mid_triangle"`
}

func NewMidStageNode(
	prevStageGeom geometry.StageGeometry,
	htCoef, htCoefNext,
	reactivityNext,
	labourCoef, etaAd,
	rpm float64,
	stageGeomGen StageGeometryGenerator,
	precision, relaxCoef, initLambda float64, iterLimit int,
) StageNode {
	prevGeom := prevStageGeom.StatorGeometry()
	dRelIn := geometry.DRel(prevGeom.XGapOut(), prevGeom)
	// zeros below are values not used in mid stage calculation
	result := NewFirstStageNode(
		dRelIn,
		htCoef, htCoefNext,
		0, reactivityNext,
		labourCoef, etaAd, 0,
		rpm, stageGeomGen,
		precision, relaxCoef, initLambda, iterLimit,
	).(*stageNode)
	result.prevStageGeom = prevStageGeom
	result.isFirstStage = false
	return result
}

func NewFirstStageNode(
	dRelIn,
	htCoef, htCoefNext,
	reactivity, reactivityNext,
	labourCoef, etaAd, caCoef,
	rpm float64,
	stageGeomGen StageGeometryGenerator,
	precision, relaxCoef, initLambda float64, iterLimit int,
) StageNode {
	result := &stageNode{
		dRelIn:         dRelIn,
		htCoef:         htCoef,
		htCoefNext:     htCoefNext,
		reactivity:     reactivity,
		reactivityNext: reactivityNext,
		labourCoef:     labourCoef,
		etaAd:          etaAd,
		caCoef:         caCoef,
		rpm:            rpm,
		stageGeomGen:   stageGeomGen,
		precision:      precision,
		relaxCoef:      relaxCoef,
		initLambda:     initLambda,
		iterLimit:      iterLimit,
		isFirstStage:   true,
	}
	result.BaseStage = common2.NewBaseStage(result)
	return result
}

type stageNode struct {
	*common2.BaseStage

	dRelIn     float64
	caCoef     float64
	labourCoef float64
	etaAd      float64

	htCoef         float64
	htCoefNext     float64
	reactivity     float64
	reactivityNext float64

	rpm float64

	stageGeomGen StageGeometryGenerator

	precision  float64
	relaxCoef  float64
	initLambda float64
	iterLimit  int

	pack *DataPack

	isFirstStage  bool
	prevStageGeom geometry.StageGeometry
}

func (node *stageNode) GetDataPack() *DataPack {
	return node.pack
}

func (node *stageNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "CompressorStage")
}

func (node *stageNode) Process() error {
	node.pack = new(DataPack)
	if node.isFirstStage {
		node.inletVelocitiesFirstStage(node.pack)
	} else {
		node.inletVelocities(node.pack)
	}
	node.hT(node.pack)
	node.labour(node.pack)
	node.temperatures(node.pack)
	node.pressures(node.pack)
	node.outletVelocities(node.pack)
	node.midVelocities(node.pack)

	node.GasOutput().SetState(states2.NewGasPortState(node.gas()))
	node.PressureOutput().SetState(states2.NewPressurePortState(node.pack.P3Stag))
	node.TemperatureOutput().SetState(states2.NewTemperaturePortState(node.pack.T3Stag))
	graph.CopyState(node.MassRateInput(), node.MassRateOutput())
	node.VelocityOutput().SetState(states.NewVelocityPortState(node.pack.OutletTriangle, states.CompressorTriangleType))
	return node.pack.Err
}

func (node *stageNode) midVelocities(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: mid_velocities", pack.Err.Error())
		return
	}
	stageGeom := pack.StageGeometry
	rotorGeom := stageGeom.RotorGeometry()

	dRelOut := geometry.DRel(rotorGeom.XGapOut(), rotorGeom)
	rRelOut := geometry.RRel(dRelOut)

	cu1 := pack.InletTriangle.CU()
	uOut1 := pack.UOut

	cu2 := 1 / rRelOut * (node.htCoef + cu1/uOut1*rRelOut)

	ca1 := pack.InletTriangle.CA()
	ca3 := pack.InletTriangle.CA()
	ca2 := (ca1 + ca3) / 2

	dOutIn := rotorGeom.OuterProfile().Diameter(0)
	dOutOut := rotorGeom.OuterProfile().Diameter(rotorGeom.XGapOut())
	u2 := uOut1 / dOutIn * dOutOut * rRelOut

	pack.MidTriangle = states.NewCompressorVelocityTriangleFromProjections(cu2, ca2, u2)
}

func (node *stageNode) outletVelocities(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: outlet_velocities", pack.Err.Error())
		return
	}
	bGeom := pack.StageGeometry.StatorGeometry()
	dRelOut := geometry.DRel(bGeom.XGapOut(), bGeom)
	rRel := geometry.RRel(dRelOut)

	gas := node.gas()
	k := gases.K(gas, pack.T3Stag)
	f3 := geometry.Area(bGeom.XGapOut(), bGeom)
	cuCoef := rRel*(1-node.reactivityNext) - node.htCoefNext/(2*rRel)
	massRate := node.massRate()
	massRateFactor := massRate * math.Sqrt(pack.T3Stag) / pack.P3Stag

	lambda3Func := func(lambda3 float64) (float64, error) {
		sinAlpha := massRateFactor / gdf.Q(lambda3, k, gas.R())
		cosAlpha := math.Sqrt(1 - sinAlpha*sinAlpha)
		return cuCoef * pack.UOut / (gdf.ACrit(k, gas.R(), pack.T3Stag) * cosAlpha * f3), nil
	}

	lambda3, err := common.SolveIterativelyWithValidation(
		lambda3Func, common2.NotNanValidator, node.initLambda,
		node.precision, node.relaxCoef, node.iterLimit,
	)
	if err != nil {
		pack.Err = fmt.Errorf("%s: outlet_velocities", err.Error())
		return
	}
	q3 := gdf.Q(lambda3, k, gas.R())
	alpha3 := math.Asin(massRateFactor / q3)
	cu := cuCoef * pack.UOut
	ca := math.Tan(alpha3) * cu

	dOutIn1 := pack.StageGeometry.RotorGeometry().OuterProfile().Diameter(0)
	dOutIn3 := pack.StageGeometry.StatorGeometry().OuterProfile().Diameter(0)
	u3 := pack.UOut / dOutIn1 * dOutIn3 * rRel

	pack.OutletTriangle = states.NewCompressorVelocityTriangleFromProjections(cu, ca, u3)
	pack.Area3 = geometry.Area(
		pack.StageGeometry.StatorGeometry().XGapOut(),
		pack.StageGeometry.StatorGeometry(),
	)
}

func (node *stageNode) pressures(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: pressure", pack.Err.Error())
		return
	}
	cpMean := gases.CpMean(node.gas(), pack.T1Stag, pack.T3Stag, nodes.DefaultN)
	kMean := gases.KMean(node.gas(), pack.T1Stag, pack.T3Stag, nodes.DefaultN)
	pi := math.Pow(
		1+pack.AdiabaticLabour/(pack.T1Stag*cpMean),
		kMean/(kMean-1),
	)
	pack.P1Stag = node.p1Stag()
	pack.PiStag = pi
	pack.P3Stag = pack.P1Stag * pi
}

func (node *stageNode) temperatures(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: temperature", pack.Err.Error())
		return
	}

	iterFunc := func(tOutCurr float64) (float64, error) {
		cp := gases.CpMean(node.gas(), node.t1Stag(), tOutCurr, nodes.DefaultN)
		return node.t1Stag() + pack.Labour/cp, nil
	}

	tOut, err := common.SolveIterativelyWithValidation(
		iterFunc,
		common2.NotNanValidator,
		node.t1Stag(), node.precision, node.relaxCoef, node.iterLimit,
	)
	if err != nil {
		pack.Err = fmt.Errorf("%s: tOut", err.Error())
	}

	pack.T1Stag = node.t1Stag()
	pack.TemperatureDrop = tOut - node.t1Stag()
	pack.T3Stag = tOut
}

func (node *stageNode) labour(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: labour", pack.Err.Error())
		return
	}
	pack.Labour = node.labourCoef * pack.HT
	pack.EtaAd = node.etaAd
	pack.AdiabaticLabour = pack.Labour * node.etaAd
}

func (node *stageNode) hT(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: hT", pack.Err.Error())
		return
	}
	pack.HTCoef = node.htCoef
	u := pack.InletTriangle.U()
	pack.HT = node.htCoef * u * u
}

func (node *stageNode) inletVelocitiesFirstStage(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: inlet_velocities_first_stage", pack.Err.Error())
		return
	}
	gas := node.gas()
	rRel := geometry.RRel(node.dRelIn)
	cuCoef := rRel*(1-node.reactivity) - node.htCoef/(2*rRel)
	alpha1 := math.Atan2(node.caCoef, cuCoef)
	kGas := gases.K(gas, pack.T1Stag)

	massRate := node.massRate()
	t1 := node.t1Stag()
	p1 := node.p1Stag()
	massRateFactor := massRate * math.Sqrt(t1) / p1
	area1Func := func(lambda1 float64) float64 {
		qLambda := gdf.Q(lambda1, kGas, gas.R())
		lambdaFactor := 1 / (qLambda * math.Sin(alpha1))
		return lambdaFactor * massRateFactor
	}

	caFactor := node.caCoef / (math.Sin(alpha1) * gdf.ACrit(kGas, gas.R(), node.t1Stag()))
	lambda1Func := func(lambda1 float64) (float64, error) {
		rpmFactor := node.rpm / 30
		f1 := area1Func(lambda1)
		areaFactor := math.Sqrt(
			math.Pi * f1 / (1 - node.dRelIn*node.dRelIn),
		)
		return rpmFactor * areaFactor * caFactor, nil
	}

	lambda1, err := common.SolveIterativelyWithValidation(
		lambda1Func,
		common2.NotNanValidator,
		node.initLambda, node.precision, node.relaxCoef, node.iterLimit,
	)
	if err != nil {
		pack.Err = fmt.Errorf("%s: inlet_velocities", err)
		return
	}

	area1 := area1Func(lambda1)
	pack.Area1 = area1

	dOutIn := math.Sqrt(4 / math.Pi * 1 / (1 - node.dRelIn*node.dRelIn) * area1)
	pack.StageGeometry = node.stageGeomGen.Generate(dOutIn)

	u1Out := math.Pi * dOutIn * node.rpm / 60
	pack.UOut = u1Out
	ca1 := node.caCoef * u1Out
	cu1 := cuCoef * u1Out

	u1 := rRel * u1Out

	pack.InletTriangle = states.NewCompressorVelocityTriangleFromProjections(cu1, ca1, u1)
}

func (node *stageNode) inletVelocities(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: inlet_velocities", pack.Err.Error())
		return
	}
	prevGeom := node.prevStageGeom.StatorGeometry()
	area1 := geometry.Area(prevGeom.XGapOut(), prevGeom)

	pack.Area1 = area1

	dOutIn := math.Sqrt(4 / math.Pi * 1 / (1 - node.dRelIn*node.dRelIn) * area1)
	pack.StageGeometry = node.stageGeomGen.Generate(dOutIn)

	u1Out := math.Pi * dOutIn * node.rpm / 60
	pack.UOut = u1Out

	pack.InletTriangle = node.VelocityInput().GetState().Value().(states.VelocityTriangle)
}

// below are private accessors

func (node *stageNode) inletTriangle() states.VelocityTriangle {
	return node.VelocityInput().GetState().Value().(states.VelocityTriangle)
}

func (node *stageNode) massRate() float64 {
	return node.MassRateInput().GetState().(states2.MassRatePortState).MassRate
}

func (node *stageNode) p1Stag() float64 {
	return node.PressureInput().GetState().(states2.PressurePortState).PStag
}

func (node *stageNode) t1Stag() float64 {
	return node.TemperatureInput().GetState().(states2.TemperaturePortState).TStag
}

func (node *stageNode) gas() gases.Gas {
	return node.GasInput().GetState().(states2.GasPortState).Gas
}

func (node *stageNode) statorInletTriangle() states.VelocityTriangle {
	return node.VelocityInput().GetState().(states.VelocityPortState).Triangle
}
