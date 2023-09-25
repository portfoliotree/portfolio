package allocation

import (
	"context"
	"time"

	"github.com/portfoliotree/portfolio/returns"
)

const EqualWeightsAlgorithmName = "Equal Weights"

type EqualWeights struct{}

func (*EqualWeights) Name() string { return EqualWeightsAlgorithmName }

func (*EqualWeights) PolicyWeights(_ context.Context, _ time.Time, _ returns.Table, ws []float64) ([]float64, error) {
	for i := range ws {
		ws[i] = 1.0 / float64(len(ws))
	}
	return ws, nil
}
