package constructive

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"gonum.org/v1/gonum/mat"
)

type GasCombiner interface {
	graph.Node
	MainInput() nodes.ComplexGasSink
	ExtraInput() nodes.ComplexGasSink
	Output() nodes.ComplexGasSource
}

func NewGasCombiner(precision, relaxCoef float64, iterLimit int) GasCombiner {
	result := &gasCombiner{
		precision: precision,
		relaxCoef: relaxCoef,
		iterLimit: iterLimit,
	}
	graph.AttachAllWithTags(
		result,
		[]*graph.Port{
			&result.gInput, &result.tInput, &result.pInput, &result.mrInput,
			&result.gEInput, &result.tEInput, &result.pEInput, &result.mrEInput,
			&result.gOutput, &result.tOutput, &result.pOutput, &result.mrOutput,
		},
		[]string{
			"gInput", "tInput", "pInput", "mrInput",
			"gEInput", "tEInput", "pEInput", "mrEInput",
			"gOutput", "tOutput", "pOutput", "mrOutput",
		},
	)
	return result
}

type gasCombiner struct {
	graph.BaseNode

	precision float64
	relaxCoef float64
	iterLimit int

	tInput  graph.Port
	pInput  graph.Port
	gInput  graph.Port
	mrInput graph.Port

	tEInput  graph.Port
	pEInput  graph.Port
	gEInput  graph.Port
	mrEInput graph.Port

	tOutput  graph.Port
	pOutput  graph.Port
	gOutput  graph.Port
	mrOutput graph.Port
}

func (node *gasCombiner) GetName() string {
	return common.EitherString(node.GetInstanceName(), "GasCombiner")
}

func (node *gasCombiner) GetPorts() []graph.Port {
	return []graph.Port{
		node.gInput, node.tInput, node.pInput, node.mrInput,
		node.gEInput, node.tEInput, node.pEInput, node.mrEInput,
		node.gOutput, node.tOutput, node.pOutput, node.mrOutput,
	}
}

func (node *gasCombiner) GetRequirePorts() ([]graph.Port, error) {
	return []graph.Port{
		node.gInput, node.tInput, node.pInput, node.mrInput,
		node.gEInput, node.tEInput, node.pEInput, node.mrEInput,
	}, nil
}

func (node *gasCombiner) GetUpdatePorts() ([]graph.Port, error) {
	return []graph.Port{
		node.gOutput, node.tOutput, node.pOutput, node.mrOutput,
	}, nil
}

// pressure is taken from main input
func (node *gasCombiner) Process() error {
	mainGas := node.gInput.GetState().Value().(gases.Gas)
	mainT := node.tInput.GetState().Value().(float64)
	mainMR := node.mrInput.GetState().Value().(float64)

	extraGas := node.gEInput.GetState().Value().(gases.Gas)
	extraT := node.tEInput.GetState().Value().(float64)
	extraMR := node.mrEInput.GetState().Value().(float64)

	mainFraction := mainMR / (mainMR + extraMR)
	extraFraction := extraMR / (mainMR + extraMR)

	oGas := gases.NewMixture(
		[]gases.Gas{mainGas, extraGas},
		[]float64{mainFraction, extraFraction},
	)

	eqSys := math.NewEquationSystem(func(tOutVec *mat.VecDense) (*mat.VecDense, error) {
		tOut := tOutVec.At(0, 0)

		cpMain := gases.CpMean(mainGas, mainT, tOut, nodes.DefaultN)
		mainHeat := mainMR * cpMain * (tOut - mainT)

		cpExtra := gases.CpMean(extraGas, extraT, tOut, nodes.DefaultN)
		extraHeat := extraMR * cpExtra * (tOut - extraT)

		return mat.NewVecDense(1, []float64{mainHeat + extraHeat}), nil
	}, 1)

	solver, err := newton.NewUniformNewtonSolver(eqSys, 1e-5, newton.NoLog)
	if err != nil {
		return err
	}

	solution, err := solver.Solve(
		mat.NewVecDense(1, []float64{mainT*mainFraction + extraT*extraFraction}),
		node.precision, node.relaxCoef, node.iterLimit,
	)
	if err != nil {
		return err
	}

	oT := solution.At(0, 0)

	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(oGas),
			states.NewTemperaturePortState(oT),
			node.pInput.GetState(),
			states.NewMassRatePortState(mainMR + extraMR),
		},
		[]graph.Port{
			node.gOutput, node.tOutput, node.pOutput, node.mrOutput,
		},
	)
	return nil
}

func (node *gasCombiner) MainInput() nodes.ComplexGasSink {
	return helper.NewPseudoComplexGasSink(
		node.gInput, node.tInput, node.pInput, node.mrInput,
	)
}

func (node *gasCombiner) ExtraInput() nodes.ComplexGasSink {
	return helper.NewPseudoComplexGasSink(
		node.gEInput, node.tEInput, node.pEInput, node.mrEInput,
	)
}

func (node *gasCombiner) Output() nodes.ComplexGasSource {
	return helper.NewPseudoComplexGasSource(
		node.gOutput, node.tOutput, node.pOutput, node.mrOutput,
	)
}
