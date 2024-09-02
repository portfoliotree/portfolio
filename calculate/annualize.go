package calculate

import (
	"math"

	"gonum.org/v1/gonum/stat"
)

const (
	PeriodsPerYear = 252.0
)

func AnnualizeRisk(risk, periodsPerYear float64) float64 {
	return risk * math.Sqrt(periodsPerYear)
}

func AnnualizedTimeWeightedReturn(returnValues []float64, periods float64) float64 {
	if len(returnValues) < 2 {
		return 0
	}
	totalNumberOfPeriods := float64(len(returnValues))
	return math.Pow(productOfReturnValues(returnValues), periods/totalNumberOfPeriods) - 1
}

// AnnualizedArithmeticReturn must receive at least 2 returns otherwise it returns 0
func AnnualizedArithmeticReturn(returnValues []float64, periods float64) float64 {
	if len(returnValues) < 2 {
		return 0
	}
	arithmeticReturn := stat.Mean(returnValues, nil)
	return arithmeticReturn * periods
}
