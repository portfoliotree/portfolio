package allocation

import (
	"context"
	"time"

	"github.com/portfoliotree/portfolio/calculations"
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
	calculations.EqualVolatilityWeights(ws, assetRisks)
	return ws, nil
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

	vols := assetReturns.RisksFromStdDev()
	calculations.EqualInverseVolatilityWeights(ws, vols)
	return ws, nil
}
