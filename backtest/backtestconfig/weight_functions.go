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

func (targetWeights ConstantWeights) PolicyWeights(_ context.Context, _ time.Time, _ returns.Table, ws []float64) ([]float64, error) {
	copy(ws, targetWeights)
	return ws, nil
}

type EqualWeights struct{}

func (EqualWeights) PolicyWeights(_ context.Context, _ time.Time, _ returns.Table, ws []float64) ([]float64, error) {
	for i := range ws {
		ws[i] = 1.0 / float64(len(ws))
	}
	return ws, nil
}

// Additional weight functions are maintained in portfoliotree.com proprietary code.
// If you'd like to read the code, feel free to ask us at support@portfoliotree.com, we are willing to share pseudocode.
