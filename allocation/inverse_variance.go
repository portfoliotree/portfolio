package allocation

import (
	"context"
	"time"

	"github.com/portfoliotree/portfolio/calculate"
	"github.com/portfoliotree/portfolio/returns"
)

type EqualInverseVariance struct{}

func (cw *EqualInverseVariance) Name() string { return "Equal Inverse Variance" }

func (*EqualInverseVariance) PolicyWeights(_ context.Context, _ time.Time, assetReturns returns.Table, ws []float64) ([]float64, error) {
	if isOnlyZeros(ws) {
		for i := range ws {
			ws[i] = 1.0
		}
		scaleToUnitRange(ws)
	}

	err := ensureEnoughReturns(assetReturns)
	if err != nil {
		return ws, err
	}

	vols := assetReturns.RisksFromStdDev()
	calculate.InverseVarianceWeights(ws, vols)
	return ws, nil
}
