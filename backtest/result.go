package backtest

import (
	"time"

	"github.com/portfoliotree/portfolio/returns"
)

const (
	PortfolioReturnsColumn       = 0
	DailyRebalancedReturnsColumn = 1
)

type Result struct {
	ReturnsTable       returns.Table `json:"returnsTable"       bson:"returnsTable"`
	Weights            [][]float64   `json:"weights"            bson:"weights"`
	FinalPolicyWeights []float64     `json:"policyWeights"      bson:"policyWeights"`
	RebalanceTimes     []time.Time   `json:"rebalanceDates"     bson:"rebalanceDates"`
	PolicyUpdateTimes  []time.Time   `json:"policyUpdatesDates" bson:"policyUpdatesDates"`
}

func (result Result) Returns() returns.List {
	return result.ReturnsTable.List(PortfolioReturnsColumn)
}

func (result Result) DailyRebalancedReturns() returns.List {
	return result.ReturnsTable.List(DailyRebalancedReturnsColumn)
}
