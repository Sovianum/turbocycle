package variator

import (
	"fmt"

	"github.com/Sovianum/turbocycle/core/graph"
	"gonum.org/v1/gonum/mat"
)

type SysCall func() (*mat.VecDense, error)

func SysCallFromNetwork(
	network graph.Network, vectorPort graph.Port,
	relaxCoef float64, skipIterations, iterLimit int,
	precision float64,
) SysCall {
	return func() (*mat.VecDense, error) {
		var converged, err = network.Solve(relaxCoef, skipIterations, iterLimit, precision)
		if err != nil {
			return nil, err
		}

		if !converged {
			return nil, fmt.Errorf("failed to converge")
		}

		var vec = vectorPort.GetState().(graph.VectorPortState).Vec
		return vec, nil
	}
}
