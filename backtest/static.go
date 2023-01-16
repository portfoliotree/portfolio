package backtest

import (
	"errors"

	"github.com/portfoliotree/portfolio/returns"
)

func DailyRebalancedWithStaticWeights(assets returns.Table, weights []float64) (returns.List, error) {
	if assets.NumberOfColumns() != len(weights) {
		return nil, errors.New("the number of weights and assets must be the same")
	}
	n := assets.NumberOfRows()

	result := make(returns.List, n)

	for i := 0; i < n; i++ {
		for ai := 0; ai < assets.NumberOfColumns(); ai++ {
			ar := assets.List(ai)
			result[i].Time = ar[i].Time
			result[i].Value += ar[i].Value * weights[ai]
		}
	}

	return result, nil
}
