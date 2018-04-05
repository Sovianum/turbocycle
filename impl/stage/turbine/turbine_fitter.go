package turbine

import (
	"github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/impl/stage/common"
	"gonum.org/v1/gonum/mat"
)

func GetPiFitEqSys1D(
	turbine StagedTurbineNode, pi float64,
	funcSetter func(common.Func1D), fg common.FuncGen1D,
) math.EquationSystem {
	sysCall := func() (*mat.VecDense, error) {
		if err := turbine.Process(); err != nil {
			return nil, err
		}
		return mat.NewVecDense(1, []float64{pi - PiStag(turbine)}), nil
	}
	return common.GetEqSys1D(sysCall, funcSetter, fg)
}
