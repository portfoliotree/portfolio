package backtest

import (
	"errors"

	"github.com/portfoliotree/timetable"
)

func DailyRebalancedWithStaticWeights(assets timetable.Compact[float64], weights []float64) (timetable.List[float64], error) {
	if assets.NumberOfColumns() != len(weights) {
		return nil, errors.New("the number of weights and assets must be the same")
	}
	var (
		times  = assets.UnderlyingTimes()
		result = make(timetable.List[float64], 0, len(times))
		values = assets.UnderlyingValues()
	)
	for i := 0; i < len(times); i++ {
		var v float64
		for j := 0; j < len(values); j++ {
			v += values[j][i] * weights[j]
		}
		result = append(result, timetable.NewCell(times[i], v))
	}
	return result, nil
}
