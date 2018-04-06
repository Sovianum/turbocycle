package compressor

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
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

func GetCycleFitEqSys(
	staged StagedCompressorNode, simple constructive.CompressorNode,
	htDistribGen, etaDistribGen common.FuncGen1D,
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
	return GetCompressorPiEtaEqSys(staged, htDistribGen, simple.PiStag(), etaDistribGen, simple.Eta())
}

func GetCompressorPiEtaEqSys(
	compressor StagedCompressorNode,
	htDistribGen common.FuncGen1D, targetPi float64,
	etaDistribGen common.FuncGen1D, targetEta float64,
) math.EquationSystem {
	return math.NewEquationSystem(func(v *mat.VecDense) (*mat.VecDense, error) {
		htParameter := v.At(0, 0)
		etaParameter := v.At(1, 0)
		compressor.SetHtLaw(common.FromDistribution(htDistribGen(htParameter)))
		compressor.SetEtaAdLaw(common.FromDistribution(etaDistribGen(etaParameter)))
		if err := compressor.Process(); err != nil {
			return nil, err
		}
		return mat.NewVecDense(2, []float64{
			PiStag(compressor) - targetPi,
			EtaStag(compressor) - targetEta,
		}), nil
	}, 2)
}
