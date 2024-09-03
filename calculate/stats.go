package calculate

import (
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

func DownsideVolatility(values []float64, periodsPerYear float64) float64 {
	negative := make([]float64, 0, len(values))
	negative = negativeValues(negative, values)
	for i := range negative {
		negative[i] = math.Pow(negative[i], 2)
	}
	x := stat.Mean(negative, nil)
	return math.Sqrt(x) * math.Sqrt(periodsPerYear)
}

func SortinoRatio(portfolio, riskFree []float64, downsideVol, periodsPerYear float64) float64 {
	pr := AnnualizedArithmeticReturn(portfolio, periodsPerYear)
	rr := AnnualizedArithmeticReturn(riskFree, periodsPerYear)
	return (pr - rr) / downsideVol
}

func MaxDrawdown(values []float64) (float64, int) {
	retained := retainedAfterDrawdown(values)
	i := floats.MinIdx(retained)
	ret := retained[i]
	return 1 - ret, i
}

func CalmarRatio(portfolio, riskFree []float64, maxDrawdown, periodsPerYear float64) float64 {
	pr := AnnualizedArithmeticReturn(portfolio, periodsPerYear)
	rr := AnnualizedArithmeticReturn(riskFree, periodsPerYear)
	return (pr - rr) / maxDrawdown
}

func UlcerIndex(values []float64, periodsPerYear float64) float64 {
	retained := retainedAfterDrawdown(values)
	for i := range retained {
		retained[i] = math.Pow(retained[i], 2)
	}
	x := stat.Mean(retained, nil)
	return math.Sqrt(x) * math.Sqrt(periodsPerYear)
}

func TrackingError(excessReturns []float64, periodsPerYear float64) float64 {
	return stat.PopStdDev(excessReturns, nil) * math.Sqrt(periodsPerYear)
}

func InformationRatio(portfolio, benchmark []float64, periodsPerYear float64) float64 {
	pr := AnnualizedTimeWeightedReturn(portfolio, periodsPerYear)
	br := AnnualizedTimeWeightedReturn(benchmark, periodsPerYear)
	er := pr - br

	excess := make([]float64, len(portfolio))
	floats.SubTo(excess, portfolio, benchmark)
	te := TrackingError(excess, periodsPerYear)

	return er / te
}

func BetaToBenchmark(portfolio, benchmark []float64) float64 {
	_, slope := stat.LinearRegression(benchmark, portfolio, nil, false)
	return slope
}

func ValueAtRisk(values []float64, portfolioValue, confidenceLevel, periodsPerYear float64) float64 {
	normal := distuv.Normal{
		Mu:    0,
		Sigma: 1.0,
	}
	zScore := normal.Quantile(confidenceLevel)
	return portfolioValue * -zScore * AnnualizeRisk(stat.PopStdDev(values, nil), periodsPerYear)
}

func negativeValues(out, in []float64) []float64 {
	for _, v := range in {
		if v < 0 {
			out = append(out, v)
		}
	}
	return out
}

func retainedAfterDrawdown(values []float64) []float64 {
	cum := make([]float64, len(values))
	for i := range values {
		index := len(values) - 1 - i
		if i == 0 {
			cum[index] = 1 + values[index]
		} else {
			cum[index] = cum[index+1] * (1 + values[index])
		}
	}
	retained := make([]float64, len(values))
	for i := range cum {
		index := len(values) - 1 - i
		if index == len(cum)-1 {
			retained[index] = 1
			continue
		}
		x := floats.Max(cum[index:])
		if x < 1 {
			x = 1
		}
		retained[index] = cum[index] / x
	}
	return retained
}
