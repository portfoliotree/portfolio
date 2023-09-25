package allocation

import (
	"context"
	"math"
	"time"

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

	assetRisks := assetReturns.RisksFromStdDev()
	for i := range assetRisks {
		assetRisks[i] = 1.0 / math.Pow(assetRisks[i], 2)
	}

	sumOfAssetRisks := 0.0
	for i := range assetRisks {
		sumOfAssetRisks += assetRisks[i]
	}

	newWeights := make([]float64, len(assetRisks))
	for i := range assetRisks {
		newWeights[i] = assetRisks[i] / sumOfAssetRisks
	}

	return newWeights, nil
}
