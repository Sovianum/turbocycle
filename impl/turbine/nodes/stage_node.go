package nodes

import (
	"fmt"
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	states2 "github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/go-errors/errors"
)

type TurbineStageNode interface {
	graph.Node
	nodes.GasChannel
	nodes.PressureChannel
	nodes.TemperatureChannel
	VelocityChannel
	MassRateChannel
	SetFirstStageMode(isFirstStage bool)
	SetAlpha1FirstStage(alpha1FirstStage float64)
	StageGeomGen() geometry.StageGeometryGenerator
	Ht() float64
	Reactivity() float64
	GetDataPack() DataPack
}

func InitFromTurbineNode(stage TurbineStageNode, turbine constructive.TurbineNode, massRate, alpha1 float64) {
	stage.GasInput().SetState(states.NewGasPortState(turbine.InputGas()))
	stage.TemperatureInput().SetState(states.NewTemperaturePortState(turbine.TStagIn()))
	stage.PressureInput().SetState(states.NewPressurePortState(turbine.PStagIn()))
	stage.MassRateInput().SetState(states.NewMassRatePortState(massRate))
	stage.SetAlpha1FirstStage(alpha1)
}

func NewTurbineStageNode(
	n, stageHeatDrop, reactivity, phi, psi, airGapRel, precision float64,
	gen geometry.StageGeometryGenerator,
) TurbineStageNode {
	var result = &turbineStageNode{
		n:                n,
		stageHeatDrop:    stageHeatDrop,
		reactivity:       reactivity,
		phi:              phi,
		psi:              psi,
		airGapRel:        airGapRel,
		precision:        precision,
		stageGeomGen:     gen,
		alpha1FirstStage: math.NaN(),
		isFirstStageNode: false,
	}
	result.gasInput = graph.NewAttachedPort(result)
	result.gasOutput = graph.NewAttachedPort(result)

	result.pressureInput = graph.NewAttachedPort(result)
	result.pressureOutput = graph.NewAttachedPort(result)

	result.temperatureInput = graph.NewAttachedPort(result)
	result.temperatureOutput = graph.NewAttachedPort(result)

	result.velocityInput = graph.NewAttachedPort(result)
	result.velocityInput.SetState(
		states2.NewVelocityPortState(
			states2.NewInletTriangle(0, 0, math.Pi/2), states2.InletTriangleType,
		),
	)

	result.velocityOutput = graph.NewAttachedPort(result)

	result.massRateInput = graph.NewAttachedPort(result)
	result.massRateOutput = graph.NewAttachedPort(result)

	return result
}

type turbineStageNode struct {
	graph.BaseNode

	gasInput          graph.Port
	gasOutput         graph.Port
	pressureInput     graph.Port
	pressureOutput    graph.Port
	temperatureInput  graph.Port
	temperatureOutput graph.Port
	velocityInput     graph.Port
	velocityOutput    graph.Port
	massRateInput     graph.Port
	massRateOutput    graph.Port

	n                float64
	stageHeatDrop    float64
	reactivity       float64
	phi              float64
	psi              float64
	airGapRel        float64
	alpha1FirstStage float64

	stageGeomGen geometry.StageGeometryGenerator

	precision float64

	isFirstStageNode bool
	pack             *DataPack
}

type DataPack struct {
	Err error

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

func (node *turbineStageNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "TurbineStage")
}

func (node *turbineStageNode) GetPorts() []graph.Port {
	return []graph.Port{
		node.gasInput,
		node.gasOutput,
		node.pressureInput,
		node.pressureOutput,
		node.temperatureInput,
		node.temperatureOutput,
		node.velocityInput,
		node.velocityOutput,
		node.massRateInput,
		node.massRateOutput,
	}
}

func (node *turbineStageNode) GetRequirePorts() []graph.Port {
	return []graph.Port{
		node.gasInput,
		node.pressureInput,
		node.temperatureInput,
		node.velocityInput,
		node.massRateInput,
	}
}

func (node *turbineStageNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{
		node.gasOutput,
		node.pressureOutput,
		node.temperatureOutput,
		node.velocityOutput,
		node.massRateOutput,
	}
}

func (node *turbineStageNode) Process() error {
	node.pack = node.getDataPack()
	if node.pack.Err != nil {
		return node.pack.Err
	}

	node.temperatureOutput.SetState(states.NewTemperaturePortState(node.pack.T2Stag))
	node.pressureOutput.SetState(states.NewPressurePortState(node.pack.P2Stag))
	node.massRateOutput.SetState(states.NewMassRatePortState(node.massRate())) // mass rate is constant
	node.velocityOutput.SetState(states2.NewVelocityPortState(node.pack.RotorOutletTriangle, states2.OutletTriangleType))
	return nil
}

func (node *turbineStageNode) GetDataPack() DataPack {
	if node.pack == nil {
		node.pack = node.getDataPack()
	}
	return *node.pack
}

func (node *turbineStageNode) SetFirstStageMode(isFirstStageNode bool) {
	node.isFirstStageNode = isFirstStageNode
}

func (node *turbineStageNode) SetAlpha1FirstStage(alpha1FirstStage float64) {
	node.alpha1FirstStage = alpha1FirstStage
}

func (node *turbineStageNode) GasOutput() graph.Port {
	return node.gasOutput
}

func (node *turbineStageNode) GasInput() graph.Port {
	return node.gasInput
}

func (node *turbineStageNode) VelocityInput() graph.Port {
	return node.velocityInput
}

func (node *turbineStageNode) VelocityOutput() graph.Port {
	return node.velocityOutput
}

func (node *turbineStageNode) PressureOutput() graph.Port {
	return node.pressureOutput
}

func (node *turbineStageNode) PressureInput() graph.Port {
	return node.pressureInput
}

func (node *turbineStageNode) TemperatureOutput() graph.Port {
	return node.temperatureOutput
}

func (node *turbineStageNode) TemperatureInput() graph.Port {
	return node.temperatureInput
}

func (node *turbineStageNode) MassRateInput() graph.Port {
	return node.massRateInput
}

func (node *turbineStageNode) MassRateOutput() graph.Port {
	return node.massRateOutput
}

func (node *turbineStageNode) StageGeomGen() geometry.StageGeometryGenerator {
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
	if node.isFirstStageNode {
		node.initCalc(pack)
	} else {
		node.initCalcFirstStage(pack)
	}

	node.rotorInletTriangle(pack)
	node.u2(pack)
	node.tw1(pack)
	node.pw1(pack)
	node.rotorHeatDrop(pack)
	node.wAd2(pack)
	node.w2(pack)
	node.t2(pack)
	node.t2Prime(pack)
	node.p2(pack)
	node.density2(pack)
	node.c2a(pack)
	node.beta2(pack)
	node.c2u(pack)
	node.rotorOutletTriangle(pack)
	node.pi(pack)
	node.meanRadiusLabour(pack)
	node.etaU(pack)
	node.statorSpecificLoss(pack)
	node.rotorSpecificLoss(pack)
	node.outletVelocitySpecificLoss(pack)
	node.airGapSpecificLoss(pack)
	node.ventilationSpecificLoss(pack)
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

func (node *turbineStageNode) initCalc(pack *DataPack) {
	node.t0(pack)
	node.p0(pack)
	node.density0(pack)
	node.getStatorMeanInletDiameter(pack)
	node.getStageGeometry(pack)
	node.statorHeatDrop(pack)
	node.t1Prime(pack)
	node.c1Ad(pack)
	node.c1(pack)
	node.t1(pack)
	node.p1(pack)
	node.density1(pack)
	node.area1(pack)
	node.c1a(pack)
	node.u1(pack)
	node.alpha1(pack)
}

func (node *turbineStageNode) initCalcFirstStage(pack *DataPack) {
	pack.Alpha1 = node.alpha1FirstStage
	node.statorHeatDrop(pack)
	node.t1Prime(pack)
	node.c1Ad(pack)
	node.c1(pack)
	node.c1aFirstStage(pack)
	node.t1(pack)
	node.p1(pack)
	node.density1(pack)
	node.area1FirstStage(pack)
	node.dRotorBladeMean(pack)
	node.getStageGeometryFirstStage(pack)
	node.u1(pack)
	node.inletVelocityTriangleFirstStage(pack)
}

func (node *turbineStageNode) etaTStag(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: etaTStag", pack.Err.Error())
		return
	}
	pack.EtaTStag = pack.StageLabour / pack.StageHeatDropStag
}

func (node *turbineStageNode) stageHeatDropStag(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: stageHeatDropStag", pack.Err.Error())
		return
	}
	var cp = gases.CpMean(node.gas(), node.t0Stag(), pack.T2, nodes.DefaultN)
	var k = gases.KMean(node.gas(), node.t0Stag(), pack.T2, nodes.DefaultN)
	pack.StageHeatDropStag = cp * node.t0Stag() * (1 - math.Pow(pack.P2Stag/node.p0Stag(), (k-1)/k))
}

func (node *turbineStageNode) stageLabour(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: stageLabour", pack.Err.Error())
		return
	}
	pack.StageLabour = node.stageHeatDrop * pack.EtaT
}

func (node *turbineStageNode) p2Stag(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: p2Stag", pack.Err.Error())
		return
	}
	var k = gases.K(node.gas(), pack.T2) // todo check if correct temperature
	pack.P2Stag = pack.P2 * math.Pow(pack.T2Stag/pack.T2, k/(k-1))
}

func (node *turbineStageNode) t2Stag(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: t2Stag", pack.Err.Error())
		return
	}
	var cp = node.gas().Cp(pack.T2) // todo check if correct temperature
	pack.T2Stag = pack.T2 +
		(pack.AirGapSpecificLoss+pack.VentilationSpecificLoss+pack.OutletVelocitySpecificLoss)/cp
}

func (node *turbineStageNode) etaT(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: etaT", pack.Err.Error())
		return
	}
	pack.EtaT = pack.EtaU - (pack.AirGapSpecificLoss+pack.VentilationSpecificLoss)/node.stageHeatDrop
}

func (node *turbineStageNode) ventilationSpecificLoss(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: ventilationSpecificLoss", pack.Err.Error())
		return
	}
	var x = pack.StageGeometry.RotorGeometry().XBladeOut()
	var d = pack.StageGeometry.RotorGeometry().MeanProfile().Diameter(x)
	var ventilationPower = 1.07 * math.Pow(d, 2) * math.Pow(pack.U2/100, 3) * pack.Density2 * 1000
	pack.VentilationSpecificLoss = ventilationPower / node.massRate()
}

func (node *turbineStageNode) airGapSpecificLoss(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: airGapSpecificLoss", pack.Err.Error())
		return
	}
	var lRel = geometry.RelativeHeight(
		pack.StageGeometry.RotorGeometry().XBladeOut(),
		pack.StageGeometry.RotorGeometry(),
	)
	pack.AirGapSpecificLoss = 1.37 * (1 + 1.6*node.reactivity) * (1 + lRel) * node.airGapRel * pack.MeanRadiusLabour
}

func (node *turbineStageNode) outletVelocitySpecificLoss(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: outletVelocitySpecificLoss", pack.Err.Error())
		return
	}
	pack.OutletVelocitySpecificLoss = math.Pow(pack.RotorOutletTriangle.C(), 2) / 2
}

func (node *turbineStageNode) rotorSpecificLoss(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: rotorSpecificLoss", pack.Err.Error())
		return
	}
	pack.RotorSpecificLoss = (math.Pow(node.psi, -2) - 1) * math.Pow(pack.W2, 2) / 2
}

func (node *turbineStageNode) statorSpecificLoss(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: statorSpecificLoss", pack.Err.Error())
		return
	}
	var defaultSpecificLoss = (math.Pow(node.phi, -2) - 1) * math.Pow(pack.C1, 2) / 2
	pack.StatorSpecificLoss = defaultSpecificLoss * pack.T2Prime / pack.T1
}

func (node *turbineStageNode) etaU(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: etaU", pack.Err.Error())
		return
	}
	pack.EtaU = pack.MeanRadiusLabour / node.stageHeatDrop
}

func (node *turbineStageNode) meanRadiusLabour(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: meanRadiusLabour", pack.Err.Error())
		return
	}
	pack.MeanRadiusLabour = pack.RotorInletTriangle.CU()*pack.RotorInletTriangle.U() +
		pack.RotorOutletTriangle.CU()*pack.RotorOutletTriangle.U()
}

func (node *turbineStageNode) pi(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: pi", pack.Err.Error())
		return
	}
	pack.Pi = node.p0Stag() / pack.P2
}

func (node *turbineStageNode) rotorOutletTriangle(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: rotorOutletTriangle", pack.Err.Error())
		return
	}
	pack.RotorOutletTriangle = states2.NewOutletTriangleFromProjections(
		pack.C2u, pack.C2a, pack.U2,
	)
}

func (node *turbineStageNode) c2u(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: c2u", pack.Err.Error())
		return
	}
	pack.C2u = pack.W2*math.Cos(pack.Beta2) - pack.U2
}

func (node *turbineStageNode) beta2(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: beta2", pack.Err.Error())
		return
	}
	pack.Beta2 = math.Asin(pack.C2a / pack.W2)
}

func (node *turbineStageNode) c2a(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: c2a", pack.Err.Error())
		return
	}
	pack.C2a = node.massRate() / (pack.Density2 * geometry.Area(
		pack.StageGeometry.RotorGeometry().XBladeOut(), pack.StageGeometry.RotorGeometry(),
	))
}

func (node *turbineStageNode) density2(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: density2", pack.Err.Error())
		return
	}
	pack.Density2 = pack.P2 / (node.gas().R() * pack.T2)
}

func (node *turbineStageNode) p2(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: p2", pack.Err.Error())
		return
	}
	var k = gases.KMean(node.gas(), pack.T2Prime, pack.T1, nodes.DefaultN)

	pack.P2 = pack.P1 * math.Pow(pack.T2Prime/pack.T1, k/(k-1))
}

func (node *turbineStageNode) t2Prime(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: t2Prime", pack.Err.Error())
		return
	}
	var cp = gases.CpMean(node.gas(), pack.T1, pack.T2, nodes.DefaultN)
	pack.T2Prime = pack.T1 - pack.RotorHeatDrop/cp
}

func (node *turbineStageNode) t2(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: t2", pack.Err.Error())
		return
	}
	var t1 = pack.T1
	var w1 = pack.RotorInletTriangle.W()
	var w2 = pack.W2
	var u1 = pack.U1
	var u2 = pack.U2

	var iterate = func(currT2, currCp float64) (newT2, newCp float64) {
		newT2 = t1 + ((w1*w1-w2*w2)+(u2*u2-u1*u1))/(2*currCp)
		newCp = gases.CpMean(node.gas(), t1, newT2, nodes.DefaultN)
		return
	}

	var currT2 = t1
	var currCp = node.gas().Cp(t1)
	var newT2, newCp = iterate(currT2, currCp)

	for !(common.Converged(currT2, newT2, node.precision) && common.Converged(currCp, newCp, node.precision)) {
		currT2, currCp = newT2, newCp
		newT2, newCp = iterate(currT2, currCp)
	}

	pack.T2 = newT2
}

func (node *turbineStageNode) w2(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: w2", pack.Err.Error())
		return
	}
	var wAd2 = pack.WAd2
	pack.W2 = wAd2 * node.psi
}

func (node *turbineStageNode) wAd2(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: wAd2", pack.Err.Error())
		return
	}
	var w1 = pack.RotorInletTriangle.W()
	var hl = pack.RotorHeatDrop
	var u1 = pack.U1
	var u2 = pack.U2
	pack.WAd2 = math.Sqrt(w1*w1 + 2*hl + (u2*u2 - u1*u1))
}

func (node *turbineStageNode) rotorHeatDrop(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: rotorHeatDrop", pack.Err.Error())
		return
	}
	var t1 = pack.T1
	var tw1 = pack.Tw1
	pack.RotorHeatDrop = node.stageHeatDrop * node.reactivity * t1 / tw1
}

func (node *turbineStageNode) pw1(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: pw1", pack.Err.Error())
		return
	}
	var t1 = pack.T1
	var k = gases.K(node.gas(), t1) // todo check if using correct cp
	var p1 = pack.P1
	var tw1 = pack.Tw1
	pack.Pw1 = p1 * math.Pow(tw1/t1, k/(k-1))
}

func (node *turbineStageNode) tw1(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: tw1", pack.Err.Error())
		return
	}

	var w1 = pack.RotorInletTriangle.W()
	var t1 = pack.T1
	pack.Tw1 = t1 + w1*w1/(2*node.gas().Cp(t1)) // todo check if using correct cp
}

func (node *turbineStageNode) u2(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: u2", pack.Err.Error())
		return
	}
	var rotorGeom = pack.StageGeometry.RotorGeometry()
	var d = rotorGeom.MeanProfile().Diameter(rotorGeom.XGapOut())
	pack.U2 = math.Pi * d * node.n / 60
}

func (node *turbineStageNode) rotorInletTriangle(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: rotorInletTriangle", pack.Err.Error())
		return
	}
	pack.RotorInletTriangle = states2.NewInletTriangle(pack.U1, pack.C1, pack.Alpha1)
}

func (node *turbineStageNode) alpha1(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: alpha1", pack.Err.Error())
		return
	}
	var alpha1 = math.Asin(pack.C1a / pack.C1)
	if math.IsNaN(alpha1) {
		pack.Err = fmt.Errorf("failed to calculate alpha_1 (c_a_1 = %v, c1 = %v)", pack.C1a, pack.C1)
		return
	}
	pack.Alpha1 = alpha1
}

// function initializes velocity input when using first stage model
// currently it sets velocity in axial direction
func (node *turbineStageNode) inletVelocityTriangleFirstStage(pack *DataPack) {
	var massRate = node.massRate()
	var area = geometry.Area(0, pack.StageGeometry.StatorGeometry())
	var density = node.p0Stag() / (node.gas().R() * node.t0Stag()) // todo use static density instead of stag
	var ca = massRate / (area * density)
	node.velocityInput.SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, ca, math.Pi/2),
		states2.InletTriangleType,
	))
}

func (node *turbineStageNode) u1(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: u1", pack.Err.Error())
		return
	}
	pack.U1 = math.Pi * pack.StageGeometry.RotorGeometry().MeanProfile().Diameter(
		pack.StageGeometry.RotorGeometry().XBladeIn(),
	) * node.n / 60
}

func (node *turbineStageNode) c1aFirstStage(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: c1a", pack.Err.Error())
		return
	}
	if math.IsNaN(node.alpha1FirstStage) {
		pack.Err = fmt.Errorf("alpha1 not set for first stage")
	}
	pack.C1a = pack.C1 * math.Sin(node.alpha1FirstStage)
}

func (node *turbineStageNode) c1a(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: c1a", pack.Err.Error())
		return
	}
	pack.C1a = node.massRate() / (pack.Area1 * pack.Density1)
}

func (node *turbineStageNode) dRotorBladeMean(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: dRotorBladeMean", pack.Err.Error())
		return
	}
	var lRelOut = node.stageGeomGen.StatorGenerator().LRelOut()
	pack.RotorMeanInletDiameter = math.Sqrt(pack.Area1 / (math.Pi * lRelOut))
}

func (node *turbineStageNode) area1(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: area1", pack.Err.Error())
		return
	}
	pack.Area1 = geometry.Area(
		pack.StageGeometry.StatorGeometry().XGapOut(),
		pack.StageGeometry.StatorGeometry(),
	)
}

func (node *turbineStageNode) area1FirstStage(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: area1", pack.Err.Error())
		return
	}
	pack.Area1 = node.massRate() / (pack.C1a * pack.Density1)
}

func (node *turbineStageNode) density1(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: density1", pack.Err.Error())
		return
	}
	pack.Density1 = pack.P1 / (node.gas().R() * pack.T1)
}

func (node *turbineStageNode) p1(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: p1", pack.Err.Error())
		return
	}

	var k = gases.KMean(node.gas(), node.t0Stag(), pack.T1Prime, nodes.DefaultN)
	var p0Stag = node.p0Stag()
	var t1Prime = pack.T1Prime
	var t0Stag = node.t0Stag()

	pack.P1 = p0Stag * math.Pow(t1Prime/t0Stag, k/(k-1))
}

func (node *turbineStageNode) t1(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: t1", pack.Err.Error())
		return
	}

	var getNewT1 = func(t1 float64) float64 {
		var t0Stag = node.t0Stag()
		var c1 = pack.C1
		var cp = gases.CpMean(node.gas(), t1, node.t0Stag(), nodes.DefaultN)
		return t0Stag - c1*c1/(2*cp)
	}

	var t1Curr = node.t0Stag()
	var t1New = getNewT1(t1Curr)

	for !common.Converged(t1Curr, t1New, node.precision) {
		if math.IsNaN(t1Curr) || math.IsNaN(t1New) {
			pack.Err = errors.New("failed to converge: try different initial guess")
			return
		}
		t1Curr = t1New
		t1New = getNewT1(t1Curr)
	}

	pack.T1 = t1New
}

func (node *turbineStageNode) c1(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: c1", pack.Err.Error())
		return
	}
	pack.C1 = pack.C1Ad * node.phi
}

func (node *turbineStageNode) c1Ad(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: c1Ad", pack.Err.Error())
		return
	}
	pack.C1Ad = math.Sqrt(2 * pack.StatorHeatDrop)
}

func (node *turbineStageNode) t1Prime(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: t1Prime", pack.Err.Error())
		return
	}

	var t0Stag = node.t0Stag()
	var cp = node.gas().Cp(t0Stag)
	var hc = pack.StatorHeatDrop
	pack.T1Prime = t0Stag - hc/cp
}

func (node *turbineStageNode) statorHeatDrop(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: statorHeatDrop", pack.Err.Error())
		return
	}
	pack.StatorHeatDrop = node.stageHeatDrop * (1 - node.reactivity)
}

func (node *turbineStageNode) getStageGeometryFirstStage(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: getStageGeometry", pack.Err.Error())
		return
	}
	pack.StageGeometry = node.stageGeomGen.GenerateFromRotorInlet(
		pack.RotorMeanInletDiameter,
	)
}

func (node *turbineStageNode) getStageGeometry(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: getStageGeometry", pack.Err.Error())
		return
	}
	pack.StageGeometry = node.stageGeomGen.GenerateFromStatorInlet(
		pack.StatorMeanInletDiameter,
	)
}

func (node *turbineStageNode) getStatorMeanInletDiameter(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: getStatorMeanDiameter", pack.Err.Error())
		return
	}

	var baRel = node.stageGeomGen.StatorGenerator().Elongation()
	var gammaIn = node.stageGeomGen.StatorGenerator().GammaIn()
	var gammaOut = node.stageGeomGen.StatorGenerator().GammaOut()
	var _, gammaMean = geometry.GetTotalAndMeanLineAngles(
		node.stageGeomGen.StatorGenerator().GammaIn(),
		node.stageGeomGen.StatorGenerator().GammaOut(),
	)
	var deltaRel = node.stageGeomGen.StatorGenerator().DeltaRel()
	var lRelOut = node.stageGeomGen.StatorGenerator().LRelOut()
	var enom = baRel - (1+deltaRel)*(math.Tan(gammaOut)-math.Tan(gammaIn))
	var denom = baRel - 2*(1+deltaRel)*lRelOut*math.Tan(gammaMean)
	var lRelIn = enom / denom

	var c0 = node.statorInletTriangle().C()
	pack.StatorMeanInletDiameter = math.Sqrt(node.massRate() / (math.Pi * pack.Density0 * c0 * lRelIn))
}

func (node *turbineStageNode) density0(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: density0", pack.Err.Error())
		return
	}
	pack.Density0 = pack.P0 / (node.gas().R() * pack.T0)
}

func (node *turbineStageNode) p0(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: p0", pack.Err.Error())
		return
	}
	var k = gases.K(node.gas(), node.t0Stag()) // todo check if correct temperature
	pack.P0 = node.p0Stag() * math.Pow(node.t0Stag()/pack.T0, -k/(k-1))
}

func (node *turbineStageNode) t0(pack *DataPack) {
	if pack.Err != nil {
		pack.Err = fmt.Errorf("%s: t0", pack.Err.Error())
		return
	}
	var c0 = node.statorInletTriangle().C()
	var cp = node.gas().Cp(node.t0Stag()) // todo check if correct temperature
	pack.T0 = node.t0Stag() - c0*c0/(2*cp)
}

// below are private accessors

func (node *turbineStageNode) massRate() float64 {
	return node.massRateInput.GetState().(states.MassRatePortState).MassRate
}

func (node *turbineStageNode) p0Stag() float64 {
	return node.pressureInput.GetState().(states.PressurePortState).PStag
}

func (node *turbineStageNode) t0Stag() float64 {
	return node.temperatureInput.GetState().(states.TemperaturePortState).TStag
}

func (node *turbineStageNode) gas() gases.Gas {
	return node.gasInput.GetState().(states.GasPortState).Gas
}

func (node *turbineStageNode) statorInletTriangle() states2.VelocityTriangle {
	return node.velocityInput.GetState().(states2.VelocityPortState).Triangle
}
