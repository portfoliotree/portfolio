package allocation

import (
	"context"
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

	err = calculations.EqualRiskContributionWeights(ctx, ws, assetReturns.RisksFromStdDev(), assetReturns.CorrelationMatrix())
	return ws, err
}
