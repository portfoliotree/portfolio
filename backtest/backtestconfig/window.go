package backtestconfig

import (
	"fmt"
	"time"

	"github.com/portfoliotree/portfolio/returns"
)

type Window string

const (
	WindowNotSet     Window = ""
	OneDayWindow     Window = "1 Day"
	OneWeekWindow    Window = "1 Week"
	OneMonthWindow   Window = "1 Month"
	OneQuarterWindow Window = "1 Quarter"
	OneYearWindow    Window = "1 Year"
	ThreeYearWindow  Window = "3 Years"
	FiveYearWindow   Window = "5 Years"
)

func (dur Window) String() string { return string(dur) }

func (dur Window) Options() []Window {
	return Windows()
}

func Windows() []Window {
	return []Window{
		OneDayWindow,
		OneWeekWindow,
		OneMonthWindow,
		OneQuarterWindow,
		OneYearWindow,
		ThreeYearWindow,
		FiveYearWindow,
	}
}

func (dur Window) IsSet() bool { return dur != "" }

func (dur Window) Validate() error {
	switch dur {
	case WindowNotSet,
		OneDayWindow,
		OneWeekWindow,
		OneMonthWindow,
		OneQuarterWindow,
		OneYearWindow,
		ThreeYearWindow,
		FiveYearWindow:
		return nil
	default:
		return fmt.Errorf("unknown named duration %q", dur)
	}
}

func (dur Window) Add(t time.Time) time.Time {
	switch dur {
	case OneDayWindow:
		return t.AddDate(0, 0, 1)
	case OneWeekWindow:
		return t.AddDate(0, 0, 7)
	case OneMonthWindow:
		return t.AddDate(0, 1, 0)
	case OneQuarterWindow:
		return t.AddDate(0, 3, 0)
	case OneYearWindow:
		return t.AddDate(1, 0, 0)
	case ThreeYearWindow:
		return t.AddDate(3, 0, 0)
	case FiveYearWindow:
		return t.AddDate(5, 0, 0)
	}
	return t
}

func (dur Window) Sub(t time.Time) time.Time {
	switch dur {
	case OneDayWindow:
		return t.AddDate(0, 0, -1)
	case OneWeekWindow:
		return t.AddDate(0, 0, -6)
	case OneMonthWindow:
		return t.AddDate(0, -1, 1)
	case OneQuarterWindow:
		return t.AddDate(0, -3, 1)
	case OneYearWindow:
		return t.AddDate(-1, 0, 1)
	case ThreeYearWindow:
		return t.AddDate(-3, 0, 1)
	case FiveYearWindow:
		return t.AddDate(-5, 0, 1)
	}
	return t
}

func (dur Window) Function(today time.Time, table returns.Table) returns.Table {
	return table.Between(today, dur.Sub(today))
}
