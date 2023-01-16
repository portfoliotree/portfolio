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

func TruncateTableWithLookBack(end, start time.Time, dur Window, table returns.Table) returns.Table {
	result := make([]returns.List, table.NumberOfColumns())
	for i, l := range table.Lists() {
		result[i] = TruncateWithLookBack(end, start, dur, l)
	}
	return returns.NewTable(result)
}
