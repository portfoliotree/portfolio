package calculate

import (
	"math"
)

// HoldingPeriodReturns calculates the holding period returns between the given quotes
// The Return's Time field remains the zero value.
func HoldingPeriodReturns(quotes []float64) []float64 {
	if len(quotes) < 2 {
		return nil
	}
	result := make([]float64, len(quotes)-1)
	for i := 0; i < len(quotes)-1; i++ {
		result[i] = quotes[i]/quotes[i+1] - 1
	}
	return result
}

func productOfReturnValues(returns []float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	product := 1.0

	for i := len(returns) - 1; i >= 0; i-- {
		product *= 1.0 + returns[i]
	}

	return product
}

func TimeWeightedReturn(returns []float64) float64 {
	return productOfReturnValues(returns) - 1
}

func SharpeRatio(portfolioReturnValues, riskFreeReturnValues []float64, periods float64) float64 {
	portfolioReturn := AnnualizedArithmeticReturn(portfolioReturnValues, periods)
	portfolioRisk := RiskFromStdDev(portfolioReturnValues)
	riskFreeReturn := AnnualizedTimeWeightedReturn(riskFreeReturnValues, periods)
	return (portfolioReturn - riskFreeReturn) / portfolioRisk
}

const MinArithmeticReturn = -0.9999

// CAGRFromArithmeticReturn calculates the ex ante compound return based on:
// Mindlin, Dimitry, On the Relationship between Arithmetic and Geometric Returns (August 14, 2011).
// Available at SSRN: https://ssrn.com/abstract=2083915 or http://dx.doi.org/10.2139/ssrn.2083915
func CAGRFromArithmeticReturn(arithmeticReturn, volatility float64) float64 {
	if arithmeticReturn <= MinArithmeticReturn {
		return 0
	}
	portfolioVariance := math.Pow(volatility, 2)
	compoundTerm := 1.0 + portfolioVariance*math.Pow(1.0+arithmeticReturn, -2.0)
	geometricReturn := (1.0+arithmeticReturn)*math.Pow(compoundTerm, -0.5) - 1.0
	return geometricReturn
}

// ArithmeticReturnFromCAGR calculates the ex ante arithmetic return based on:
// Mindlin, Dimitry, On the Relationship between Arithmetic and Geometric Returns (August 14, 2011).
// Available at SSRN: https://ssrn.com/abstract=2083915 or http://dx.doi.org/10.2139/ssrn.2083915
func ArithmeticReturnFromCAGR(geometricReturn, volatility float64) float64 {
	x1 := geometricReturn
	x2 := volatility * volatility
	return (1+x1)*math.Sqrt(1.0/2+(1.0/2*math.Sqrt(1+(4*x2)/(1+x1)/(1+x1)))) - 1
}
