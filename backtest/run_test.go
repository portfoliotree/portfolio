package backtest_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/portfoliotree/round"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/portfoliotree/portfolio/allocation"
	"github.com/portfoliotree/portfolio/backtest"
	"github.com/portfoliotree/portfolio/backtest/backtestconfig"
	"github.com/portfoliotree/portfolio/returns"
)

func TestSpec_Run(t *testing.T) {
	t.Run("end date is before start date", func(t *testing.T) {
		assets := returns.NewTable([]returns.List{
			{{Time: date("2020-01-03")}, {Time: date("2020-01-02")}, {Time: date("2020-01-01")}},
			{{Time: date("2020-01-03")}, {Time: date("2020-01-02")} /*{Time: date("2020-01-01")}*/},
		})

		_, err := backtest.Run(context.Background(), date("2020-01-01"), date("2020-01-03"), assets, nil, nil, nil, nil)
		assert.Error(t, err)
	})
	t.Run("when zero value times are passed", func(t *testing.T) {
		assets := returns.NewTable([]returns.List{{
			{Time: date("2021-01-04"), Value: 0.8},
			{Time: date("2021-01-03"), Value: 0.4},
			{Time: date("2021-01-02"), Value: 0.2},
			{Time: date("2021-01-01"), Value: 0.1},
		}})
		alg := testAlgorithm()
		windowFunc := backtestconfig.WindowNotSet.Function
		rebalanceIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()
		policyUpdateIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()

		bt, err := backtest.Run(context.Background(), time.Time{}, time.Time{}, assets, alg, windowFunc, rebalanceIntervalFunc, policyUpdateIntervalFunc)
		assert.NoError(t, err)
		assert.Equal(t, date("2021-01-04"), bt.ReturnsTable.LastTime())
		assert.Equal(t, date("2021-01-01"), bt.ReturnsTable.FirstTime())
	})
	t.Run("when alg returns values they are not changed", func(t *testing.T) {
		assets := returns.NewTable([]returns.List{
			{
				{Time: date("2021-01-04"), Value: 0.8},
				{Time: date("2021-01-03"), Value: 0.4},
				{Time: date("2021-01-02"), Value: 0.2},
				{Time: date("2021-01-01"), Value: 0.1},
			},
			{
				{Time: date("2021-01-04"), Value: 0.8},
				{Time: date("2021-01-03"), Value: 0.4},
				{Time: date("2021-01-02"), Value: 0.2},
				{Time: date("2021-01-01"), Value: 0.1},
			},
		})
		windowFunc := backtestconfig.WindowNotSet.Function
		rebalanceIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()
		policyUpdateIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()

		ws := []float64{.715, .315}
		_, err := backtest.Run(context.Background(), date("2021-01-04"), date("2021-01-01"), assets, allocationFunction(func(_ context.Context, _ time.Time, _ returns.Table, currentWeights []float64) ([]float64, error) {
			return ws, nil
		}), windowFunc, rebalanceIntervalFunc, policyUpdateIntervalFunc)
		assert.NoError(t, err)
		assert.Equal(t, []float64{.715, .315}, ws)
	})
	t.Run("start date does not have a return", func(t *testing.T) {
		rs := returns.NewTable([]returns.List{
			{{Time: date("2020-01-03")}, {Time: date("2020-01-02")}, {Time: date("2020-01-01")}},
			{{Time: date("2020-01-03")}, {Time: date("2020-01-02")} /*{Time: date("2020-01-01")}*/},
		})

		_, err := backtest.Run(context.Background(), date("2020-01-01"), date("2020-01-03"), rs, nil, nil, nil, nil)
		assert.Error(t, err)
	})
	t.Run("end date does not have a return", func(t *testing.T) {
		alg := testAlgorithm()
		windowFunc := backtestconfig.WindowNotSet.Function
		rebalanceIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()
		policyUpdateIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()

		assets := returns.NewTable([]returns.List{
			{{Time: date("2020-01-03")}, {Time: date("2020-01-02")}, {Time: date("2020-01-01")}},
			{ /*{Time: date("2020-01-03")},*/ {Time: date("2020-01-02")}, {Time: date("2020-01-01")}},
		})
		start := assets.FirstTime()
		end := date("2020-01-03")

		_, err := backtest.Run(context.Background(), end, start, assets, alg, windowFunc, rebalanceIntervalFunc, policyUpdateIntervalFunc)
		assert.Error(t, err)
	})
	t.Run("with no returns", func(t *testing.T) {
		assets := returns.Table{}

		alg := testAlgorithm()
		windowFunc := backtestconfig.WindowNotSet.Function
		rebalanceIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()
		policyUpdateIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()

		end, start := time.Time{}, time.Time{}

		result, err := backtest.Run(context.Background(), end, start, assets, alg, windowFunc, rebalanceIntervalFunc, policyUpdateIntervalFunc)
		assert.Error(t, err)

		assert.Equal(t, result.ReturnsTable.NumberOfRows(), 0)
	})

	t.Run("when there is one asset", func(t *testing.T) {
		alg := testAlgorithm()
		windowFunc := backtestconfig.OneDayWindow.Function
		rebalanceIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()
		policyUpdateIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()

		asset := returns.List{
			{Time: date("2021-01-04"), Value: 0.8},
			{Time: date("2021-01-03"), Value: 0.4},
			{Time: date("2021-01-02"), Value: 0.2},
			{Time: date("2021-01-01"), Value: 0.1},
		}
		assets := returns.NewTable([]returns.List{asset})
		end, start, _ := assets.EndAndStartDates()

		result, err := backtest.Run(context.Background(), end, start, assets, alg, windowFunc, rebalanceIntervalFunc, policyUpdateIntervalFunc)

		assert.NoError(t, err)

		assert.Equal(t, result.ReturnsTable.NumberOfRows(), asset.Returns().Len())
	})

	t.Run("when called repeatedly", func(t *testing.T) {
		alg := testAlgorithm()
		windowFunc := backtestconfig.OneDayWindow.Function
		rebalanceIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()
		policyUpdateIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()

		asset := returns.List{
			{Time: date("2021-01-07"), Value: 6.4},
			{Time: date("2021-01-06"), Value: 3.2},
			{Time: date("2021-01-05"), Value: 1.6},
			{Time: date("2021-01-04"), Value: 0.8},
			{Time: date("2021-01-03"), Value: 0.4},
			{Time: date("2021-01-02"), Value: 0.2},
			{Time: date("2021-01-01"), Value: 0.1},
		}
		assets := returns.NewTable([]returns.List{asset})
		end, start, _ := assets.EndAndStartDates()

		result, err := backtest.Run(context.Background(), end, start, assets, alg, windowFunc, rebalanceIntervalFunc, policyUpdateIntervalFunc)
		require.NoError(t, err)
		require.Equal(t, result.ReturnsTable.FirstTime(), start)
		require.Equal(t, result.ReturnsTable.LastTime(), end)
		for i := 0; i < 100; i++ {
			result, err = backtest.Run(context.Background(), result.ReturnsTable.LastTime(), result.ReturnsTable.FirstTime(), assets, alg, windowFunc, rebalanceIntervalFunc, policyUpdateIntervalFunc)
			require.NoError(t, err)
		}
		assert.Equal(t, result.ReturnsTable.FirstTime(), start)
		assert.Equal(t, result.ReturnsTable.LastTime(), end)
	})

	t.Run("it responds to context cancellation", func(t *testing.T) {
		asset := returns.List{
			{Time: date("2021-01-04"), Value: 0.8},
			{Time: date("2021-01-03"), Value: 0.4},
			{Time: date("2021-01-02"), Value: 0.2},
			{Time: date("2021-01-01"), Value: 0.1},
		}
		assets := returns.NewTable([]returns.List{asset})
		end, start := asset[0].Time, asset[len(asset)-1].Time

		ctx, cancel := context.WithCancel(context.Background())
		c := make(chan struct{})
		var err error
		go func() {
			<-c
			cancel()
		}()
		alg := allocationFunction(func(ctx context.Context, _ time.Time, _ returns.Table, ws []float64) (targetWeights []float64, err error) {
			close(c)
			<-ctx.Done()
			return ws, ctx.Err()
		})
		windowFunc := backtestconfig.OneDayWindow.Function
		rebalanceIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()
		policyUpdateIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()

		_, err = backtest.Run(ctx, end, start, assets, alg, windowFunc, rebalanceIntervalFunc, policyUpdateIntervalFunc)

		assert.Equal(t, err, context.Canceled)
	})

	t.Run("daily rebalancing", func(t *testing.T) {
		asset1 := returns.List{
			{Time: date("2021-01-07"), Value: -0.1},
			{Time: date("2021-01-06"), Value: 0.25},
			{Time: date("2021-01-05"), Value: -0.1},
			{Time: date("2021-01-04"), Value: 0.3},
			{Time: date("2021-01-03"), Value: 0.05},
			{Time: date("2021-01-02"), Value: 0.1},
			{Time: date("2021-01-01"), Value: 0},
		}
		asset2 := returns.List{
			{Time: date("2021-01-07"), Value: 0.3},
			{Time: date("2021-01-06"), Value: 0.2},
			{Time: date("2021-01-05"), Value: -0.5},
			{Time: date("2021-01-04"), Value: -0.1},
			{Time: date("2021-01-03"), Value: 0.2},
			{Time: date("2021-01-02"), Value: 0.5},
			{Time: date("2021-01-01"), Value: -0.5},
		}
		assets := returns.NewTable([]returns.List{asset1, asset2})

		alg := testAlgorithm()
		windowFunc := backtestconfig.OneDayWindow.Function
		rebalanceIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()
		policyUpdateIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()

		end, start, _ := assets.EndAndStartDates()
		result, err := backtest.Run(context.Background(), end, start, assets, alg, windowFunc, rebalanceIntervalFunc, policyUpdateIntervalFunc)

		assert.NoError(t, err)
		values := result.Returns().Values()
		_ = round.Recursive(values, 3)
		assert.Equal(t, values, []float64{0.1, 0.225, -0.3, 0.1, 0.125, 0.3, -0.25})
	})

	t.Run("when the policy is not implementable at first data", func(t *testing.T) {
		asset1 := returns.List{
			{Time: date("2021-04-23"), Value: -0.1},
			{Time: date("2021-04-22"), Value: 0.25},
			{Time: date("2021-04-21"), Value: -0.1},
			{Time: date("2021-04-20"), Value: 0.30},
			{Time: date("2021-04-19"), Value: 0.05},
			{Time: date("2021-04-16"), Value: 0.1},
			{Time: date("2021-04-15"), Value: 0},
		}
		asset2 := returns.List{
			{Time: date("2021-04-23"), Value: -0.1},
			{Time: date("2021-04-22"), Value: 0.25},
			{Time: date("2021-04-21"), Value: -0.1},
			{Time: date("2021-04-20"), Value: 0.30},
			{Time: date("2021-04-19"), Value: 0.05},
			{Time: date("2021-04-16"), Value: 0.1},
			{Time: date("2021-04-15"), Value: 0},
		}
		assets := returns.NewTable([]returns.List{asset1, asset2})
		fallback := testAlgorithm()
		alg := allocationFunction(func(ctx context.Context, t time.Time, assetReturns returns.Table, currentWeights []float64) ([]float64, error) {
			if t.Before(date("2021-04-20")) {
				return nil, backtest.ErrorNotEnoughData{}
			}
			return fallback.PolicyWeights(ctx, t, assetReturns, currentWeights)
		})
		windowFunc := backtestconfig.WindowNotSet.Function
		rebalanceIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()
		policyUpdateIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()

		end, start, _ := assets.EndAndStartDates()
		result, err := backtest.Run(context.Background(), end, start, assets, alg, windowFunc, rebalanceIntervalFunc, policyUpdateIntervalFunc)

		assert.NoError(t, err)

		expected := returns.List{
			{Time: date("2021-04-23"), Value: -0.10},
			{Time: date("2021-04-22"), Value: 0.25},
			{Time: date("2021-04-21"), Value: -0.10},
			{Time: date("2021-04-20"), Value: 0.30},
		}

		rs := result.Returns()
		_ = round.Recursive(rs, 2)
		assert.Equal(t, rs, expected)
	})

	t.Run("composite returns are calculated correctly", func(t *testing.T) {
		asset1 := returns.List{
			{Time: date("2021-04-04"), Value: 0.20},
			{Time: date("2021-04-03"), Value: 0.10},
			{Time: date("2021-04-02"), Value: 0.00},
			{Time: date("2021-04-01"), Value: 0.50},
		}
		asset2 := returns.List{
			{Time: date("2021-04-04"), Value: 0.00},
			{Time: date("2021-04-03"), Value: 0.10},
			{Time: date("2021-04-02"), Value: 0.20},
			{Time: date("2021-04-01"), Value: -0.30},
		}
		assets := returns.NewTable([]returns.List{asset1, asset2})

		alg := testAlgorithm()
		windowFunc := backtestconfig.OneWeekWindow.Function
		rebalanceIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()
		policyUpdateIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()

		end, start, _ := assets.EndAndStartDates()
		result, err := backtest.Run(context.Background(), end, start, assets, alg, windowFunc, rebalanceIntervalFunc, policyUpdateIntervalFunc)

		assert.NoError(t, err)

		expected := returns.List{
			{Time: date("2021-04-04"), Value: 0.10},
			{Time: date("2021-04-03"), Value: 0.10},
			{Time: date("2021-04-02"), Value: 0.10},
			{Time: date("2021-04-01"), Value: 0.10},
		}

		rs := result.Returns()
		_ = round.Recursive(rs, 2)
		assert.Equal(t, rs, expected)
	})

	t.Run("when a look back is set", func(t *testing.T) {
		asset1 := returns.List{
			{Time: date("2021-04-23"), Value: -0.1},
			{Time: date("2021-04-22"), Value: 0.25},
			{Time: date("2021-04-21"), Value: -0.1},
			{Time: date("2021-04-20"), Value: 0.3},
			{Time: date("2021-04-19"), Value: 0.05},
			{Time: date("2021-04-16"), Value: 0.1},
			{Time: date("2021-04-15"), Value: 0},
			{Time: date("2021-04-14"), Value: 0},
			{Time: date("2021-04-13"), Value: 0},
			{Time: date("2021-04-12"), Value: 0},
		}

		callCount := 0
		alg := allocationFunction(func(_ context.Context, tm time.Time, assetReturns returns.Table, currentWeights []float64) ([]float64, error) {
			callCount++
			assert.Equalf(t, assetReturns.NumberOfColumns(), 1, "call count %d", callCount)
			for c := 0; c < assetReturns.NumberOfColumns(); c++ {
				rs := assetReturns.List(c)
				switch callCount {
				case 1:
					assert.Equalf(t, rs.Times(), []time.Time{
						date("2021-04-19"),
						date("2021-04-16"),
						date("2021-04-15"),
						date("2021-04-14"),
						date("2021-04-13"),
					}, "call count %d", callCount)
				case 2:
					assert.Equalf(t, rs.Times(), []time.Time{
						date("2021-04-20"),
						date("2021-04-19"),
						date("2021-04-16"),
						date("2021-04-15"),
						date("2021-04-14"),
					}, "call count %d", callCount)
				case 7:
					assert.Equalf(t, rs.Times(), []time.Time{
						date("2021-04-23"),
						date("2021-04-22"),
						date("2021-04-21"),
						date("2021-04-20"),
						date("2021-04-19"),
					}, "call count %d", callCount)
				}
				assert.Lenf(t, rs, 5, "call count %d", callCount)
			}
			return (&allocation.EqualWeights{}).PolicyWeights(context.Background(), tm, assetReturns, currentWeights)
		})

		windowFunc := backtestconfig.OneWeekWindow.Function
		rebalanceIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()
		policyUpdateIntervalFunc := backtestconfig.IntervalWeekly.CheckFunction()

		assets := returns.NewTable([]returns.List{asset1})
		end, start, _ := assets.EndAndStartDates()
		start = backtestconfig.OneWeekWindow.Add(start)
		result, err := backtest.Run(context.Background(), end, start, assets,
			alg,
			windowFunc,
			rebalanceIntervalFunc,
			policyUpdateIntervalFunc,
		)
		assert.Equal(t, callCount, 5)
		assert.NoError(t, err)
		assert.Equal(t, result.ReturnsTable.NumberOfRows(), 5)
	})
}

func TestSpec_Run_weightHistory(t *testing.T) {
	t.Run("single asset", func(t *testing.T) {
		asset := returns.List{
			{Time: date("2021-01-22")},
			{Time: date("2021-01-21")},
			{Time: date("2021-01-20")},
			{Time: date("2021-01-19")},
			// {Time: date("2021-01-18")}, // MLK day

			{Time: date("2021-01-15")},
			{Time: date("2021-01-14")},
			{Time: date("2021-01-13")},
			{Time: date("2021-01-12")},
			{Time: date("2021-01-11")},

			{Time: date("2021-01-08")},
			{Time: date("2021-01-07")},
		}

		assets := returns.NewTable([]returns.List{asset})

		alg := allocationFunction(randomWeights)
		windowFunc := backtestconfig.WindowNotSet.Function
		rebalanceIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()
		policyUpdateIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()

		end, start, _ := assets.EndAndStartDates()

		result, err := backtest.Run(context.Background(), end, start, assets, alg, windowFunc, policyUpdateIntervalFunc, rebalanceIntervalFunc)

		assert.NoError(t, err)

		assert.Equal(t, [][]float64{
			{1},
			{1},
			{1},
			{1},
			{1},
			{1},
			{1},
			{1},
			{1},
			{1},
			{1},
		}, result.Weights)
		assert.Equal(t, result.PolicyUpdateTimes, []time.Time{
			date("2021-01-22"),
			date("2021-01-21"),
			date("2021-01-20"),
			date("2021-01-19"),
			date("2021-01-15"),
			date("2021-01-14"),
			date("2021-01-13"),
			date("2021-01-12"),
			date("2021-01-11"),
			date("2021-01-08"),
		})

		assert.Equal(t, result.RebalanceTimes, []time.Time{
			date("2021-01-22"),
			date("2021-01-21"),
			date("2021-01-20"),
			date("2021-01-19"),
			date("2021-01-15"),
			date("2021-01-14"),
			date("2021-01-13"),
			date("2021-01-12"),
			date("2021-01-11"),
			date("2021-01-08"),
			date("2021-01-07"),
		})
	})

	t.Run("two assets with weekly rebalancing", func(t *testing.T) {
		asset1 := returns.List{
			{Time: date("2021-01-22"), Value: 0.0},
			{Time: date("2021-01-21"), Value: 0.0},
			{Time: date("2021-01-20"), Value: 0.0},
			{Time: date("2021-01-19"), Value: 0.0},
			// MLK day
			// {Time: date("2021-01-18"), Value: 0.0},

			{Time: date("2021-01-15"), Value: 0.0},
			{Time: date("2021-01-14"), Value: 0.0},
			{Time: date("2021-01-13"), Value: 0.0},
			{Time: date("2021-01-12"), Value: 0.0},
			{Time: date("2021-01-11"), Value: 0.0},

			{Time: date("2021-01-08"), Value: 0.0},
			{Time: date("2021-01-07"), Value: 0.0},
		}

		asset2 := returns.List{
			{Time: date("2021-01-22"), Value: 0.1},
			{Time: date("2021-01-21"), Value: 0.1},
			{Time: date("2021-01-20"), Value: 0.1},
			{Time: date("2021-01-19"), Value: 0.1},
			// MLK day
			// {Time: date("2021-01-18"), Value: 0.1},

			{Time: date("2021-01-15"), Value: 0.1},
			{Time: date("2021-01-14"), Value: 0.1},
			{Time: date("2021-01-13"), Value: 0.1},
			{Time: date("2021-01-12"), Value: 0.1},
			{Time: date("2021-01-11"), Value: 0.1},

			{Time: date("2021-01-08"), Value: 0.1},
			{Time: date("2021-01-07"), Value: 0.1},
		}
		assets := returns.NewTable([]returns.List{asset1, asset2})

		alg := testAlgorithm()
		windowFunc := backtestconfig.WindowNotSet.Function
		rebalanceIntervalFunc := backtestconfig.IntervalWeekly.CheckFunction()
		policyUpdateIntervalFunc := backtestconfig.IntervalMonthly.CheckFunction()

		end, start, _ := assets.EndAndStartDates()
		result, err := backtest.Run(context.Background(), end, start, assets, alg, windowFunc, policyUpdateIntervalFunc, rebalanceIntervalFunc)

		assert.NoError(t, err)

		assert.Len(t, result.Weights, len(asset1))
		assert.Len(t, result.RebalanceTimes, 3)
		assert.Len(t, result.PolicyUpdateTimes, 0, "calculating the initial policy weights is not an update")
	})

	t.Run("daily rebalanced returns is the same when daily rebalancing", func(t *testing.T) {
		asset1 := returns.List{
			{Time: date("2021-01-22"), Value: 0.1},
			{Time: date("2021-01-21"), Value: 0.1},
			{Time: date("2021-01-20"), Value: 0.2},
			{Time: date("2021-01-19"), Value: 0.1},
			// MLK day
			// {Time: date("2021-01-18"), Value: 0.1},

			{Time: date("2021-01-15"), Value: 0.1},
			{Time: date("2021-01-14"), Value: 0.2},
			{Time: date("2021-01-13"), Value: 0.1},
			{Time: date("2021-01-12"), Value: 0.3},
			{Time: date("2021-01-11"), Value: 0.1},

			{Time: date("2021-01-08"), Value: 0.1},
			{Time: date("2021-01-07"), Value: 0.1},
			{Time: date("2021-01-06"), Value: 0.1},
			{Time: date("2021-01-05"), Value: 0.3},
		}
		asset2 := returns.List{
			{Time: date("2021-01-22"), Value: 0.2},
			{Time: date("2021-01-21"), Value: 0.1},
			{Time: date("2021-01-20"), Value: 0.1},
			{Time: date("2021-01-19"), Value: 0.2},
			// MLK day
			// {Time: date("2021-01-18"), Value: 0.1},

			{Time: date("2021-01-15"), Value: 0.2},
			{Time: date("2021-01-14"), Value: -0.1},
			{Time: date("2021-01-13"), Value: 0.6},
			{Time: date("2021-01-12"), Value: 0.1},
			{Time: date("2021-01-11"), Value: 0.1},

			{Time: date("2021-01-08"), Value: 0.1},
			{Time: date("2021-01-07"), Value: 0.1},
			{Time: date("2021-01-06"), Value: 0.3},
			{Time: date("2021-01-05"), Value: 0.1},
		}
		assets := returns.NewTable([]returns.List{asset1, asset2})

		alg := testAlgorithm()
		windowFunc := backtestconfig.WindowNotSet.Function
		rebalanceIntervalFunc := backtestconfig.IntervalDaily.CheckFunction()
		policyUpdateIntervalFunc := backtestconfig.IntervalWeekly.CheckFunction()

		end, start, _ := assets.EndAndStartDates()
		result, err := backtest.Run(context.Background(), end, start, assets, alg, windowFunc, policyUpdateIntervalFunc, rebalanceIntervalFunc)

		assert.NoError(t, err)

		assert.Equal(t, result.Returns(), result.DailyRebalancedReturns())
	})
}

func randomWeights(_ context.Context, _ time.Time, _ returns.Table, currentWeights []float64) (targetWeights []float64, err error) {
	for i := range currentWeights {
		currentWeights[i] = rand.Float64()
	}
	return currentWeights, nil
}

func date(str string) time.Time {
	d, _ := time.Parse(time.DateOnly, str)
	return d
}

type allocationFunction func(_ context.Context, _ time.Time, _ returns.Table, currentWeights []float64) (targetWeights []float64, err error)

func (function allocationFunction) PolicyWeights(ctx context.Context, today time.Time, assets returns.Table, ws []float64) (targetWeights []float64, err error) {
	return function(ctx, today, assets, ws)
}

func testAlgorithm() allocation.Algorithm {
	return new(allocation.EqualWeights)
}
