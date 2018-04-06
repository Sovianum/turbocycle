package turbine

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/stage/common"
	"gonum.org/v1/gonum/mat"
)

func GetCycleFitEqSys(
	staged StagedTurbineNode, simple constructive.TurbineNode,
	phiDistribGen, psiDistribGen common.FuncGen1D,
) math.EquationSystem {
	graph.CopyAll(
		[]graph.Port{
			simple.GasInput(), simple.TemperatureInput(),
			simple.PressureInput(),
		},
		[]graph.Port{
			staged.GasInput(), staged.TemperatureInput(),
			staged.PressureInput(),
		},
	)
	return GetTurbinePiEtaEqSys(staged, phiDistribGen, psiDistribGen, simple.PiTStag(), simple.Eta())
}

// efficiency is changed by phi and distribution only
func GetTurbinePiEtaEqSys(
	turbine StagedTurbineNode,
	phiDistribGen, psiDistribGen common.FuncGen1D,
	targetPi float64, targetEta float64,
) math.EquationSystem {
	return math.NewEquationSystem(func(v *mat.VecDense) (*mat.VecDense, error) {
		ht := v.At(0, 0)
		efficiencyFitParameter := v.At(1, 0)
		turbine.SetHt(ht)
		turbine.SetPhiFunc(phiDistribGen(efficiencyFitParameter))
		turbine.SetPsiFunc(psiDistribGen(efficiencyFitParameter))
		if err := turbine.Process(); err != nil {
			return nil, err
		}
		return mat.NewVecDense(2, []float64{
			PiStag(turbine) - targetPi,
			EtaStag(turbine) - targetEta,
		}), nil
	}, 2)
}
