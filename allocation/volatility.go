package allocation

import (
	"context"
	"time"

	"github.com/portfoliotree/portfolio/returns"
)

type EqualVolatility struct{}

func (*EqualVolatility) Name() string { return "Equal Volatility" }

func (*EqualVolatility) PolicyWeights(_ context.Context, _ time.Time, assetReturns returns.Table, ws []float64) ([]float64, error) {
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

type EqualInverseVolatility struct{}

func (*EqualInverseVolatility) Name() string { return "Equal Inverse Volatility" }

func (*EqualInverseVolatility) PolicyWeights(_ context.Context, _ time.Time, assetReturns returns.Table, ws []float64) ([]float64, error) {
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
		assetRisks[i] = 1.0 / assetRisks[i]
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
