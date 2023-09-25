package allocation

import (
	"context"
	"time"

	"github.com/portfoliotree/portfolio/returns"
)

const ConstantWeightsAlgorithmName = "Constant Weights"

type ConstantWeights struct {
	weights []float64
}

func (cw *ConstantWeights) Name() string { return ConstantWeightsAlgorithmName }

func (cw *ConstantWeights) PolicyWeights(_ context.Context, _ time.Time, _ returns.Table, ws []float64) ([]float64, error) {
	copy(ws, cw.weights)
	return ws, nil
}

func (cw *ConstantWeights) SetWeights(in []float64) {
	cw.weights = in
}
