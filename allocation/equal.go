package allocation

import (
	"context"
	"time"

	"github.com/portfoliotree/portfolio/calculations"
	"github.com/portfoliotree/portfolio/returns"
)

const EqualWeightsAlgorithmName = "Equal Weights"

type EqualWeights struct{}

func (*EqualWeights) Name() string { return EqualWeightsAlgorithmName }

func (*EqualWeights) PolicyWeights(_ context.Context, _ time.Time, _ returns.Table, ws []float64) ([]float64, error) {
	calculations.EqualWeights(ws)
	return ws, nil
}
