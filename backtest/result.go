package backtest

import (
	"time"

	"github.com/portfoliotree/portfolio/returns"
)

type Result struct {
	ReturnsTable       returns.Table `json:"returnsTable"`
	Weights            [][]float64   `json:"historicalWeights"`
	FinalPolicyWeights []float64     `json:"finalPolicyWeights"`
	RebalanceTimes     []time.Time   `json:"rebalanceDates"`
	PolicyUpdateTimes  []time.Time   `json:"policyUpdatesDates"`
}

func (result Result) Returns() returns.List {
	return result.ReturnsTable.List(0)
}

func (result Result) DailyRebalancedReturns() returns.List {
	return result.ReturnsTable.List(1)
}
