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
	d, _ := time.Parse(time.DateOnly, str)
	return d
}

func TestInterval_StartDate(t *testing.T) {
	for _, tt := range []struct {
		Name     string
		Interval backtestconfig.Interval
		Expected time.Time
		Now      time.Time
	}{
		{
			Name:     "never gives the same date",
			Interval: backtestconfig.IntervalNever,
			Expected: date("2020-03-01"),
			Now:      date("2020-03-01"),
		},
		{
			Name:     "daily gives the same day",
			Interval: backtestconfig.IntervalDaily,
			Expected: date("2020-03-01"),
			Now:      date("2020-03-01"),
		},
		{
			Name:     "weekly gives the most recent monday",
			Interval: backtestconfig.IntervalWeekly,
			Expected: date("2024-01-01"),
			Now:      date("2024-01-03"),
		},
		{
			Name:     "monthly gives the 1st of a month",
			Interval: backtestconfig.IntervalMonthly,
			Expected: date("2024-02-01"),
			Now:      date("2024-02-07"),
		},
		{
			Name:     "q1 gives the most recent 1st quarter day",
			Interval: backtestconfig.IntervalQuarterly,
			Expected: date("2024-01-01"),
			Now:      date("2024-02-23"),
		},
		{
			Name:     "q2 gives the most recent 1st quarter day",
			Interval: backtestconfig.IntervalQuarterly,
			Expected: date("2024-04-01"),
			Now:      date("2024-05-23"),
		},
		{
			Name:     "q3 gives the most recent 1st quarter day",
			Interval: backtestconfig.IntervalQuarterly,
			Expected: date("2024-07-01"),
			Now:      date("2024-09-23"),
		},
		{
			Name:     "q4 gives the most recent 1st quarter day",
			Interval: backtestconfig.IntervalQuarterly,
			Expected: date("2024-10-01"),
			Now:      date("2024-12-03"),
		},
		{
			Name:     "weekly gives the most recent monday",
			Interval: backtestconfig.IntervalAnnually,
			Expected: date("2024-01-01"),
			Now:      date("2024-04-23"),
		},
		{
			Name:     "weekly on monday gives the same day",
			Interval: backtestconfig.IntervalWeekly,
			Expected: date("2024-01-01"),
			Now:      date("2024-01-01"),
		},
		{
			Name:     "monthly on the first gives the same day",
			Interval: backtestconfig.IntervalMonthly,
			Expected: date("2024-01-01"),
			Now:      date("2024-01-01"),
		},
		{
			Name:     "monthly on the first gives the same day",
			Interval: backtestconfig.IntervalQuarterly,
			Expected: date("2024-01-01"),
			Now:      date("2024-01-01"),
		},
		{
			Name:     "annually on jan 1st gives the same day",
			Interval: backtestconfig.IntervalAnnually,
			Expected: date("2024-01-01"),
			Now:      date("2024-01-01"),
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			out := tt.Interval.StartDate(tt.Now)
			assert.Equal(t, tt.Expected, out)
		})
	}
}
