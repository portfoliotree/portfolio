package calculate

import (
	"errors"
	"math"
	"slices"

	"gonum.org/v1/gonum/stat"
)

func Variance(volatility float64) float64 {
	return volatility * volatility
}

func NumberOfBets(weightedAverageRisk float64, portfolioRisk float64) (float64, error) {
	if portfolioRisk == 0 {
		return 0, errors.New("can't divide by 0")
	}

	abRatio := weightedAverageRisk / portfolioRisk
	bets := abRatio * abRatio

	return bets, nil
}

func RiskFromStdDev(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	// PopStdDev may be the correct method to call here for backwards looking stuff.
	// For now, we use sample standard deviation.
	return stat.StdDev(values, nil)
}

func WeightedAverageRisk(weights, risks []float64) float64 {
	return stat.Mean(risks, weights)
}

func RiskWeights(portfolioVol float64, riskContributions []float64) []float64 {
	result := slices.Clone(riskContributions)
	for i := range result {
		result[i] = result[i] / portfolioVol
	}
	return result
}

func PortfolioVolatility(weights, stdDevs []float64, correlations [][]float64) (float64, []float64) {
	var variance float64
	n := len(weights)
	covariances := calculateCovarianceMatrix(weights, stdDevs, correlations)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			variance += covariances[i][j]
		}
	}

	totalRisk := math.Sqrt(variance)
	riskContributions := make([]float64, n)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			riskContributions[i] += covariances[i][j]
		}
		riskContributions[i] = riskContributions[i] / totalRisk
	}

	return totalRisk, riskContributions
}

func calculateCovarianceMatrix(weights, vols []float64, correlations [][]float64) [][]float64 {
	n := len(weights)
	covariances := make([][]float64, n)
	for i := range covariances {
		covariances[i] = make([]float64, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			covariances[i][j] = weights[i] * weights[j] * vols[i] * vols[j] * correlations[i][j]
		}
	}
	return covariances
}
