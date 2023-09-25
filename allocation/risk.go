package allocation

import (
	"context"
	"math"
	"time"

	"github.com/portfoliotree/portfolio/calculations"
	"github.com/portfoliotree/portfolio/returns"
)

type EqualRiskContribution struct{}

func (*EqualRiskContribution) Name() string { return "Equal Risk Contribution" }

func (*EqualRiskContribution) PolicyWeights(ctx context.Context, _ time.Time, assetReturns returns.Table, ws []float64) ([]float64, error) {
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

	target := 1.0 / float64(len(assetRisks))

	cm := assetReturns.CorrelationMatrix()

	weights := make([]float64, len(ws))
	copy(weights, ws)

	return weights, optWeights(ctx, weights, func(ws []float64) float64 {
		_, _, riskWeights := calculations.RiskFromRiskContribution(assetRisks, ws, cm)
		var diff float64
		for i := range riskWeights {
			diff += math.Abs(target - riskWeights[i])
		}
		return diff
	})
}
