package compressor

import (
	"github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/impl/stage/common"
	"gonum.org/v1/gonum/mat"
)

func GetPiFitEqSys1D(
	compressor StagedCompressorNode, pi float64,
	funcSetter func(common.Func1D), fg common.FuncGen1D,
) math.EquationSystem {
	sysCall := func() (*mat.VecDense, error) {
		if err := compressor.Process(); err != nil {
			return nil, err
		}
		return mat.NewVecDense(1, []float64{pi - PiStag(compressor)}), nil
	}
	return common.GetEqSys1D(sysCall, funcSetter, fg)
}
