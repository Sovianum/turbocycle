package compressor

import (
	"fmt"

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
	htLimit, etaLimit float64,
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
	return GetCompressorPiEtaEqSys(
		staged,
		htDistribGen, htLimit, simple.PiStag(),
		etaDistribGen, etaLimit, simple.Eta(),
	)
}

func GetCompressorPiEtaEqSys(
	compressor StagedCompressorNode,
	htDistribGen common.FuncGen1D, htLimit, targetPi float64,
	etaDistribGen common.FuncGen1D, etaLimit, targetEta float64,
) math.EquationSystem {
	return math.NewEquationSystem(func(v *mat.VecDense) (*mat.VecDense, error) {
		htParameter := v.At(0, 0)
		etaParameter := v.At(1, 0)

		htLaw := htDistribGen(htParameter)
		etaLaw := etaDistribGen(etaParameter)

		compressor.SetHtLaw(common.FromDistribution(htLaw))
		compressor.SetEtaAdLaw(common.FromDistribution(etaLaw))
		if err := compressor.Process(); err != nil {
			return nil, err
		}

		for i := range compressor.Stages() {
			floatI := float64(i)
			eta := etaLaw(floatI)
			ht := htLaw(floatI)
			if eta > etaLimit {
				return nil, fmt.Errorf("eta exceeded limit %.3f > %.3f on stage %d", eta, etaLimit, i)
			}
			if ht > htLimit {
				return nil, fmt.Errorf("ht exceeded limit %.3f > %.3f on stage %d", ht, htLimit, i)
			}
		}

		currPi := PiStag(compressor)
		currEta := EtaStag(compressor)
		return mat.NewVecDense(2, []float64{
			currPi - targetPi,
			currEta - targetEta,
		}), nil
	}, 2)
}
