package backtestconfig_test

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/portfoliotree/portfolio/backtest"
	"github.com/portfoliotree/portfolio/backtest/backtestconfig"
)

var _ backtest.TriggerFunc = backtestconfig.IntervalDaily.CheckFunction()

func TestNamedInterval_CheckFunc(t *testing.T) {
	for _, interval := range backtestconfig.Intervals() {
		t.Run(interval.String(), func(t *testing.T) {
			o := NewWithT(t)
			fn := interval.CheckFunction()
			o.Expect(fn).NotTo(BeNil())
		})
	}

	t.Run("not set", func(t *testing.T) {
		o := NewWithT(t)
		fn := backtestconfig.Interval("").CheckFunction()
		o.Expect(fn).NotTo(BeNil())
	})

	t.Run("Cat", func(t *testing.T) {
		o := NewWithT(t)
		fn := backtestconfig.Interval("Cat").CheckFunction()
		o.Expect(fn).NotTo(BeNil())
	})
}

func TestDaily(t *testing.T) {
	please := NewWithT(t)

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

	please.Expect(trueCount).To(Equal(31))
}

func TestWeekly(t *testing.T) {
	please := NewWithT(t)

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

	please.Expect(trueCount).To(Equal(4))
}

func TestMonthly(t *testing.T) {
	please := NewWithT(t)

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

	please.Expect(trueCount).To(Equal(13))
}

func TestQuarterly(t *testing.T) {
	please := NewWithT(t)

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

	please.Expect(trueCount).To(Equal(10))
}

func TestAnnually(t *testing.T) {
	please := NewWithT(t)

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

	please.Expect(trueCount).To(Equal(6))
}

func date(str string) time.Time {
	d, _ := time.Parse("2006-01-02", str)
	return d
}
