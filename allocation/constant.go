package allocation

import (
	"context"
	"errors"
	"time"

	"github.com/portfoliotree/timetable"
)

const ConstantWeightsAlgorithmName = "Constant Weights"

type ConstantWeights struct {
	weights []float64
}

func (cw *ConstantWeights) Name() string { return ConstantWeightsAlgorithmName }

func (cw *ConstantWeights) PolicyWeights(_ context.Context, _ time.Time, _ timetable.Compact[float64], ws []float64) ([]float64, error) {
	if len(cw.weights) != len(ws) {
		return nil, errors.New("expected the number of policy weights to be the same as the number of assets")
	}
	copy(ws, cw.weights)
	return ws, nil
}

func (cw *ConstantWeights) SetWeights(in []float64) {
	cw.weights = in
}
