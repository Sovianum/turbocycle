package turbine

import (
	"github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/impl/stage/common"
	"gonum.org/v1/gonum/mat"
)

// efficiency is changed by phi distribution only
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
