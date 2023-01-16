package backtestconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/portfolio/backtest/backtestconfig"
	"github.com/portfoliotree/portfolio/returns"
)

func TestNamedDuration_TruncateWithLookBack(t *testing.T) {
	everyOtherWednesday := returns.List{
		{Time: date("2021-12-22")},
		{Time: date("2021-12-08")},
		{Time: date("2021-11-24")},
		{Time: date("2021-11-10")},
		{Time: date("2021-10-27")},
		{Time: date("2021-10-13")},
		{Time: date("2021-09-29")},
		{Time: date("2021-09-15")},
		{Time: date("2021-09-01")},
		{Time: date("2021-08-18")},
		{Time: date("2021-08-04")},
		{Time: date("2021-07-21")},
		{Time: date("2021-07-07")},
		{Time: date("2021-06-23")},
		{Time: date("2021-06-09")},
		{Time: date("2021-05-26")},
		{Time: date("2021-05-12")},
		{Time: date("2021-04-28")},
		{Time: date("2021-04-14")},
		{Time: date("2021-03-31")},
		{Time: date("2021-03-17")},
		{Time: date("2021-03-03")},
		{Time: date("2021-02-17")},
		{Time: date("2021-02-03")},
		{Time: date("2021-01-20")},
		{Time: date("2021-01-06")},

		{Time: date("2020-12-23")},
		{Time: date("2020-12-09")},
		{Time: date("2020-11-25")},
		{Time: date("2020-11-11")},
		{Time: date("2020-10-28")},
		{Time: date("2020-10-14")},
		{Time: date("2020-09-30")},
		{Time: date("2020-09-16")},
		{Time: date("2020-09-02")},
		{Time: date("2020-08-19")},
		{Time: date("2020-08-05")},
		{Time: date("2020-07-22")},
		{Time: date("2020-07-08")},
		{Time: date("2020-06-24")},
		{Time: date("2020-06-10")},
		{Time: date("2020-05-27")},
		{Time: date("2020-05-13")},
		{Time: date("2020-04-29")},
		{Time: date("2020-04-15")},
		{Time: date("2020-04-01")},
		{Time: date("2020-03-18")},
		{Time: date("2020-03-04")},
		{Time: date("2020-02-19")},
		{Time: date("2020-02-05")},
		{Time: date("2020-01-22")},
		{Time: date("2020-01-08")},
	}

	assert.Equal(t,
		backtestconfig.TruncateWithLookBack(
			date("2021-08-18"), date("2021-08-04"),
			backtestconfig.OneDayWindow, everyOtherWednesday,
		).Times(),
		everyOtherWednesday.Between(date("2021-08-18"), date("2021-07-21")).Times())

	assert.Equal(t,
		backtestconfig.TruncateWithLookBack(
			date("2020-12-30"), date("2020-12-23"),
			backtestconfig.OneMonthWindow, everyOtherWednesday,
		).Times(),
		everyOtherWednesday.Between(date("2020-12-30"), date("2020-11-11")).Times(),
	)

	assert.Equal(t,
		backtestconfig.TruncateWithLookBack(
			date("2021-09-04"), date("2021-08-04"),
			backtestconfig.OneQuarterWindow, everyOtherWednesday,
		).Times(),
		everyOtherWednesday.Between(date("2021-09-04"), date("2021-04-28")).Times(),
	)

	assert.Equal(t,
		backtestconfig.TruncateWithLookBack(
			date("2020-06-01"), date("2020-05-27"),
			backtestconfig.OneYearWindow, everyOtherWednesday,
		).Times(),
		everyOtherWednesday.Between(date("2020-06-01"), date("2020-01-08")).Times(),
	)
}
