package backtest

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/slices"

	"github.com/portfoliotree/portfolio/returns"
)

type (
	PolicyWeightCalculator interface {
		PolicyWeights(ctx context.Context, today time.Time, assets returns.Table, currentWeights []float64) ([]float64, error)
	}
	WindowFunc  func(today time.Time, table returns.Table) returns.Table
	TriggerFunc func(t time.Time, currentWeights []float64) bool
)

// Run runs a portfolio back-test. It calls function parameters for policy updates and to check
// when a policy update or rebalancing is required. Generally you should call Run or RunWithFullData.
func Run(ctx context.Context, end, start time.Time, assetReturns returns.Table,
	alg PolicyWeightCalculator,
	lookBackWindow WindowFunc,
	shouldCalculatePolicy,
	shouldRebalanceAssetWeights TriggerFunc,
) (Result, error) {
	if assetReturns.NumberOfColumns() == 0 {
		return Result{}, errors.New("no asset returns provided")
	}

	if end.After(assetReturns.LastTime()) ||
		start.Before(assetReturns.FirstTime()) {
		return Result{}, ErrorNotEnoughData{}
	}

	firstPolicyDate, policyWeights, err := fetchPolicy(ctx, end, start, alg, assetReturns, lookBackWindow)
	if err != nil {
		return Result{}, err
	}
	start = firstPolicyDate

	updatedWeights := make([]float64, assetReturns.NumberOfColumns())
	copy(updatedWeights, policyWeights)

	updatedDailyWeights := make([]float64, assetReturns.NumberOfColumns())
	copy(updatedDailyWeights, policyWeights)

	var (
		dailyRebalancedReturns returns.List
		historicReturns        returns.Table
		assetReturnValuesToday []float64
		backTestReturns        returns.List

		rebalanceCount, recalculatePolicyCount = 0, 0

		next    time.Time
		hasNext = true

		result = Result{
			Weights:            make([][]float64, 0, assetReturns.NumberOfRows()),
			RebalanceTimes:     make([]time.Time, 0, assetReturns.NumberOfRows()),
			PolicyUpdateTimes:  make([]time.Time, 0, assetReturns.NumberOfRows()),
			FinalPolicyWeights: make([]float64, assetReturns.NumberOfColumns()),
		}
	)

	for today := start; hasNext; today = next {
		next, hasNext = assetReturns.TimeAfter(today)
		scaleToUnitRange(updatedWeights)
		scaleToUnitRange(updatedDailyWeights)

		if today.After(end) {
			break
		}

		historicReturns = lookBackWindow(today, assetReturns)
		assetReturnValuesToday = historicReturns.MostRecentValues()

		if shouldCalculatePolicy(today, updatedDailyWeights) && start != today {
			var err error
			policyWeights, err = alg.PolicyWeights(ctx, today, historicReturns, updatedWeights)
			if err != nil {
				return Result{}, err
			}

			scaleToUnitRange(policyWeights)

			recalculatePolicyCount++
			result.PolicyUpdateTimes = append(result.PolicyUpdateTimes, today)
			copy(result.FinalPolicyWeights, policyWeights)
		}

		weightsToday := make([]float64, len(updatedWeights))
		copy(weightsToday, updatedWeights)
		result.Weights = append(result.Weights, weightsToday)

		ret := returns.Return{
			Time: today,
		}
		for j := 0; j < assetReturns.NumberOfColumns(); j++ {
			ret.Value += updatedWeights[j] * assetReturnValuesToday[j]
		}
		backTestReturns = append(backTestReturns, ret)

		dailyRebalancedRet := returns.Return{
			Time: today,
		}
		for j := 0; j < assetReturns.NumberOfColumns(); j++ {
			dailyRebalancedRet.Value += updatedDailyWeights[j] * assetReturnValuesToday[j]
		}
		dailyRebalancedReturns = append(dailyRebalancedReturns, dailyRebalancedRet)

		// calculate drift
		lnp := historicReturns.NumberOfRows()
		for j := 0; lnp > 0 && j < assetReturns.NumberOfColumns(); j++ {
			updatedWeights[j] *= 1 + assetReturnValuesToday[j]      // drift
			updatedDailyWeights[j] *= 1 + assetReturnValuesToday[j] // drift
		}

		if shouldRebalanceAssetWeights(today, updatedDailyWeights) {
			copy(updatedWeights, policyWeights)
			rebalanceCount++
			result.RebalanceTimes = append(result.RebalanceTimes, today)
		}
		copy(updatedDailyWeights, policyWeights)
	}

	dailyRebalancedReturns.Reverse()
	backTestReturns.Reverse()
	reverseInPlace(result.Weights)
	result.Weights = slices.Clip(result.Weights)
	reverseInPlace(result.RebalanceTimes)
	result.RebalanceTimes = slices.Clip(result.RebalanceTimes)
	reverseInPlace(result.PolicyUpdateTimes)
	result.PolicyUpdateTimes = slices.Clip(result.PolicyUpdateTimes)
	result.ReturnsTable = returns.NewTable([]returns.List{
		backTestReturns,
		dailyRebalancedReturns,
	})

	return result, nil
}

func reverseInPlace[E any](s []E) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func fetchPolicy(ctx context.Context, end, start time.Time, alg PolicyWeightCalculator, assetReturns returns.Table, window WindowFunc) (time.Time, []float64, error) {
	ws := make([]float64, assetReturns.NumberOfColumns())

	var historicReturns returns.Table

	var (
		next    time.Time
		hasNext = true
	)

	for today := start; hasNext; today = next {
		next, hasNext = assetReturns.TimeAfter(today)

		if today.Before(start) {
			continue
		}
		if today.After(end) {
			break
		}
		historicReturns = window(today, assetReturns)

		setFloat64Slice(ws, 0)

		policyWeights, err := alg.PolicyWeights(ctx, today, historicReturns, ws)
		if err != nil {
			if errors.Is(err, ErrorNotEnoughData{}) {
				continue
			}
			return time.Time{}, nil, err
		}

		if len(policyWeights) != assetReturns.NumberOfColumns() {
			return time.Time{}, nil, fmt.Errorf("expected policy to have %d weights but got %d", assetReturns.NumberOfColumns(), len(policyWeights))
		}

		scaleToUnitRange(policyWeights)

		return today, policyWeights, nil
	}

	return time.Time{}, nil, ErrorNotEnoughData{}
}

func scaleToUnitRange(list []float64) {
	sum := 0.0
	for _, v := range list {
		sum += v
	}
	if sum == 0 {
		return
	}
	for i := range list {
		list[i] /= sum
	}
}

func setFloat64Slice(a []float64, v float64) {
	for i := range a {
		a[i] = v
	}
}
