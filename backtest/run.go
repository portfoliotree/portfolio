package backtest

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/slices"
	"gonum.org/v1/gonum/floats"

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
// when a policy update or rebalancing is required.
func Run(ctx context.Context, end, start time.Time, assetReturns returns.Table,
	alg PolicyWeightCalculator,
	lookBackWindow WindowFunc,
	shouldCalculatePolicy,
	shouldRebalanceAssetWeights TriggerFunc,
) (Result, error) {
	if assetReturns.NumberOfColumns() == 0 {
		return Result{}, errors.New("no asset returns provided")
	}

	end, start, err := ensureDatesAreWithinAssetTableRange(end, start, assetReturns)
	if err != nil {
		return Result{}, err
	}

	firstPolicyDate, policyWeights, err := fetchPolicy(ctx, end, start, alg, assetReturns, lookBackWindow)
	if err != nil {
		return Result{}, err
	}

	start = firstPolicyDate

	var (
		updatedWeights      = slices.Clone(policyWeights)
		updatedDailyWeights = slices.Clone(policyWeights)

		backTestReturns,
		dailyRebalancedReturns returns.List
		historicReturns        returns.Table
		assetReturnValuesToday = make([]float64, assetReturns.NumberOfColumns())

		rebalanceCount, recalculatePolicyCount = 0, 0

		next    time.Time
		hasNext = true

		currentWeightsPolicyWeightsInput = make([]float64, assetReturns.NumberOfColumns())

		result = Result{
			Weights:            make([][]float64, 0, assetReturns.NumberOfRows()),
			RebalanceTimes:     make([]time.Time, 0, assetReturns.NumberOfRows()),
			PolicyUpdateTimes:  make([]time.Time, 0, assetReturns.NumberOfRows()),
			FinalPolicyWeights: make([]float64, assetReturns.NumberOfColumns()),
		}
	)

	for today, i := start, 0; hasNext && !today.After(end) && i < assetReturns.NumberOfRows(); today, i = next, i+1 {
		next, hasNext = assetReturns.TimeAfter(today)
		scaleToUnitRange(updatedWeights)
		scaleToUnitRange(updatedDailyWeights)

		historicReturns = lookBackWindow(today, assetReturns)
		assetReturnValuesToday = mostRecentValues(assetReturnValuesToday, historicReturns)

		if shouldCalculatePolicy(today, updatedDailyWeights) && start != today {
			copy(currentWeightsPolicyWeightsInput, updatedWeights)
			pw, err := alg.PolicyWeights(ctx, today, historicReturns, currentWeightsPolicyWeightsInput)
			if err != nil {
				return Result{}, err
			}
			copy(policyWeights, pw)

			scaleToUnitRange(policyWeights)

			recalculatePolicyCount++
			result.PolicyUpdateTimes = append(result.PolicyUpdateTimes, today)
			copy(result.FinalPolicyWeights, policyWeights)
		}

		backTestReturns = append(backTestReturns, returns.Return{
			Time:  today,
			Value: floats.Dot(updatedWeights, assetReturnValuesToday),
		})
		dailyRebalancedReturns = append(dailyRebalancedReturns, returns.Return{
			Time:  today,
			Value: floats.Dot(updatedDailyWeights, assetReturnValuesToday),
		})

		// calculate drift
		for j := 0; historicReturns.NumberOfRows() > 0 && j < assetReturns.NumberOfColumns(); j++ {
			updatedWeights[j] *= 1 + assetReturnValuesToday[j]      // drift
			updatedDailyWeights[j] *= 1 + assetReturnValuesToday[j] // drift
		}

		if shouldRebalanceAssetWeights(today, updatedDailyWeights) {
			copy(updatedWeights, policyWeights)
			rebalanceCount++
			result.RebalanceTimes = append(result.RebalanceTimes, today)
		}

		result.Weights = append(result.Weights, slices.Clone(updatedWeights))

		copy(updatedDailyWeights, policyWeights)
	}

	slices.Reverse(backTestReturns)
	slices.Reverse(dailyRebalancedReturns)
	slices.Reverse(result.Weights)
	result.Weights = slices.Clip(result.Weights)
	slices.Reverse(result.RebalanceTimes)
	result.RebalanceTimes = slices.Clip(result.RebalanceTimes)
	slices.Reverse(result.PolicyUpdateTimes)
	result.PolicyUpdateTimes = slices.Clip(result.PolicyUpdateTimes)
	result.ReturnsTable = returns.NewTable([]returns.List{
		backTestReturns,
		dailyRebalancedReturns,
	})

	return result, nil
}

func mostRecentValues(mostRecentValues []float64, table returns.Table) []float64 {
	values := table.ColumnValues()
	for i := range values {
		mostRecentValues[i] = values[i][0]
	}
	return mostRecentValues
}

func ensureDatesAreWithinAssetTableRange(end, start time.Time, assetReturns returns.Table) (time.Time, time.Time, error) {
	if start.IsZero() {
		start = assetReturns.FirstTime()
	}
	if end.IsZero() {
		end = assetReturns.LastTime()
	}
	if end.After(assetReturns.LastTime()) || start.Before(assetReturns.FirstTime()) {
		return time.Time{}, time.Time{}, ErrorNotEnoughData{}
	}
	return end, start, nil
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

		pw, err := alg.PolicyWeights(ctx, today, historicReturns, ws)
		if err != nil {
			if errors.Is(err, ErrorNotEnoughData{}) {
				continue
			}
			return time.Time{}, nil, err
		}
		policyWeights := slices.Clone(pw)

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
