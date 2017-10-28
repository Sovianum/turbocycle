package nodes

import (
	"fmt"
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/helpers/gases"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	states2 "github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/go-errors/errors"
)

type TurbineStageNode interface {
	core.Node
	nodes.GasChannel
	nodes.PressureChannel
	nodes.TemperatureChannel
	VelocityChannel
	MassRateChannel
}

func NewTurbineStageNode(
	n, stageHeatDrop, reactivity, phi, psi, airGapRel, precision float64,
	gen geometry.StageGeometryGenerator,
) TurbineStageNode {
	var result = &turbineStageNode{
		ports:         make(core.PortsType),
		n:             n,
		stageHeatDrop: stageHeatDrop,
		reactivity:    reactivity,
		phi:           phi,
		psi:           psi,
		airGapRel:     airGapRel,
		precision:     precision,
		stageGeomGen:  gen,
	}
	result.ports[nodes.GasInput] = core.NewPort()
	result.ports[nodes.GasInput].SetInnerNode(result)
	result.ports[nodes.GasInput].SetState(states.NewGasPortState(gases.GetAir()))

	result.ports[nodes.GasOutput] = core.NewPort()
	result.ports[nodes.GasOutput].SetInnerNode(result)
	result.ports[nodes.GasOutput].SetState(states.NewGasPortState(gases.GetAir()))

	result.ports[nodes.PressureInput] = core.NewPort()
	result.ports[nodes.PressureInput].SetInnerNode(result)
	result.ports[nodes.PressureInput].SetState(states.NewPressurePortState(common.AtmPressure))

	result.ports[nodes.PressureOutput] = core.NewPort()
	result.ports[nodes.PressureOutput].SetInnerNode(result)
	result.ports[nodes.PressureOutput].SetState(states.NewPressurePortState(common.AtmPressure))

	result.ports[nodes.TemperatureInput] = core.NewPort()
	result.ports[nodes.TemperatureInput].SetInnerNode(result)
	result.ports[nodes.TemperatureInput].SetState(states.NewTemperaturePortState(common.AtmTemperature))

	result.ports[nodes.TemperatureOutput] = core.NewPort()
	result.ports[nodes.TemperatureOutput].SetInnerNode(result)
	result.ports[nodes.TemperatureOutput].SetState(states.NewTemperaturePortState(common.AtmTemperature))

	result.ports[VelocityInput] = core.NewPort()
	result.ports[VelocityInput].SetInnerNode(result)
	result.ports[VelocityInput].SetState(
		states2.NewVelocityPortState(
			states2.NewInletTriangle(0, 0, math.Pi/2), states2.InletTriangleType,
		),
	)

	result.ports[velocityOutput] = core.NewPort()
	result.ports[velocityOutput].SetInnerNode(result)
	result.ports[velocityOutput].SetState(
		states2.NewVelocityPortState(
			states2.NewInletTriangle(0, 0, math.Pi/2), states2.InletTriangleType,
		),
	)

	result.ports[massRateInput] = core.NewPort()
	result.ports[massRateInput].SetInnerNode(result)
	result.ports[massRateInput].SetState(states2.NewMassRatePortState(0))

	result.ports[massRateOutput] = core.NewPort()
	result.ports[massRateOutput].SetInnerNode(result)
	result.ports[massRateOutput].SetState(states2.NewMassRatePortState(0))

	return result
}

type turbineStageNode struct {
	ports         core.PortsType
	n             float64
	stageHeatDrop float64
	reactivity    float64
	phi           float64
	psi           float64
	airGapRel     float64

	stageGeomGen  geometry.StageGeometryGenerator
	stageGeometry geometry.StageGeometry

	precision float64
}

type dataPack struct {
	err error

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

func (node *turbineStageNode) MarshalJSON() ([]byte, error) {
	return nil, nil // todo add real functional
}

func (node *turbineStageNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *turbineStageNode) Process() error {
	return nil
}

func (node *turbineStageNode) GetRequirePortTags() ([]string, error) {
	return []string{
		nodes.TemperatureInput,
		nodes.PressureInput,
		nodes.GasInput,
		massRateInput,
		dimensionInput,
	}, nil
}

func (node *turbineStageNode) GetUpdatePortTags() ([]string, error) {
	return []string{
		nodes.TemperatureOutput,
		nodes.PressureOutput,
		nodes.GasOutput,
		massRateOutput,
		dimensionOutput,
	}, nil
}

func (node *turbineStageNode) GetPortTags() []string {
	return []string{
		nodes.TemperatureInput, nodes.TemperatureOutput,
		nodes.PressureInput, nodes.PressureOutput,
		nodes.GasInput, nodes.GasOutput,
		massRateInput, massRateOutput,
		dimensionInput, dimensionOutput,
	}
}

func (node *turbineStageNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.TemperatureInput:
		return node.temperatureInput(), nil
	case nodes.TemperatureOutput:
		return node.temperatureOutput(), nil
	case nodes.PressureInput:
		return node.pressureInput(), nil
	case nodes.PressureOutput:
		return node.pressureOutput(), nil
	case nodes.GasInput:
		return node.gasInput(), nil
	case nodes.GasOutput:
		return node.gasOutput(), nil
	case massRateInput:
		return node.massRateInput(), nil
	case massRateOutput:
		return node.massRateOutput(), nil
	case VelocityInput:
		return node.velocityInput(), nil
	case velocityOutput:
		return node.velocityOutput(), nil
	default:
		return nil, fmt.Errorf("port with tag \"%s\" not found", tag)
	}
}

func (node *turbineStageNode) ContextDefined() bool {
	return true
}

func (node *turbineStageNode) GasOutput() core.Port {
	return node.gasOutput()
}

func (node *turbineStageNode) GasInput() core.Port {
	return node.gasInput()
}

func (node *turbineStageNode) VelocityInput() core.Port {
	return node.velocityInput()
}

func (node *turbineStageNode) VelocityOutput() core.Port {
	return node.velocityOutput()
}

func (node *turbineStageNode) PressureOutput() core.Port {
	return node.pressureOutput()
}

func (node *turbineStageNode) PressureInput() core.Port {
	return node.pressureInput()
}

func (node *turbineStageNode) TemperatureOutput() core.Port {
	return node.temperatureOutput()
}

func (node *turbineStageNode) TemperatureInput() core.Port {
	return node.temperatureInput()
}

func (node *turbineStageNode) MassRateInput() core.Port {
	return node.massRateInput()
}

func (node *turbineStageNode) MassRateOutput() core.Port {
	return node.massRateOutput()
}

func (node *turbineStageNode) getDataPack() *dataPack {
	var pack = new(dataPack)

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
	node.etaT(pack)

	return pack
}

func (node *turbineStageNode) etaTStag(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: etaTStag", pack.err.Error())
		return
	}
	pack.EtaTStag = pack.StageLabour / node.stageHeatDrop
}

func (node *turbineStageNode) stageHeatDropStag(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: stageHeatDropStag", pack.err.Error())
		return
	}
	var cp = gases.CpMean(node.gas(), node.t0Stag(), pack.T2, nodes.DefaultN)
	var k = gases.KMean(node.gas(), node.t0Stag(), pack.T2, nodes.DefaultN)
	pack.StageHeatDropStag = cp * node.t0Stag() * (1 - math.Pow(pack.P2Stag/node.p0Stag(), (k-1)/k))
}

func (node *turbineStageNode) stageLabour(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: stageLabour", pack.err.Error())
		return
	}
	pack.StageLabour = node.stageHeatDrop * pack.EtaT
}

func (node *turbineStageNode) p2Stag(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: p2Stag", pack.err.Error())
		return
	}
	var k = gases.K(node.gas(), pack.T2) // todo check if correct temperature
	pack.P2Stag = pack.P2 * math.Pow(pack.T2Stag/pack.T2, k/(k-1))
}

func (node *turbineStageNode) t2Stag(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: t2Stag", pack.err.Error())
		return
	}
	var cp = node.gas().Cp(pack.T2) // todo check if correct temperature
	pack.T2Stag = pack.T2 +
		(pack.AirGapSpecificLoss+pack.VentilationSpecificLoss+pack.OutletVelocitySpecificLoss)/cp
}

func (node *turbineStageNode) etaT(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: etaT", pack.err.Error())
		return
	}
	pack.EtaT = pack.EtaU - (pack.AirGapSpecificLoss+pack.VentilationSpecificLoss)/node.stageHeatDrop
}

func (node *turbineStageNode) ventilationSpecificLoss(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: ventilationSpecificLoss", pack.err.Error())
		return
	}
	var x = pack.StageGeometry.RotorGeometry().XBladeOut()
	var d = pack.StageGeometry.RotorGeometry().MeanProfile().Diameter(x)
	var ventilationPower = 1.07 * math.Pow(d, 2) * math.Pow(pack.U2/100, 3) * pack.Density2 * 1000
	pack.VentilationSpecificLoss = ventilationPower / node.massRate()
}

func (node *turbineStageNode) airGapSpecificLoss(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: airGapSpecificLoss", pack.err.Error())
		return
	}
	var lRel = geometry.RelativeHeight(
		pack.StageGeometry.RotorGeometry().XBladeOut(),
		pack.StageGeometry.RotorGeometry(),
	)
	pack.AirGapSpecificLoss = 1.37 * (1 + 1.6*node.reactivity) * (1 + lRel) * node.airGapRel
}

func (node *turbineStageNode) outletVelocitySpecificLoss(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: outletVelocitySpecificLoss", pack.err.Error())
		return
	}
	pack.OutletVelocitySpecificLoss = math.Pow(pack.RotorOutletTriangle.C(), 2) / 2
}

func (node *turbineStageNode) rotorSpecificLoss(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: rotorSpecificLoss", pack.err.Error())
		return
	}
	pack.RotorSpecificLoss = (math.Pow(node.psi, -2) - 1) * math.Pow(pack.W2, 2) / 2
}

func (node *turbineStageNode) statorSpecificLoss(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: statorSpecificLoss", pack.err.Error())
		return
	}
	var defaultSpecificLoss = (math.Pow(node.phi, -2) - 1) * math.Pow(pack.C1, 2) / 2
	pack.StatorSpecificLoss = defaultSpecificLoss * pack.T2Prime / pack.T1
}

func (node *turbineStageNode) etaU(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: etaU", pack.err.Error())
		return
	}
	pack.EtaU = pack.MeanRadiusLabour / node.stageHeatDrop
}

func (node *turbineStageNode) meanRadiusLabour(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: meanRadiusLabour", pack.err.Error())
		return
	}
	pack.MeanRadiusLabour = pack.RotorInletTriangle.CU()*pack.RotorInletTriangle.U() +
		pack.RotorOutletTriangle.CU()*pack.RotorOutletTriangle.U()
}

func (node *turbineStageNode) pi(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: pi", pack.err.Error())
		return
	}
	pack.Pi = node.p0Stag() / pack.P2
}

func (node *turbineStageNode) rotorOutletTriangle(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: rotorOutletTriangle", pack.err.Error())
		return
	}
	pack.RotorOutletTriangle = states2.NewOutletTriangleFromProjections(
		pack.C2u, pack.C2a, pack.U2,
	)
}

func (node *turbineStageNode) c2u(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: c2u", pack.err.Error())
		return
	}
	pack.C2u = pack.W2*math.Cos(pack.Beta2) - pack.U2
}

func (node *turbineStageNode) beta2(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: beta2", pack.err.Error())
		return
	}
	pack.Beta2 = math.Asin(pack.C2a / pack.W2)
}

func (node *turbineStageNode) c2a(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: c2a", pack.err.Error())
		return
	}
	pack.C2a = node.massRate() / (pack.Density2 * geometry.Area(pack.StageGeometry.RotorGeometry().XBladeOut(), pack.StageGeometry.RotorGeometry()))
}

func (node *turbineStageNode) density2(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: density2", pack.err.Error())
		return
	}
	pack.Density2 = pack.P2 / (node.gas().R() * pack.T2)
}

func (node *turbineStageNode) p2(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: p2", pack.err.Error())
		return
	}
	var k = gases.KMean(node.gas(), pack.T2Prime, pack.T1, nodes.DefaultN)

	pack.P2 = pack.P1 * math.Pow(pack.T2Prime/pack.T1, k/(k-1))
}

func (node *turbineStageNode) t2Prime(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: t2Prime", pack.err.Error())
		return
	}
	var cp = gases.CpMean(node.gas(), pack.T1, pack.T2, nodes.DefaultN)
	pack.T2Prime = pack.T1 - pack.RotorHeatDrop/cp
}

func (node *turbineStageNode) t2(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: t2", pack.err.Error())
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
		newT2, newCp = iterate(currT2, currCp)
	}

	pack.T2 = newT2
}

func (node *turbineStageNode) w2(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: w2", pack.err.Error())
		return
	}
	var wAd2 = pack.WAd2
	pack.W2 = wAd2 * node.psi
}

func (node *turbineStageNode) wAd2(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: wAd2", pack.err.Error())
		return
	}
	var w1 = pack.RotorInletTriangle.W()
	var hl = pack.RotorHeatDrop
	var u1 = pack.U1
	var u2 = pack.U2
	pack.WAd2 = math.Sqrt(w1*w1 + 2*hl + (u2*u2 - u1*u1))
}

func (node *turbineStageNode) rotorHeatDrop(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: rotorHeatDrop", pack.err.Error())
		return
	}
	var t1 = pack.T1
	var tw1 = pack.Tw1
	pack.RotorHeatDrop = node.stageHeatDrop * node.reactivity * t1 / tw1
}

func (node *turbineStageNode) pw1(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: pw1", pack.err.Error())
		return
	}
	var t1 = pack.T1
	var k = gases.K(node.gas(), t1) // todo check if using correct cp
	var p1 = pack.P1
	var tw1 = pack.Tw1
	pack.Pw1 = p1 * math.Pow(tw1/t1, k/(k-1))
}

func (node *turbineStageNode) tw1(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: tw1", pack.err.Error())
		return
	}

	var w1 = pack.RotorInletTriangle.W()
	var t1 = pack.T1
	pack.Tw1 = t1 + w1*w1/(2*node.gas().Cp(t1)) // todo check if using correct cp
}

func (node *turbineStageNode) u2(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: u2", pack.err.Error())
		return
	}
	var rotorGeom = pack.StageGeometry.RotorGeometry()
	var d = rotorGeom.MeanProfile().Diameter(rotorGeom.XGapOut())
	pack.U2 = math.Pi * d * node.n / 60
}

func (node *turbineStageNode) rotorInletTriangle(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: rotorInletTriangle", pack.err.Error())
		return
	}
	pack.RotorInletTriangle = states2.NewInletTriangle(pack.U1, pack.C1, pack.Alpha1)
}

func (node *turbineStageNode) alpha1(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: alpha1", pack.err.Error())
		return
	}
	pack.Alpha1 = math.Asin(pack.C1a / pack.C1)
}

func (node *turbineStageNode) u1(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: u1", pack.err.Error())
		return
	}
	pack.U1 = math.Pi * pack.StageGeometry.RotorGeometry().MeanProfile().Diameter(
		pack.StageGeometry.RotorGeometry().XBladeIn(),
	) * node.n / 60
}

func (node *turbineStageNode) c1a(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: c1a", pack.err.Error())
		return
	}
	pack.C1a = node.massRate() / (pack.Area1 * pack.Density1)
}

func (node *turbineStageNode) area1(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: area1", pack.err.Error())
		return
	}
	pack.Area1 = geometry.Area(pack.StageGeometry.StatorGeometry().XGapOut(), pack.StageGeometry.StatorGeometry())
}

func (node *turbineStageNode) density1(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: density1", pack.err.Error())
		return
	}
	pack.Density1 = pack.P1 / (node.gas().R() * pack.T1)
}

func (node *turbineStageNode) p1(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: p1", pack.err.Error())
		return
	}

	var k = gases.KMean(node.gas(), node.t0Stag(), pack.T1Prime, nodes.DefaultN)
	var p0Stag = node.p0Stag()
	var t1Prime = pack.T1Prime
	var t0Stag = node.t0Stag()

	pack.P1 = p0Stag * math.Pow(t1Prime/t0Stag, k/(k-1))
}

func (node *turbineStageNode) t1(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: t1", pack.err.Error())
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
			pack.err = errors.New("failed to converge: try different initial guess")
			return
		}
		t1Curr = t1New
		t1New = getNewT1(t1Curr)
	}

	pack.T1 = t1New
}

func (node *turbineStageNode) c1(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: c1", pack.err.Error())
		return
	}
	pack.C1 = pack.C1Ad * node.phi
}

func (node *turbineStageNode) c1Ad(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: c1Ad", pack.err.Error())
		return
	}
	pack.C1Ad = math.Sqrt(2 * pack.StatorHeatDrop)
}

func (node *turbineStageNode) t1Prime(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: t1Prime", pack.err.Error())
		return
	}

	var t0Stag = node.t0Stag()
	var cp = node.gas().Cp(t0Stag)
	var hc = pack.StatorHeatDrop
	pack.T1Prime = t0Stag - hc/cp
}

func (node *turbineStageNode) statorHeatDrop(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: statorHeatDrop", pack.err.Error())
		return
	}
	pack.StatorHeatDrop = node.stageHeatDrop * (1 - node.reactivity)
}

func (node *turbineStageNode) getStageGeometry(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: getStageGeometry", pack.err.Error())
		return
	}
	pack.StageGeometry = node.stageGeomGen.GenerateFromStatorInlet(pack.StatorMeanInletDiameter)
}

func (node *turbineStageNode) getStatorMeanInletDiameter(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: getStatorMeanDiameter", pack.err.Error())
		return
	}

	var baRel = node.stageGeomGen.StatorGenerator().Elongation()
	var gammaIn = node.stageGeomGen.StatorGenerator().GammaIn()
	var gammaOut = node.stageGeomGen.StatorGenerator().GammaOut()
	var _, gammaMean = geometry.GetTotalAndMeanLineAngles(
		node.stageGeomGen.StatorGenerator().GammaIn(),
		node.stageGeomGen.StatorGenerator().GammaOut(),
	)
	var lRelOut = node.stageGeomGen.StatorGenerator().LRelOut()
	var lRelIn = lRelOut * (baRel - (math.Tan(gammaOut) - math.Tan(gammaIn))) / (baRel - 2*lRelOut*math.Tan(gammaMean))

	var c0 = node.statorInletTriangle().C()
	pack.StatorMeanInletDiameter = math.Sqrt(node.massRate() / (math.Pi * pack.Density0 * c0 * lRelIn))
}

func (node *turbineStageNode) density0(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: density0", pack.err.Error())
		return
	}
	pack.Density0 = pack.P0 / (node.gas().R() * pack.T0)
}

func (node *turbineStageNode) p0(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: p0", pack.err.Error())
		return
	}
	var k = gases.K(node.gas(), node.t0Stag()) // todo check if correct temperature
	pack.P0 = node.p0Stag() * math.Pow(node.t0Stag()/pack.T0, -k/(k-1))
}

func (node *turbineStageNode) t0(pack *dataPack) {
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: t0", pack.err.Error())
		return
	}
	var c0 = node.statorInletTriangle().C()
	var cp = node.gas().Cp(node.t0Stag()) // todo check if correct temperature
	pack.T0 = node.t0Stag() - c0*c0/(2*cp)
}

func (node *turbineStageNode) massRate() float64 {
	return node.massRateInput().GetState().(states2.MassRatePortState).MassRate
}

func (node *turbineStageNode) p0Stag() float64 {
	return node.pressureInput().GetState().(states.PressurePortState).PStag
}

func (node *turbineStageNode) t0Stag() float64 {
	return node.temperatureInput().GetState().(states.TemperaturePortState).TStag
}

func (node *turbineStageNode) gas() gases.Gas {
	return node.gasInput().GetState().(states.GasPortState).Gas
}

func (node *turbineStageNode) statorInletTriangle() states2.VelocityTriangle {
	return node.velocityInput().GetState().(states2.VelocityPortState).Triangle
}

func (node *turbineStageNode) velocityInput() core.Port {
	return node.ports[VelocityInput]
}

func (node *turbineStageNode) velocityOutput() core.Port {
	return node.ports[velocityOutput]
}

func (node *turbineStageNode) gasOutput() core.Port {
	return node.ports[nodes.GasOutput]
}

func (node *turbineStageNode) gasInput() core.Port {
	return node.ports[nodes.GasInput]
}

func (node *turbineStageNode) pressureOutput() core.Port {
	return node.ports[nodes.PressureOutput]
}

func (node *turbineStageNode) pressureInput() core.Port {
	return node.ports[nodes.PressureInput]
}

func (node *turbineStageNode) temperatureOutput() core.Port {
	return node.ports[nodes.TemperatureOutput]
}

func (node *turbineStageNode) temperatureInput() core.Port {
	return node.ports[nodes.TemperatureInput]
}

func (node *turbineStageNode) massRateInput() core.Port {
	return node.ports[massRateInput]
}

func (node *turbineStageNode) massRateOutput() core.Port {
	return node.ports[massRateOutput]
}
