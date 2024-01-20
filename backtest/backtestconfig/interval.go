package backtestconfig

import (
	"fmt"
	"time"
)

type (
	Interval string

	IntervalCheckFunc func(t time.Time) bool
)

const (
	IntervalDefault            = IntervalNever
	IntervalNever     Interval = "Never"
	IntervalDaily     Interval = "Daily"
	IntervalWeekly    Interval = "Weekly"
	IntervalMonthly   Interval = "Monthly"
	IntervalQuarterly Interval = "Quarterly"
	IntervalAnnually  Interval = "Annually"
)

func (t Interval) Options() []Interval {
	return Intervals()
}

func Intervals() []Interval {
	return []Interval{
		IntervalNever,
		IntervalDaily,
		IntervalWeekly,
		IntervalMonthly,
		IntervalQuarterly,
		IntervalAnnually,
	}
}

func (t Interval) CheckFunction() func(t time.Time, currentWeights []float64) bool {
	switch t {
	case IntervalDaily:
		return Daily()
	case IntervalWeekly:
		return Weekly()
	case IntervalMonthly:
		return Monthly()
	case IntervalQuarterly:
		return Quarterly()
	case IntervalAnnually:
		return Annually()
	case IntervalNever, "":
		fallthrough
	default:
		return Never()
	}
}

func (t Interval) String() string { return string(t) }

func (t Interval) Validate() error {
	switch t {
	case IntervalNever,
		IntervalDaily,
		IntervalWeekly,
		IntervalMonthly,
		IntervalQuarterly,
		IntervalAnnually,
		"":
		return nil
	default:
		return fmt.Errorf("unknown trigger interval %q", t)
	}
}

func Never() func(t time.Time, currentWeights []float64) bool {
	return func(_ time.Time, _ []float64) bool {
		return false
	}
}

func Daily() func(t time.Time, _ []float64) bool {
	return func(current time.Time, _ []float64) bool {
		return true
	}
}

func Weekly() func(t time.Time, _ []float64) bool {
	var previous time.Time
	return func(current time.Time, _ []float64) bool {
		isStartOfPeriod := previous.IsZero() || current.Weekday() < previous.Weekday()
		previous = current
		return isStartOfPeriod
	}
}

func Monthly() func(t time.Time, currentWeights []float64) bool {
	var previous time.Time
	return func(current time.Time, _ []float64) bool {
		isStartOfPeriod := previous.IsZero() || current.Day() < previous.Day()
		previous = current
		return isStartOfPeriod
	}
}

func isFirstMonthInQuarter(month time.Month) bool {
	return (month-1)%3 == 0
}

func Quarterly() func(t time.Time, _ []float64) bool {
	monthly := Monthly()
	return func(current time.Time, _ []float64) bool {
		return monthly(current, nil) && isFirstMonthInQuarter(current.Month())
	}
}

func Annually() func(t time.Time, _ []float64) bool {
	monthly := Monthly()
	return func(current time.Time, _ []float64) bool {
		return monthly(current, nil) && current.Month() == time.January
	}
}

func (t Interval) StartDate(now time.Time) time.Time {
	switch t {
	case IntervalWeekly:
		wd := now.Weekday()
		return now.AddDate(0, 0, -int(wd)+1)
	case IntervalMonthly:
		y, m, _ := now.Date()
		return time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
	case IntervalQuarterly:
		y, m, _ := now.Date()
		m = quarterMonth(m)
		return time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
	case IntervalAnnually:
		return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	case IntervalNever, IntervalDaily, "":
		fallthrough
	default:
		return now
	}
}

func quarterMonth(m time.Month) time.Month {
	switch m {
	case time.January, time.February, time.March:
		fallthrough
	default:
		return time.January
	case time.April, time.May, time.June:
		return time.April
	case time.July, time.August, time.September:
		return time.July
	case time.October, time.November, time.December:
		return time.October
	}
}
