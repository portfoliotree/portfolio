package calculations

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
