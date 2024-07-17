package calculate

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
	if len(ws) != len(vols) {
		panic("length of weights and volatilizes must be equal")
	}

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
	if len(ws) != len(vols) {
		panic("length of weights and volatilizes must be equal")
	}
	if len(ws) != len(correlations) {
		panic("length of weights and correlations must be equal")
	}
	for _, row := range correlations {
		if len(row) != len(correlations) {
			panic("correlations must be a square matrix")
		}
	}

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
	if len(ws) != len(vols) {
		panic("length of weights and volatilizes must be equal")
	}

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
	if len(ws) != len(vols) {
		panic("length of weights and volatilizes must be equal")
	}

	sum := 0.0
	for i := range vols {
		sum += vols[i]
	}
	for i := range vols {
		ws[i] = vols[i] / sum
	}
}
