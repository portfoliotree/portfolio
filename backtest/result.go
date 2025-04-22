package backtest

import (
	"time"

	"github.com/portfoliotree/timetable"
)

const (
	PortfolioReturnsColumn       = 0
	DailyRebalancedReturnsColumn = 1
)

type Result struct {
	ReturnsTable       timetable.Compact[float64] `json:"returnsTable"       bson:"returnsTable"`
	Weights            [][]float64                `json:"weights"            bson:"weights"`
	FinalPolicyWeights []float64                  `json:"policyWeights"      bson:"policyWeights"`
	RebalanceTimes     []time.Time                `json:"rebalanceDates"     bson:"rebalanceDates"`
	PolicyUpdateTimes  []time.Time                `json:"policyUpdatesDates" bson:"policyUpdatesDates"`
}

func (result Result) Returns() timetable.List[float64] {
	list, ok := result.ReturnsTable.Column(PortfolioReturnsColumn)
	if !ok {
		return timetable.List[float64](nil)
	}
	return list
}

func (result Result) DailyRebalancedReturns() timetable.List[float64] {
	list, ok := result.ReturnsTable.Column(DailyRebalancedReturnsColumn)
	if !ok {
		return timetable.List[float64](nil)
	}
	return list
}
