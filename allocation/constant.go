package allocation

import (
	"context"
	"errors"
	"time"

	"github.com/portfoliotree/portfolio/returns"
)

type ConstantWeights struct {
	weights []float64
}

func (cw *ConstantWeights) Name() string { return "Constant Weights" }

func (cw *ConstantWeights) PolicyWeights(_ context.Context, _ time.Time, _ returns.Table, ws []float64) ([]float64, error) {
	if len(cw.weights) != len(ws) {
		return nil, errors.New("expected the number of policy weights to be the same as the number of assets")
	}
	copy(ws, cw.weights)
	return ws, nil
}

func (cw *ConstantWeights) SetWeights(in []float64) {
	cw.weights = in
}
