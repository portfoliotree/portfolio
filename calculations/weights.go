package calculations

import (
	"context"
	"math"
)

func EqualWeights(ws []float64) {
	for i := range ws {
		ws[i] = 1.0 / float64(len(ws))
	}
}

func InverseVarianceWeights(ws []float64, vols []float64) {
	inverseVars := make([]float64, len(vols))
	for i := range vols {
		inverseVars[i] = 1.0 / Variance(vols[i])
	}

	sumOfInverseVars := 0.0
	for i := range inverseVars {
		sumOfInverseVars += inverseVars[i]
	}

	for i := range inverseVars {
		ws[i] = inverseVars[i] / sumOfInverseVars
	}
}

func EqualRiskContributionWeights(ctx context.Context, ws []float64, vols []float64, correlations [][]float64) error {
	target := 1.0 / float64(len(vols))
	return optWeights(ctx, ws, func(ws []float64) float64 {
		riskWeights := RiskWeights(PortfolioVolatility(vols, ws, correlations))
		var diff float64
		for i := range riskWeights {
			diff += math.Abs(target - riskWeights[i])
		}
		return diff
	})
}

func EqualInverseVolatilityWeights(ws []float64, vols []float64) {
	inverseVols := make([]float64, len(vols))
	for i := range vols {
		inverseVols[i] = 1.0 / vols[i]
	}

	sumOfInverseVars := 0.0
	for i := range inverseVols {
		sumOfInverseVars += inverseVols[i]
	}

	for i := range inverseVols {
		ws[i] = inverseVols[i] / sumOfInverseVars
	}
}

func EqualVolatilityWeights(ws []float64, vols []float64) {
	sum := 0.0
	for i := range vols {
		sum += vols[i]
	}
	for i := range vols {
		ws[i] = vols[i] / sum
	}
}
