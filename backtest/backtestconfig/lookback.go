package backtestconfig

import (
	"time"

	"github.com/portfoliotree/portfolio/returns"
)

func TruncateWithLookBack(end, start time.Time, dur Window, returns returns.List) returns.List {
	for i, r := range returns {
		if !r.Time.After(dur.Sub(start)) || i == len(returns)-1 {
			start = r.Time
			break
		}
	}
	return returns.Between(end, start)
}
