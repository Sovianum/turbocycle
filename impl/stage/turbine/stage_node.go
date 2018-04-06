package turbine

import (
	"fmt"
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	common2 "github.com/Sovianum/turbocycle/impl/stage/common"
	"github.com/Sovianum/turbocycle/impl/stage/geometry"
	states2 "github.com/Sovianum/turbocycle/impl/stage/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type StageMode int8

const (
	only StageMode = 1 + iota // todo remove only mode
	first
	mid
)

type StageNode interface {
	common2.StageChannel
	SetStageMode(mode StageMode)
	SetAlpha1FirstStage(alpha1FirstStage float64)
	StageGeomGen() StageGeometryGenerator
	Ht() float64
	Reactivity() float64
	GetDataPack() DataPack
}

func InitFromTurbineNode(stage StageNode, turbine constructive.TurbineNode, massRate, alpha1 float64) {
	stage.GasInput().SetState(states.NewGasPortState(turbine.InputGas()))
	stage.TemperatureInput().SetState(states.NewTemperaturePortState(turbine.TStagIn()))
	stage.PressureInput().SetState(states.NewPressurePortState(turbine.PStagIn()))
	stage.MassRateInput().SetState(states.NewMassRatePortState(massRate))
	stage.SetAlpha1FirstStage(alpha1)
}

func NewTurbineMidStageNode(
	dMeanIn, n, stageHeatDrop,
	reactivity, phi, psi, airGapRel, precision float64,
	gen StageGeometryGenerator,
) StageNode {
	result := NewTurbineSingleStageNode(n, stageHeatDrop, reactivity, phi, psi, airGapRel, precision, gen).(*turbineStageNode)
	result.mode = mid
	result.dMeanIn = dMeanIn
	return result
}

func NewTurbineFirstStageNode(
	alpha1, n, stageHeatDrop, reactivity, phi, psi, airGapRel, precision float64,
	gen StageGeometryGenerator,
) StageNode {
	result := &turbineStageNode{
		n:                n,
		stageHeatDrop:    stageHeatDrop,
		reactivity:       reactivity,
		phi:              phi,
		psi:              psi,
		airGapRel:        airGapRel,
		precision:        precision,
		stageGeomGen:     gen,
		alpha1FirstStage: alpha1,
		mode:             first,
	}
	result.BaseStage = common2.NewBaseStage(result)
	result.VelocityInput().SetState(
		states2.NewVelocityPortState(
			states2.NewInletTriangle(0, 0, math.Pi/2), states2.InletTriangleType,
		),
	)
	return result
}

func NewTurbineSingleStageNode(
	n, stageHeatDrop, reactivity, phi, psi, airGapRel, precision float64,
	gen StageGeometryGenerator,
) StageNode {
	result := &turbineStageNode{
		n:                n,
		stageHeatDrop:    stageHeatDrop,
		reactivity:       reactivity,
		phi:              phi,
		psi:              psi,
		airGapRel:        airGapRel,
		precision:        precision,
		stageGeomGen:     gen,
		alpha1FirstStage: math.NaN(),
		mode:             only,
	}
	result.BaseStage = common2.NewBaseStage(result)
	result.VelocityInput().SetState(
		states2.NewVelocityPortState(
			states2.NewInletTriangle(0, 0, math.Pi/2), states2.InletTriangleType,
		),
	)
	return result
}

type turbineStageNode struct {
	*common2.BaseStage

	lRelIn float64

	dMeanIn float64

	n                float64
	stageHeatDrop    float64
	reactivity       float64
	phi              float64
	psi              float64
	airGapRel        float64
	alpha1FirstStage float64

	stageGeomGen StageGeometryGenerator

	precision float64

	mode StageMode
	pack *DataPack
}

type DataPack struct {
	Err     error
	chapter string

	RPM        float64
	Reactivity float64
	Phi        float64
	Psi        float64
	AirGapRel  float64

	EtaTStag                   float64                  `json:"eta_t_stag"`
	StageHeatDropStag          float64                  `json:"stage_heat_drop_stag"`
	StageLabour                float64                  `json:"stage_labour"`
	P2Stag                     float64                  `json:"p_2_stag"`
	T2Stag                     float64                  `json:"t_2_stag"`
	EtaT                       float64                  `json:"eta_t"`
	VentilationSpecificLoss    float64                  `json:"ventilation_specific_loss"`
	AirGapSpecificLoss         float64                  `json:"air_gap_specific_loss"`
	OutletVelocitySpecificLoss float64                  `json:"outlet_velocity_specific_loss"`
	RotorSpecificLoss          float64                  `json:"rotor_specific_loss"`
	StatorSpecificLoss         float64                  `json:"stator_specific_loss"`
	EtaU                       float64                  `json:"eta_u"`
	MeanRadiusLabour           float64                  `json:"mean_radius_labour"`
	Pi                         float64                  `json:"pi"`
	PiStag                     float64                  `json:"pi_stag"`
	RotorOutletTriangle        states2.VelocityTriangle `json:"rotor_outlet_triangle"`
	C2u                        float64                  `json:"c_2_u"`
	Beta2                      float64                  `json:"beta_2"`
	C2a                        float64                  `json:"c_2_a"`
	Density2                   float64                  `json:"density_2"`
	P2                         float64                  `json:"p_2"`
	T2Prime                    float64                  `json:"t_2_prime"`
	T2                         float64                  `json:"t_2"`
	W2                         float64                  `json:"w_2"`
	WAd2                       float64                  `json:"w_ad_2"`
	RotorHeatDrop              float64                  `json:"rotor_heat_drop"`
	Pw1                        float64                  `json:"pw_1"`
	Tw1                        float64                  `json:"tw_1"`
	U2                         float64                  `json:"u_2"`
	RotorInletTriangle         states2.VelocityTriangle `json:"rotor_inlet_triangle"`
	Alpha1                     float64                  `json:"alpha_1"`
	U1                         float64                  `json:"u_1"`
	C1a                        float64                  `json:"c_1_a"`
	RotorMeanInletDiameter     float64                  `json:"d_rotor_blade_in_mean"`
	Area1                      float64                  `json:"area_1"`
	Density1                   float64                  `json:"density_1"`
	P1                         float64                  `json:"p_1"`
	T1                         float64                  `json:"t_1"`
	C1                         float64                  `json:"c_1"`
	C1Ad                       float64                  `json:"c_1_ad"`
	T1Prime                    float64                  `json:"t_1_prime"`
	StatorHeatDrop             float64                  `json:"stator_heat_drop"`
	StageGeometry              geometry.StageGeometry   `json:"stage_geometry"`
	StatorMeanInletDiameter    float64                  `json:"stator_mean_inlet_diameter"`
	Density0                   float64                  `json:"density_0"`
	P0                         float64                  `json:"p_0"`
	T0                         float64                  `json:"t_0"`
}

func (pack *DataPack) setChapterName(name string) (isErr bool) {
	isErr = pack.Err != nil
	pack.chapter = name
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: %s", pack.Err.Error(), name)
	}
	return
}

func (pack *DataPack) checkAllAndSet(value float64, dest *float64) {
	pack.validateAndSet(value, dest, common2.ComplexPositiveValidator)
}

func (pack *DataPack) checkFiniteAndSet(value float64, dest *float64) {
	pack.validateAndSet(value, dest, common2.FiniteValidator)
}

func (pack *DataPack) validateAndSet(value float64, dest *float64, validator common2.Validator) {
	if pack.Err != nil {
		return
	}
	if err := validator(value); err != nil {
		pack.Err = fmt.Errorf("%s: %s", err.Error(), pack.chapter)
		return
	}
	*dest = value
}

func (node *turbineStageNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "TurbineStage")
}

func (node *turbineStageNode) Process() error {
	node.pack = node.getDataPack()
	if node.pack.Err != nil {
		return node.pack.Err
	}
	node.pack.PiStag = node.p0Stag() / node.pack.P2Stag

	graph.CopyState(node.GasInput(), node.GasOutput())
	node.TemperatureOutput().SetState(states.NewTemperaturePortState(node.pack.T2Stag))
	node.PressureOutput().SetState(states.NewPressurePortState(node.pack.P2Stag))
	node.MassRateOutput().SetState(states.NewMassRatePortState(node.massRate())) // mass rate is constant
	node.VelocityOutput().SetState(states2.NewVelocityPortState(node.pack.RotorOutletTriangle, states2.OutletTriangleType))
	return nil
}

func (node *turbineStageNode) GetDataPack() DataPack {
	if node.pack == nil {
		node.pack = node.getDataPack()
	}
	return *node.pack
}

func (node *turbineStageNode) SetStageMode(mode StageMode) {
	node.mode = mode
}

func (node *turbineStageNode) SetAlpha1FirstStage(alpha1FirstStage float64) {
	node.alpha1FirstStage = alpha1FirstStage
}

func (node *turbineStageNode) StageGeomGen() StageGeometryGenerator {
	return node.stageGeomGen
}

func (node *turbineStageNode) Ht() float64 {
	return node.stageHeatDrop
}

func (node *turbineStageNode) Reactivity() float64 {
	return node.reactivity
}

func (node *turbineStageNode) getDataPack() *DataPack {
	var pack = new(DataPack)
	pack.T0 = node.t0Stag()
	pack.P0 = node.p0Stag()
	switch node.mode {
	case first:
		node.initCalcFirst(pack)
	case only:
		node.initCalcOnly(pack)
	case mid:
		node.initCalcMid(pack)
	default:
		pack.Err = fmt.Errorf("invalid mode")
	}

	node.relativeThermo(pack)
	node.thermo2(pack)
	node.velocity2(pack)
	node.pi(pack)
	node.meanRadiusLabour(pack)
	node.etaU(pack)
	node.losses(pack)
	node.etaT(pack)
	node.t2Stag(pack)
	node.p2Stag(pack)
	node.stageLabour(pack)
	node.stageHeatDropStag(pack)
	node.etaTStag(pack)

	node.pushExtraData(pack)
	return pack
}

func (node *turbineStageNode) pushExtraData(pack *DataPack) {
	pack.RPM = node.n
	pack.Reactivity = node.reactivity
	pack.Phi = node.phi
	pack.Psi = node.psi
	pack.AirGapRel = node.airGapRel
}

func (node *turbineStageNode) initCalcOnly(pack *DataPack) {
	node.thermo0(pack)
	node.getStageGeometry(pack)
	node.thermo1(pack)
	node.velocity1(pack)
}

func (node *turbineStageNode) initCalcFirst(pack *DataPack) {
	pack.Alpha1 = node.alpha1FirstStage
	node.thermo1(pack)
	node.velocity1FirstStage(pack)
}

func (node *turbineStageNode) initCalcMid(pack *DataPack) {
	pack.StageGeometry = node.stageGeomGen.GenerateFromStatorInlet(node.dMeanIn)
	node.thermo0(pack)
	node.thermo1(pack)
	node.velocity1(pack)
}

func (node *turbineStageNode) etaTStag(pack *DataPack) {
	if pack.setChapterName("etaTStag") {
		return
	}
	pack.EtaTStag = pack.StageLabour / pack.StageHeatDropStag
}

func (node *turbineStageNode) stageHeatDropStag(pack *DataPack) {
	if pack.setChapterName("stageHeatDropStag") {
		return
	}
	var cp = gases.CpMean(node.gas(), node.t0Stag(), pack.T2, nodes.DefaultN)
	var k = gases.KMean(node.gas(), node.t0Stag(), pack.T2, nodes.DefaultN)
	pack.StageHeatDropStag = cp * node.t0Stag() * (1 - math.Pow(pack.P2Stag/node.p0Stag(), (k-1)/k))
}

func (node *turbineStageNode) stageLabour(pack *DataPack) {
	if pack.setChapterName("stageLabour") {
		return
	}
	pack.StageLabour = node.stageHeatDrop * pack.EtaT
}

func (node *turbineStageNode) p2Stag(pack *DataPack) {
	if pack.setChapterName("p2Stag") {
		return
	}
	var k = gases.K(node.gas(), pack.T2) // todo check if correct temperature
	pack.P2Stag = pack.P2 * math.Pow(pack.T2Stag/pack.T2, k/(k-1))
}

func (node *turbineStageNode) t2Stag(pack *DataPack) {
	if pack.setChapterName("t2Stag") {
		return
	}
	var cp = node.gas().Cp(pack.T2) // todo check if correct temperature
	pack.T2Stag = pack.T2 +
		(pack.AirGapSpecificLoss+pack.VentilationSpecificLoss+pack.OutletVelocitySpecificLoss)/cp
}

func (node *turbineStageNode) etaT(pack *DataPack) {
	if pack.setChapterName("etaT") {
		return
	}
	pack.EtaT = pack.EtaU - (pack.AirGapSpecificLoss+pack.VentilationSpecificLoss)/node.stageHeatDrop
}

func (node *turbineStageNode) losses(pack *DataPack) {
	if pack.setChapterName("losses") {
		return
	}
	//node.ventilationSpecificLoss(pack)

	defaultSpecificLoss := (math.Pow(node.phi, -2) - 1) * math.Pow(pack.C1, 2) / 2
	pack.checkAllAndSet(
		defaultSpecificLoss*pack.T2Prime/pack.T1,
		&pack.StatorSpecificLoss,
	)
	pack.checkAllAndSet(
		(math.Pow(node.psi, -2)-1)*math.Pow(pack.W2, 2)/2,
		&pack.RotorSpecificLoss,
	)
	pack.checkAllAndSet(
		math.Pow(pack.RotorOutletTriangle.C(), 2)/2,
		&pack.OutletVelocitySpecificLoss,
	)

	lRel := geometry.RelativeHeight(
		pack.StageGeometry.RotorGeometry().XBladeOut(),
		pack.StageGeometry.RotorGeometry(),
	)
	pack.checkAllAndSet(
		1.37*(1+1.6*node.reactivity)*(1+lRel)*node.airGapRel*pack.MeanRadiusLabour,
		&pack.AirGapSpecificLoss,
	)

	x := pack.StageGeometry.RotorGeometry().XBladeOut()
	d := pack.StageGeometry.RotorGeometry().MeanProfile().Diameter(x)
	ventilationPower := 1.07 * math.Pow(d, 2) * math.Pow(pack.U2/100, 3) * pack.Density2 * 1000
	pack.checkAllAndSet(
		ventilationPower/node.massRate(),
		&pack.VentilationSpecificLoss,
	)
}

func (node *turbineStageNode) etaU(pack *DataPack) {
	if pack.setChapterName("etaU") {
		return
	}
	pack.checkAllAndSet(
		pack.MeanRadiusLabour/node.stageHeatDrop,
		&pack.EtaU,
	)
}

func (node *turbineStageNode) meanRadiusLabour(pack *DataPack) {
	if pack.setChapterName("meanRadiusLabour") {
		return
	}
	pack.checkAllAndSet(
		pack.RotorInletTriangle.CU()*pack.RotorInletTriangle.U()+pack.RotorOutletTriangle.CU()*pack.RotorOutletTriangle.U(),
		&pack.MeanRadiusLabour,
	)
}

func (node *turbineStageNode) pi(pack *DataPack) {
	if pack.setChapterName("pi") {
		return
	}
	pack.checkAllAndSet(
		node.p0Stag()/pack.P2,
		&pack.Pi,
	)
}

func (node *turbineStageNode) velocity2(pack *DataPack) {
	if pack.setChapterName("velocity2") {
		return
	}
	area := geometry.Area(pack.StageGeometry.RotorGeometry().XBladeOut(), pack.StageGeometry.RotorGeometry())

	pack.checkAllAndSet(
		node.massRate()/(pack.Density2*area),
		&pack.C2a,
	)
	pack.checkAllAndSet(
		math.Asin(pack.C2a/pack.W2),
		&pack.Beta2,
	)
	pack.checkFiniteAndSet(
		pack.W2*math.Cos(pack.Beta2)-pack.U2,
		&pack.C2u,
	)
	pack.RotorOutletTriangle = states2.NewOutletTriangleFromProjections(
		pack.C2u, pack.C2a, pack.U2,
	)
}

func (node *turbineStageNode) thermo2(pack *DataPack) {
	if pack.setChapterName("thermo2") {
		return
	}
	t1 := pack.T1
	w1 := pack.RotorInletTriangle.W()
	w2 := pack.W2
	u1 := pack.U1
	u2 := pack.U2

	t2Func := func(currT2 float64) (float64, error) {
		cp := gases.CpMean(node.gas(), t1, currT2, nodes.DefaultN)
		newT2 := t1 + ((w1*w1-w2*w2)+(u2*u2-u1*u1))/(2*cp)
		return newT2, nil
	}
	t2, err := common.SolveIterativelyWithValidation(t2Func, common2.NotNanValidator, t1, node.precision, 1, nodes.DefaultN)
	if err != nil {
		pack.Err = fmt.Errorf("%s: t2: thermo2", err.Error())
		return
	}
	pack.checkAllAndSet(
		t2,
		&pack.T2,
	)

	cp := gases.CpMean(node.gas(), pack.T1, pack.T2, nodes.DefaultN)
	pack.checkAllAndSet(
		pack.T1-pack.RotorHeatDrop/cp,
		&pack.T2Prime,
	)

	k := gases.KMean(node.gas(), pack.T2Prime, pack.T1, nodes.DefaultN)
	pack.checkAllAndSet(
		pack.P1*math.Pow(pack.T2Prime/pack.T1, k/(k-1)),
		&pack.P2,
	)
	pack.checkAllAndSet(
		pack.P2/(node.gas().R()*pack.T2),
		&pack.Density2,
	)
}

func (node *turbineStageNode) relativeThermo(pack *DataPack) {
	if pack.setChapterName("relativeThermo") {
		return
	}
	pack.RotorInletTriangle = states2.NewInletTriangle(pack.U1, pack.C1, pack.Alpha1)

	w1 := pack.RotorInletTriangle.W()
	t1 := pack.T1
	pack.checkAllAndSet(
		t1+w1*w1/(2*node.gas().Cp(t1)), // todo check if using correct cp
		&pack.Tw1,
	)

	k := gases.K(node.gas(), t1) // todo check if using correct cp
	p1 := pack.P1
	tw1 := pack.Tw1

	pack.checkAllAndSet(
		p1*math.Pow(tw1/t1, k/(k-1)),
		&pack.Pw1,
	)
	pack.checkAllAndSet(
		node.stageHeatDrop*node.reactivity*t1/tw1,
		&pack.RotorHeatDrop,
	)

	rotorGeom := pack.StageGeometry.RotorGeometry()
	d := rotorGeom.MeanProfile().Diameter(rotorGeom.XGapOut())
	pack.checkAllAndSet(
		math.Pi*d*node.n/60,
		&pack.U2,
	)

	hl := pack.RotorHeatDrop
	u1 := pack.U1
	u2 := pack.U2
	pack.checkAllAndSet(
		math.Sqrt(w1*w1+2*hl+(u2*u2-u1*u1)),
		&pack.WAd2,
	)
	pack.checkAllAndSet(
		pack.WAd2*node.psi,
		&pack.W2,
	)
}

func (node *turbineStageNode) velocity1FirstStage(pack *DataPack) {
	if pack.setChapterName("velocity1FirstStage") {
		return
	}
	pack.checkAllAndSet(
		node.massRate()/(pack.C1a*pack.Density1),
		&pack.Area1,
	)

	lRelOut := node.stageGeomGen.StatorGenerator().LRelOut()
	pack.checkAllAndSet(
		math.Sqrt(pack.Area1/(math.Pi*lRelOut)),
		&pack.RotorMeanInletDiameter,
	)

	pack.StageGeometry = node.stageGeomGen.GenerateFromRotorInlet(
		pack.RotorMeanInletDiameter,
	)
	node.u1(pack)

	// initialize velocity input when using first stage model
	// currently it sets velocity in axial direction
	massRate := node.massRate()
	area := geometry.Area(0, pack.StageGeometry.StatorGeometry())
	density := node.p0Stag() / (node.gas().R() * node.t0Stag()) // todo use static density instead of stag
	ca := massRate / (area * density)
	node.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, ca, math.Pi/2),
		states2.InletTriangleType,
	))
}

func (node *turbineStageNode) velocity1(pack *DataPack) {
	if pack.setChapterName("velocity1") {
		return
	}
	pack.checkAllAndSet(
		geometry.Area(
			pack.StageGeometry.StatorGeometry().XGapOut(),
			pack.StageGeometry.StatorGeometry(),
		),
		&pack.Area1,
	)
	pack.checkAllAndSet(
		node.massRate()/(pack.Area1*pack.Density1),
		&pack.C1a,
	)
	node.u1(pack)

	alpha1 := math.Asin(pack.C1a / pack.C1)
	if math.IsNaN(alpha1) {
		pack.Err = fmt.Errorf("failed to calculate alpha_1 (c_a_1 = %v, c1 = %v)", pack.C1a, pack.C1)
		return
	}
	pack.checkAllAndSet(alpha1, &pack.Alpha1)
}

func (node *turbineStageNode) u1(pack *DataPack) {
	if pack.setChapterName("u1") {
		return
	}
	d := pack.StageGeometry.RotorGeometry().MeanProfile().Diameter(
		pack.StageGeometry.RotorGeometry().XBladeIn(),
	)
	pack.checkAllAndSet(
		math.Pi*d*node.n/60,
		&pack.U1,
	)
}

func (node *turbineStageNode) thermo1(pack *DataPack) {
	if pack.setChapterName("thermo1") {
		return
	}

	pack.checkAllAndSet(
		node.stageHeatDrop*(1-node.reactivity),
		&pack.StatorHeatDrop,
	)
	pack.checkAllAndSet(
		math.Sqrt(2*pack.StatorHeatDrop),
		&pack.C1Ad,
	)
	pack.checkAllAndSet(
		pack.C1Ad*node.phi,
		&pack.C1,
	)
	if node.mode == first {
		node.c1aFirstStage(pack)
	}

	t0Stag := node.t0Stag()
	cp := node.gas().Cp(t0Stag)
	hc := pack.StatorHeatDrop
	pack.checkAllAndSet(
		t0Stag-hc/cp,
		&pack.T1Prime,
	)

	t1Func := func(t1 float64) (float64, error) {
		t0Stag := node.t0Stag()
		c1 := pack.C1
		cp := gases.CpMean(node.gas(), t1, node.t0Stag(), nodes.DefaultN)
		result := t0Stag - c1*c1/(2*cp)
		return result, nil
	}
	t1, err := common.SolveIterativelyWithValidation(
		t1Func, common2.NotNanValidator, node.t0Stag(), node.precision, 1, nodes.DefaultN,
	)
	if err != nil {
		pack.Err = fmt.Errorf("%s: t1: thermo1", err.Error())
		return
	}
	pack.checkAllAndSet(t1, &pack.T1)

	k := gases.KMean(node.gas(), node.t0Stag(), pack.T1Prime, nodes.DefaultN)
	pack.checkAllAndSet(
		node.p0Stag()*math.Pow(pack.T1Prime/t0Stag, k/(k-1)),
		&pack.P1,
	)
	pack.checkAllAndSet(
		pack.P1/(node.gas().R()*pack.T1),
		&pack.Density1,
	)
}

func (node *turbineStageNode) c1aFirstStage(pack *DataPack) {
	if pack.setChapterName("c1a") {
		return
	}

	if math.IsNaN(node.alpha1FirstStage) {
		pack.Err = fmt.Errorf("alpha1 not set for first stage")
	}
	pack.checkAllAndSet(
		pack.C1*math.Sin(node.alpha1FirstStage),
		&pack.C1a,
	)
}

func (node *turbineStageNode) getStageGeometry(pack *DataPack) {
	if pack.setChapterName("StageGeometry") {
		return
	}

	gammaIn := node.stageGeomGen.StatorGenerator().GammaIn()
	gammaOut := node.stageGeomGen.StatorGenerator().GammaOut()
	baRel := node.stageGeomGen.StatorGenerator().Elongation()
	deltaRel := node.stageGeomGen.StatorGenerator().DeltaRel()
	xRel := -(1 + deltaRel) / baRel
	lRelOut := node.stageGeomGen.StatorGenerator().LRelOut()

	lRelIn := RecalculateLRel(lRelOut, xRel, gammaIn, gammaOut)

	c0 := node.statorInletTriangle().C()

	pack.checkAllAndSet(
		math.Sqrt(node.massRate()/(math.Pi*pack.Density0*c0*lRelIn)),
		&pack.StatorMeanInletDiameter,
	)
	pack.StageGeometry = node.stageGeomGen.GenerateFromStatorInlet(
		pack.StatorMeanInletDiameter,
	)
}

func (node *turbineStageNode) thermo0(pack *DataPack) {
	if pack.setChapterName("thermo0") {
		return
	}
	c0 := node.statorInletTriangle().C()
	cp := node.gas().Cp(node.t0Stag()) // todo check if correct temperature
	pack.checkAllAndSet(node.t0Stag()-c0*c0/(2*cp), &pack.T0)

	k := gases.K(node.gas(), node.t0Stag()) // todo check if correct temperature
	pack.checkAllAndSet(
		node.p0Stag()*math.Pow(node.t0Stag()/pack.T0, -k/(k-1)),
		&pack.P0,
	)

	pack.checkAllAndSet(
		pack.P0/(node.gas().R()*pack.T0),
		&pack.Density0,
	)
}

// below are private accessors

func (node *turbineStageNode) massRate() float64 {
	return node.MassRateInput().GetState().(states.MassRatePortState).MassRate
}

func (node *turbineStageNode) p0Stag() float64 {
	return node.PressureInput().GetState().(states.PressurePortState).PStag
}

func (node *turbineStageNode) t0Stag() float64 {
	return node.TemperatureInput().GetState().(states.TemperaturePortState).TStag
}

func (node *turbineStageNode) gas() gases.Gas {
	return node.GasInput().GetState().(states.GasPortState).Gas
}

func (node *turbineStageNode) statorInletTriangle() states2.VelocityTriangle {
	return node.VelocityInput().GetState().(states2.VelocityPortState).Triangle
}
