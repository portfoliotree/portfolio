package allocation

import (
	"context"
	"time"

	"github.com/portfoliotree/timetable"

	"github.com/portfoliotree/portfolio/calculate"
)

const EqualWeightsAlgorithmName = "Equal Weights"

type EqualWeights struct{}

func (*EqualWeights) Name() string { return EqualWeightsAlgorithmName }

func (*EqualWeights) PolicyWeights(_ context.Context, _ time.Time, _ timetable.Compact[float64], ws []float64) ([]float64, error) {
	calculate.EqualWeights(ws)
	return ws, nil
}
