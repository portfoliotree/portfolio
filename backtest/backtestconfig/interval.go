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
