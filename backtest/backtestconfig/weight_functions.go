package backtestconfig

import (
	"context"
	"time"

	"github.com/portfoliotree/portfolio/returns"
)

// PolicyWeightCalculatorFunc can be used to wrap a function and pass it into Run as a PolicyWeightCalculator
type PolicyWeightCalculatorFunc func(ctx context.Context, today time.Time, assets returns.Table, currentWeights []float64) ([]float64, error)

func (p PolicyWeightCalculatorFunc) PolicyWeights(ctx context.Context, today time.Time, assets returns.Table, currentWeights []float64) ([]float64, error) {
	return p(ctx, today, assets, currentWeights)
}

type ConstantWeights []float64

func (targetWeights ConstantWeights) PolicyWeights(_ context.Context, _ time.Time, _ returns.Table, _ []float64) ([]float64, error) {
	return targetWeights, nil
}

type EqualWeights struct{}

func (EqualWeights) PolicyWeights(_ context.Context, _ time.Time, assets returns.Table, _ []float64) ([]float64, error) {
	n := float64(assets.NumberOfColumns())
	targetWeights := make([]float64, assets.NumberOfColumns())
	for i := range targetWeights {
		targetWeights[i] = 1 / n
	}
	return targetWeights, nil
}

// Additional weight functions are maintained in portfoliotree.com proprietary code.
// If you'd like to read the code, feel free to ask us at support@portfoliotree.com, we are willing to share pseudocode.
