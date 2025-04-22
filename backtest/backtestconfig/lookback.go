package backtestconfig

import (
	"time"

	"github.com/portfoliotree/timetable"
)

func TruncateWithLookBack(end, start time.Time, dur Window, returns timetable.List[float64]) timetable.List[float64] {
	for i, r := range returns {
		if tm := r.Time(); !tm.After(dur.Sub(start)) || i == len(returns)-1 {
			start = tm
			break
		}
	}
	return returns.Between(end, start)
}
