package backtestconfig_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/portfolio/backtest"
	"github.com/portfoliotree/portfolio/backtest/backtestconfig"
)

var _ backtest.TriggerFunc = backtestconfig.IntervalDaily.CheckFunction()

func TestNamedInterval_CheckFunc(t *testing.T) {
	for _, interval := range backtestconfig.Intervals() {
		t.Run(interval.String(), func(t *testing.T) {
			fn := interval.CheckFunction()
			assert.NotNil(t, fn)
		})
	}

	t.Run("not set", func(t *testing.T) {
		fn := backtestconfig.Interval("").CheckFunction()
		assert.NotNil(t, fn)
	})

	t.Run("Cat", func(t *testing.T) {
		fn := backtestconfig.Interval("Cat").CheckFunction()
		assert.NotNil(t, fn)
	})
}

func TestDaily(t *testing.T) {
	fn := backtestconfig.Daily()

	tm := date("2020-03-01")
	final := date("2020-03-31")

	trueCount := 0
	for tm.Before(final) || tm.Equal(final) {
		if fn(tm, nil) {
			trueCount++
		}

		tm = tm.AddDate(0, 0, 1)
	}

	assert.Equal(t, trueCount, 31)
}

func TestWeekly(t *testing.T) {
	fn := backtestconfig.Weekly()

	tm := date("2020-03-03")
	final := date("2020-03-23")

	trueCount := 0
	for tm.Before(final) || tm.Equal(final) {
		if fn(tm, nil) {
			trueCount++
		}

		tm = tm.AddDate(0, 0, 1)
	}

	assert.Equal(t, trueCount, 4)
}

func TestMonthly(t *testing.T) {
	fn := backtestconfig.Monthly()

	tm := date("2020-03-03")
	final := date("2021-03-03")

	trueCount := 0
	for tm.Before(final) || tm.Equal(final) {
		if fn(tm, nil) {
			trueCount++
		}

		tm = tm.AddDate(0, 0, 1)
	}

	assert.Equal(t, trueCount, 13)
}

func TestQuarterly(t *testing.T) {
	fn := backtestconfig.Quarterly()

	tm := date("2018-01-02")
	final := date("2020-04-01")

	trueCount := 0
	for tm.Before(final) || tm.Equal(final) {
		if fn(tm, nil) {
			trueCount++
		}

		tm = tm.AddDate(0, 0, 1)
	}

	assert.Equal(t, trueCount, 10)
}

func TestAnnually(t *testing.T) {
	fn := backtestconfig.Annually()

	tm := date("2015-01-02")
	final := date("2020-04-01")

	trueCount := 0
	for tm.Before(final) || tm.Equal(final) {
		if fn(tm, nil) {
			trueCount++
		}

		tm = tm.AddDate(0, 0, 1)
	}

	assert.Equal(t, trueCount, 6)
}

func date(str string) time.Time {
	d, _ := time.Parse("2006-01-02", str)
	return d
}
