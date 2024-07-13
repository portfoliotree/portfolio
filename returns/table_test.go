package returns_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/portfoliotree/round"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/portfoliotree/portfolio/internal/fixtures"
	"github.com/portfoliotree/portfolio/returns"
)

func TestNewTable(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		// given
		lists := []returns.List(nil)

		// when
		tab := returns.NewTable(lists)

		// then
		assert.Len(t, tab.Times(), 0)
		assert.Equal(t, tab.NumberOfColumns(), 0)
	})
	t.Run("empty lists", func(t *testing.T) {
		// given
		lists := []returns.List{
			{},
			{},
		}

		// when
		tab := returns.NewTable(lists)

		// then
		assert.Len(t, tab.Times(), 0)
		assert.Equal(t, tab.NumberOfColumns(), 2)
	})

	t.Run("a single return per list", func(t *testing.T) {
		// then
		t.Run("when times are fetched", func(t *testing.T) {
			lists := returns.NewTable([]returns.List{
				{rtn(t, fixtures.Day0, 0.1)},
				{rtn(t, fixtures.Day0, 0.2)},
				{rtn(t, fixtures.Day0, 0.3)},
			})
			assert.Len(t, lists.Times(), 1, "it returns one row")
		})

		t.Run("when the number of columns is fetched", func(t *testing.T) {
			lists := returns.NewTable([]returns.List{
				{rtn(t, fixtures.Day0, 0.1)},
				{rtn(t, fixtures.Day0, 0.2)},
				{rtn(t, fixtures.Day0, 0.3)},
			})
			assert.Equal(t, lists.NumberOfColumns(), 3, "it gives the correct column count")
		})

		t.Run("when the number of rows is fetched", func(t *testing.T) {
			lists := returns.NewTable([]returns.List{
				{rtn(t, fixtures.Day0, 0.1)},
				{rtn(t, fixtures.Day0, 0.2)},
				{rtn(t, fixtures.Day0, 0.3)},
			})
			assert.Equal(t, lists.NumberOfRows(), 1, "it gives the correct row count")
		})

		// then
		t.Run("when a row is fetched", func(t *testing.T) {
			lists := returns.NewTable([]returns.List{
				{rtn(t, fixtures.Day0, 0.1)},
				{rtn(t, fixtures.Day0, 0.2)},
				{rtn(t, fixtures.Day0, 0.3)},
			})
			row, found := lists.Row(fixtures.T(t, fixtures.Day0))
			assert.True(t, found)
			assert.Equal(t, row, []float64{0.1, 0.2, 0.3})
		})
	})

	t.Run("three aligned returns per list", func(t *testing.T) {
		lists := []returns.List{
			{rtn(t, fixtures.Day2, 0.1), rtn(t, fixtures.Day1, 0.01), rtn(t, fixtures.Day0, 0.001)},
			{rtn(t, fixtures.Day2, 0.2), rtn(t, fixtures.Day1, 0.02), rtn(t, fixtures.Day0, 0.002)},
			{rtn(t, fixtures.Day2, 0.3), rtn(t, fixtures.Day1, 0.03), rtn(t, fixtures.Day0, 0.003)},
		}

		// when
		tab := returns.NewTable(lists)

		// then
		t.Run("when times are fetched", func(t *testing.T) {
			assert.Equal(t, tab.Times(), []time.Time{fixtures.T(t, fixtures.Day2), fixtures.T(t, fixtures.Day1), fixtures.T(t, fixtures.Day0)}, "it has the right dates")
		})

		t.Run("when the number of columns is fetched", func(t *testing.T) {
			assert.Equal(t, tab.NumberOfColumns(), 3, "it gives the correct number")
		})

		// then
		t.Run("when a row is fetched", func(t *testing.T) {
			type testRow struct {
				name     string
				date     time.Time
				expFound bool
				expVal   []float64
			}

			for _, tt := range []testRow{
				{
					name:     "found initial",
					date:     fixtures.T(t, fixtures.Day0),
					expFound: true,
					expVal:   []float64{0.001, 0.002, 0.003},
				},
				{
					name:     "found final",
					date:     fixtures.T(t, fixtures.Day2),
					expFound: true,
					expVal:   []float64{0.1, 0.2, 0.3},
				},
				{
					name:     "found middle",
					date:     fixtures.T(t, fixtures.Day1),
					expFound: true,
					expVal:   []float64{0.01, 0.02, 0.03},
				},
				{
					name:     "missing before",
					date:     fixtures.T(t, fixtures.Day0).AddDate(-1, 0, 0),
					expFound: false,
					expVal:   []float64{0, 0, 0},
				},
				{
					name:     "missing after",
					date:     fixtures.T(t, fixtures.Day2).AddDate(1, 0, 0),
					expFound: false,
					expVal:   []float64{0, 0, 0},
				},
				{
					name:     "missing between",
					date:     fixtures.T(t, fixtures.Day1).AddDate(0, 1, 0),
					expFound: false,
					expVal:   []float64{0, 0, 0},
				},
			} {
				t.Run(tt.name, func(t *testing.T) {
					// when
					row, found := tab.Row(tt.date)

					// then
					assert.Equal(t, found, tt.expFound, "it finds the row")
					assert.Equal(t, row, tt.expVal, "it returns the correct values")
				})
			}
		})
	})

	t.Run("non aligned returns", func(t *testing.T) {
		t.Run("single returns with different dates", func(t *testing.T) {
			// given
			lists := returns.NewTable([]returns.List{
				{rtn(t, fixtures.Day1, 0.2)},
				{rtn(t, fixtures.Day0, 0.1)},
				{rtn(t, fixtures.Day2, 0.3)},
			})
			assert.Equal(t, lists.NumberOfRows(), 0)
			assert.Len(t, lists.Times(), 0)
		})

		t.Run("the second asset has more history", func(t *testing.T) {
			// given
			lists := returns.NewTable([]returns.List{
				{rtn(t, fixtures.Day1, 0.3)},
				{rtn(t, fixtures.Day1, 0.2), rtn(t, fixtures.Day0, 0.1)},
			})
			assert.Equal(t, lists.NumberOfRows(), 1)
			values, found := lists.Row(fixtures.T(t, fixtures.Day1))
			assert.True(t, found)
			assert.Equal(t, values, []float64{0.3, 0.2})
		})

		t.Run("the first asset has more history", func(t *testing.T) {
			// given
			lists := returns.NewTable([]returns.List{
				{rtn(t, fixtures.Day1, 0.2), rtn(t, fixtures.Day0, 0.1)},
				{rtn(t, fixtures.Day1, 0.3)},
			})
			assert.Equal(t, lists.NumberOfRows(), 1)
			values, found := lists.Row(fixtures.T(t, fixtures.Day1))
			assert.True(t, found)
			assert.Equal(t, values, []float64{0.2, 0.3})
		})

		t.Run("the first asset has more recent returns", func(t *testing.T) {
			// given
			lists := returns.NewTable([]returns.List{
				{rtn(t, fixtures.Day1, 0.2), rtn(t, fixtures.Day0, 0.1)},
				{ /*                       ,*/ rtn(t, fixtures.Day0, 0.3)},
			})
			assert.Equal(t, lists.NumberOfRows(), 1)

			values, found := lists.Row(fixtures.T(t, fixtures.Day0))
			assert.True(t, found)
			assert.Equal(t, values, []float64{0.1, 0.3})
		})

		t.Run("the second asset has more recent returns", func(t *testing.T) {
			// given
			lists := returns.NewTable([]returns.List{
				{ /*                       ,*/ rtn(t, fixtures.Day0, 0.3)},
				{rtn(t, fixtures.Day1, 0.2), rtn(t, fixtures.Day0, 0.1)},
			})
			assert.Equal(t, lists.NumberOfRows(), 1)
			values, found := lists.Row(fixtures.T(t, fixtures.Day0))
			assert.True(t, found)
			assert.Equal(t, values, []float64{0.3, 0.1})
		})
	})
}

func TestTable_Between(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		// given
		tab := returns.NewTable(nil)

		// when

		assert.NotPanics(t, func() {
			_ = tab.Between(fixtures.T(t, fixtures.Day0), fixtures.T(t, fixtures.Day1))
		})
	})

	t.Run("times are outside of table", func(t *testing.T) {
		// given
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.Day2, 0.10), rtn(t, fixtures.Day1, 0.20), rtn(t, fixtures.Day0, 0.30)},
			{rtn(t, fixtures.Day2, 0.01), rtn(t, fixtures.Day1, 0.02), rtn(t, fixtures.Day0, 0.03)},
		})

		queryEnd := fixtures.T(t, fixtures.DayAfter)
		queryStart := fixtures.T(t, fixtures.DayBefore)

		// when

		slice := table.Between(queryEnd, queryStart)

		// then

		assert.Equal(t, slice.ColumnValues(), table.ColumnValues())
		assert.Equal(t, slice.Times(), []time.Time{fixtures.T(t, fixtures.Day2), fixtures.T(t, fixtures.Day1), fixtures.T(t, fixtures.Day0)})
	})

	t.Run("times are inside the table", func(t *testing.T) {
		// given
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.Day3, -0.10), rtn(t, fixtures.Day2, 0.10), rtn(t, fixtures.Day1, 0.20), rtn(t, fixtures.Day0, 0.30)},
			{rtn(t, fixtures.Day3, -0.01), rtn(t, fixtures.Day2, 0.01), rtn(t, fixtures.Day1, 0.02), rtn(t, fixtures.Day0, 0.03)},
		})

		queryEnd, queryStart := fixtures.T(t, fixtures.Day2), fixtures.T(t, fixtures.Day1)

		// when

		slice := table.Between(queryEnd, queryStart)

		// then

		assert.Equal(t, slice.ColumnValues(), [][]float64{
			{0.10, 0.20},
			{0.01, 0.02},
		})
		assert.Equal(t, slice.Times(), []time.Time{fixtures.T(t, fixtures.Day2), fixtures.T(t, fixtures.Day1)})
	})
}

func TestTable_AddColumn(t *testing.T) {
	t.Run("when adding list with an additional row", func(t *testing.T) {
		t.Skip(`
AddColumn now does not add a column to the table if the table does not already have a row.
`)
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.Day3, .1), rtn(t, fixtures.Day1, .1), rtn(t, fixtures.Day0, .1)},
		})
		table = table.AddColumn(returns.List{
			rtn(t, fixtures.Day3, .1), rtn(t, fixtures.Day2, .1), rtn(t, fixtures.Day1, .1), rtn(t, fixtures.Day0, .1),
		})
		assert.Equal(t, table.Lists(), []returns.List{
			{rtn(t, fixtures.Day3, .1), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, .1), rtn(t, fixtures.Day0, .1)},
			{rtn(t, fixtures.Day3, .1), rtn(t, fixtures.Day2, .1), rtn(t, fixtures.Day1, .1), rtn(t, fixtures.Day0, .1)},
		})
	})

	t.Run("when adding list with no overlap", func(t *testing.T) {
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.Day1, .1)},
		})
		table = table.AddColumn(returns.List{
			rtn(t, fixtures.Day0, .1),
		})
		assert.Equal(t, table.Lists(), []returns.List{
			{},
			{},
		})
	})

	//t.Run("when adding to a sliced column", func(t *testing.T) {
	//	t.SkipNow()
	//	table := returns.NewTable([]returns.List{{}})
	//	slice := table.Between(fixtures.T(t, fixtures.Day1), fixtures.T(t, fixtures.Day0))
	//	o.Expect(func() {
	//		slice.AddColumn(returns.List{})
	//	}).To(Panic())
	//})
}

func TestTable_CorrelationMatrix(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		tab := returns.NewTable(nil)
		assert.Len(t, tab.CorrelationMatrix(), 0)
	})
	returnsFromQuotes := func(quotes ...float64) []float64 {
		if len(quotes) < 2 {
			return nil
		}
		result := make([]float64, len(quotes)-1)
		for i := 0; i < len(quotes)-1; i++ {
			result[i] = quotes[i]/quotes[i+1] - 1
		}
		return result
	}
	t.Run("perfectly positively correlated", func(t *testing.T) {
		rs1 := returnsFromQuotes(10.00, 20.00, 10.00, 20.00)
		rs2 := returnsFromQuotes(10.00, 20.00, 10.00, 20.00)
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.Day2, rs1[2]), rtn(t, fixtures.Day1, rs1[1]), rtn(t, fixtures.Day0, rs1[0])},
			{rtn(t, fixtures.Day2, rs2[2]), rtn(t, fixtures.Day1, rs2[1]), rtn(t, fixtures.Day0, rs2[0])},
		})
		assert.Equal(t, table.CorrelationMatrix(), [][]float64{
			{1, 1},
			{1, 1},
		})
	})
	t.Run("perfectly negatively correlated", func(t *testing.T) {
		rs1 := returnsFromQuotes(10.00, 20.00, 10.00, 20.00)
		rs2 := returnsFromQuotes(20.00, 10.00, 20.00, 10.00)
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.Day2, rs1[2]), rtn(t, fixtures.Day1, rs1[1]), rtn(t, fixtures.Day0, rs1[0])},
			{rtn(t, fixtures.Day2, rs2[2]), rtn(t, fixtures.Day1, rs2[1]), rtn(t, fixtures.Day0, rs2[0])},
		})
		assert.Equal(t, table.CorrelationMatrix(), [][]float64{
			{1, -1},
			{-1, 1},
		})
	})
}

func TestTable_Risks(t *testing.T) {
	table := returns.NewTable([]returns.List{
		{rtn(t, fixtures.Day2, -0.02), rtn(t, fixtures.Day1, 0.03), rtn(t, fixtures.Day0, -0.01)},
		{rtn(t, fixtures.Day2, +0.03), rtn(t, fixtures.Day1, 0.01), rtn(t, fixtures.Day0, +0.01)},
	})
	result := table.RisksFromStdDev()
	_ = round.Recursive(result, 4)
	assert.Equal(t, result, []float64{0.0265, 0.0115})
}

func TestTable_TimeWeightedReturns(t *testing.T) {
	table := returns.NewTable([]returns.List{
		{rtn(t, fixtures.Day2, -0.01), rtn(t, fixtures.Day1, 0.03), rtn(t, fixtures.Day0, -0.02)},
		{rtn(t, fixtures.Day2, +0.00), rtn(t, fixtures.Day1, 0.00), rtn(t, fixtures.Day0, +0.01)},
	})
	result := table.TimeWeightedReturns()
	_ = round.Recursive(result, 4)
	assert.Equal(t, result, []float64{-0.0566, 1.3067})
}

func TestTable_AnnualizedArithmeticReturns(t *testing.T) {
	table := returns.NewTable([]returns.List{
		{rtn(t, fixtures.Day2, -0.01), rtn(t, fixtures.Day1, 0.03), rtn(t, fixtures.Day0, -0.02)},
		{rtn(t, fixtures.Day2, +0.00), rtn(t, fixtures.Day1, 0.00), rtn(t, fixtures.Day0, +0.01)},
	})
	result := table.AnnualizedArithmeticReturns()
	_ = round.Recursive(result, 4)
	assert.Equal(t, result, []float64{0, 0.84})
}

func TestTable_ExpectedRisk(t *testing.T) {
	table := returns.NewTable([]returns.List{
		{rtn(t, fixtures.Day2, -0.01), rtn(t, fixtures.Day1, 0.03), rtn(t, fixtures.Day0, -0.02)},
		{rtn(t, fixtures.Day2, +0.00), rtn(t, fixtures.Day1, -0.01), rtn(t, fixtures.Day0, +0.01)},
	})
	result := table.ExpectedRisk([]float64{0.5, 0.5})
	result = round.Decimal(result, 4)
	assert.Equal(t, result, 0.0087)
	risks := table.RisksFromStdDev()
	assert.Less(t, result, risks[0])
	assert.Less(t, result, risks[1])
}

func TestTable_TimeAfter(t *testing.T) {
	t.Run("after data", func(t *testing.T) {
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
		})
		_, hasReturn := table.TimeAfter(fixtures.T(t, fixtures.DayAfter))
		assert.False(t, hasReturn)
	})
	t.Run("before data", func(t *testing.T) {
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
		})
		tm, hasReturn := table.TimeAfter(fixtures.T(t, fixtures.DayBefore))
		assert.True(t, hasReturn)
		assert.Equal(t, tm, fixtures.T(t, fixtures.FirstDay))
	})
	t.Run("on a friday", func(t *testing.T) {
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
		})
		after, hasReturn := table.TimeAfter(fixtures.T(t, fixtures.Day1))
		assert.True(t, hasReturn)
		assert.Equal(t, after, fixtures.T(t, fixtures.Day2))
	})
	t.Run("on a monday", func(t *testing.T) {
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
		})
		after, hasReturn := table.TimeAfter(fixtures.T(t, fixtures.Day2))
		assert.True(t, hasReturn)
		assert.Equal(t, after, fixtures.T(t, fixtures.Day3))
	})
}

func TestTable_TimeBefore(t *testing.T) {
	t.Run("after data", func(t *testing.T) {
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
		})
		tm, hasReturn := table.TimeBefore(fixtures.T(t, fixtures.DayAfter))
		assert.True(t, hasReturn)
		assert.Equal(t, tm, fixtures.T(t, fixtures.LastDay))
	})
	t.Run("before data", func(t *testing.T) {
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
		})
		_, hasReturn := table.TimeBefore(fixtures.T(t, fixtures.DayBefore))
		assert.False(t, hasReturn)
	})
	t.Run("on a Monday", func(t *testing.T) {
		in := fixtures.T(t, fixtures.Day2)
		require.Equal(t, time.Monday, in.Weekday())
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
		})
		result, hasReturn := table.TimeBefore(in)
		assert.True(t, hasReturn)
		assert.Equal(t, fixtures.T(t, fixtures.Day1), result)
	})
	t.Run("on a Friday", func(t *testing.T) {
		in := fixtures.T(t, fixtures.Day1)
		require.Equal(t, in.Weekday(), time.Friday)
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
		})
		result, hasReturn := table.TimeBefore(in)
		assert.True(t, hasReturn)
		assert.Equal(t, fixtures.T(t, fixtures.Day0), result)
	})
}

func TestTable_ClosestTimeOnOrBefore(t *testing.T) {
	t.Run("on a Friday", func(t *testing.T) {
		in := fixtures.T(t, fixtures.Day1)
		require.Equal(t, in.Weekday(), time.Friday)
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.LastDay, 0), rtn(t, fixtures.Day2, 0), rtn(t, fixtures.Day1, 0), rtn(t, fixtures.FirstDay, 0)},
		})
		result, hasReturn := table.ClosestTimeOnOrBefore(in)
		assert.True(t, hasReturn)
		assert.Equal(t, fixtures.T(t, fixtures.Day1), result)
	})
	t.Run("exactly between", func(t *testing.T) {
		in := fixtures.T(t, fixtures.Day2)
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.Day3, 0), rtn(t, fixtures.Day1, 0)},
		})
		result, hasReturn := table.ClosestTimeOnOrBefore(in)
		assert.True(t, hasReturn)
		assert.Equal(t, fixtures.T(t, fixtures.Day1), result)
	})
	t.Run("between closer to final day", func(t *testing.T) {
		in := fixtures.T(t, fixtures.Day2)
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.Day3, 0), rtn(t, fixtures.Day0, 0)},
		})
		result, hasReturn := table.ClosestTimeOnOrBefore(in)
		assert.True(t, hasReturn)
		assert.Equal(t, fixtures.T(t, fixtures.Day0), result)
	})
	t.Run("between closer to first day", func(t *testing.T) {
		in := fixtures.T(t, fixtures.Day1)
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.Day3, 0), rtn(t, fixtures.Day0, 0)},
		})
		result, hasReturn := table.ClosestTimeOnOrBefore(in)
		assert.True(t, hasReturn)
		assert.Equal(t, fixtures.T(t, fixtures.Day0), result)
	})
	t.Run("exactly first", func(t *testing.T) {
		in := fixtures.T(t, fixtures.Day0)
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.Day1, 0), rtn(t, fixtures.Day0, 0)},
		})
		result, hasReturn := table.ClosestTimeOnOrBefore(in)
		assert.True(t, hasReturn)
		assert.Equal(t, fixtures.T(t, fixtures.Day0), result)
	})
	t.Run("exactly last", func(t *testing.T) {
		in := fixtures.T(t, fixtures.Day1)
		table := returns.NewTable([]returns.List{
			{rtn(t, fixtures.Day1, 0), rtn(t, fixtures.Day0, 0)},
		})
		result, hasReturn := table.ClosestTimeOnOrBefore(in)
		assert.True(t, hasReturn)
		assert.Equal(t, fixtures.T(t, fixtures.Day1), result)
	})
}

func TestTable_Lists(t *testing.T) {
	table := returns.NewTable([]returns.List{
		{ /*                         ,*/ rtn(t, fixtures.Day2, 0.01), rtn(t, fixtures.Day1, -0.01), rtn(t, fixtures.Day0, 0.001)},
		{rtn(t, fixtures.Day3, 0.02), rtn(t, fixtures.Day2, -0.02), rtn(t, fixtures.Day1, 0.002)},
	})

	lists := table.Lists()
	assert.Len(t, lists, 2)
	assert.Equal(t, table.ColumnValues(), [][]float64{
		{0.01, -0.01},
		{-0.02, 0.002},
	})
}

func TestTable_Join(t *testing.T) {
	a := returns.NewTable([]returns.List{
		{ /*                         ,*/ rtn(t, fixtures.Day2, 0.01), rtn(t, fixtures.Day1, -0.01), rtn(t, fixtures.Day0, 0.001)},
	})
	b := returns.NewTable([]returns.List{
		{rtn(t, fixtures.Day3, 0.02), rtn(t, fixtures.Day2, -0.02), rtn(t, fixtures.Day1, 0.002)},
	})
	table := a.Join(b)

	assert.Equal(t, table.NumberOfRows(), 2)
	assert.Equal(t, table.ColumnValues(), [][]float64{
		{+0.01, -0.01},
		{-0.02, 0.002},
	})
}

func TestColumnGroup(t *testing.T) {
	t.Run("when group returns are outside of table range", func(t *testing.T) {
		group, updated := returns.NewTable([]returns.List{
			{rtn(t, fixtures.Day1, 100), rtn(t, fixtures.Day0, 420)},
			{rtn(t, fixtures.Day1, 100), rtn(t, fixtures.Day0, 420)},
		}).AddColumnGroup([]returns.List{
			{rtn(t, fixtures.Day1, 1), rtn(t, fixtures.Day0, 2)},
			{rtn(t, fixtures.Day1, 3), rtn(t, fixtures.Day0, 4)},
		})
		updated = updated.AddColumn(returns.List{rtn(t, fixtures.Day1, 9000)})

		assert.Equal(t, updated.ColumnValues(), [][]float64{
			{100},
			{100},
			{1},
			{3},
			{9000},
		})

		groupReturns := updated.ColumnGroupLists(group)

		assert.Equal(t, group.Length(), 2)

		assert.Len(t, groupReturns, 2)
		assert.Equal(t, groupReturns[0], returns.List{
			rtn(t, fixtures.Day1, 1),
		})
		assert.Equal(t, groupReturns[1], returns.List{
			rtn(t, fixtures.Day1, 3),
		})
	})

	t.Run("when returns outside of group time range", func(t *testing.T) {
		table := returns.NewTable(nil)
		table = table.AddColumn(returns.List{rtn(t, fixtures.Day0, 420)})
		group, updated := table.AddColumnGroup([]returns.List{
			{rtn(t, fixtures.Day2, 3), rtn(t, fixtures.Day1, 2)},
			{rtn(t, fixtures.Day2, 4), rtn(t, fixtures.Day1, 5)},
		})
		updated = updated.AddColumn(returns.List{rtn(t, fixtures.Day3, 420)})

		groupReturns := updated.ColumnGroupLists(group)
		assert.Len(t, groupReturns, 2)
		assert.Len(t, groupReturns[0], 0)
		assert.Len(t, groupReturns[1], 0)
	})
}

func TestTable_HasRow(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var table returns.Table

		assert.False(t, table.HasRow(fixtures.T(t, fixtures.Day0)))
	})

	t.Run("one return", func(t *testing.T) {
		table := returns.NewTable([]returns.List{
			{returns.New(fixtures.T(t, fixtures.Day0), 0)},
		})

		assert.True(t, table.HasRow(fixtures.T(t, fixtures.Day0)))
	})

	t.Run("between", func(t *testing.T) {
		table := returns.NewTable([]returns.List{
			{returns.New(fixtures.T(t, fixtures.Day0), 0)},
			{returns.New(fixtures.T(t, fixtures.Day2), 0)},
		})

		assert.False(t, table.HasRow(fixtures.T(t, fixtures.Day1)))
	})
}

func TestTable_ColumnGroupValues(t *testing.T) {
	var table returns.Table
	group, updated := table.AddColumnGroup([]returns.List{
		{rtn(t, fixtures.Day2, 3), rtn(t, fixtures.Day1, 2)},
		{rtn(t, fixtures.Day2, 4), rtn(t, fixtures.Day1, 5)},
	})
	assert.Equal(t, updated.ColumnGroupValues(group), [][]float64{
		{3, 2},
		{4, 5},
	})
}

func TestTable_Row(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var table returns.Table
		_, found := table.Row(fixtures.T(t, fixtures.Day0))
		assert.False(t, found)
	})
}

func TestDateAlignedReturns(t *testing.T) {
	t.Run("one asset has newer data", func(t *testing.T) {
		r1 := returns.List{
			{2, fixtures.T(t, "2020-01-05")},
			{2, fixtures.T(t, "2020-01-04")},
			{5, fixtures.T(t, "2020-01-03")},
			{1, fixtures.T(t, "2020-01-02")},
			{1, fixtures.T(t, "2020-01-01")},
		}

		r2 := returns.List{
			{2, fixtures.T(t, "2020-01-04")},
			{1, fixtures.T(t, "2020-01-03")},
			{3, fixtures.T(t, "2020-01-02")},
			{1, fixtures.T(t, "2020-01-01")},
		}

		table := returns.NewTable([]returns.List{r1, r2})

		assert.Equal(t, table.Lists(), returns.NewTable([]returns.List{
			{
				{2, fixtures.T(t, "2020-01-04")},
				{5, fixtures.T(t, "2020-01-03")},
				{1, fixtures.T(t, "2020-01-02")},
				{1, fixtures.T(t, "2020-01-01")},
			},
			{
				{2, fixtures.T(t, "2020-01-04")},
				{1, fixtures.T(t, "2020-01-03")},
				{3, fixtures.T(t, "2020-01-02")},
				{1, fixtures.T(t, "2020-01-01")},
			},
		}).Lists())
	})

	t.Run("both assets have same duration of data", func(t *testing.T) {
		r1 := returns.List{
			{2, fixtures.T(t, "2020-01-04")},
			{5, fixtures.T(t, "2020-01-03")},
			{1, fixtures.T(t, "2020-01-02")},
			{1, fixtures.T(t, "2020-01-01")},
		}

		r2 := returns.List{
			{2, fixtures.T(t, "2020-01-04")},
			{1, fixtures.T(t, "2020-01-03")},
			{3, fixtures.T(t, "2020-01-02")},
			{1, fixtures.T(t, "2020-01-01")},
		}

		table := returns.NewTable([]returns.List{r1, r2})

		assert.Equal(t, table.Lists(), returns.NewTable([]returns.List{
			{
				{2, fixtures.T(t, "2020-01-04")},
				{5, fixtures.T(t, "2020-01-03")},
				{1, fixtures.T(t, "2020-01-02")},
				{1, fixtures.T(t, "2020-01-01")},
			},
			{
				{2, fixtures.T(t, "2020-01-04")},
				{1, fixtures.T(t, "2020-01-03")},
				{3, fixtures.T(t, "2020-01-02")},
				{1, fixtures.T(t, "2020-01-01")},
			},
		}).Lists())
	})

	t.Run("one asset has longer history", func(t *testing.T) {
		r1 := returns.List{
			{2, fixtures.T(t, "2020-01-05")},
			{2, fixtures.T(t, "2020-01-04")},
			{5, fixtures.T(t, "2020-01-03")},
			{1, fixtures.T(t, "2020-01-02")},
		}

		r2 := returns.List{
			{4, fixtures.T(t, "2020-01-05")},
			{2, fixtures.T(t, "2020-01-04")},
			{1, fixtures.T(t, "2020-01-03")},
			{3, fixtures.T(t, "2020-01-02")},
			{1, fixtures.T(t, "2020-01-01")},
		}

		table := returns.NewTable([]returns.List{r1, r2})

		assert.Equal(t, table.Lists(), returns.NewTable([]returns.List{
			{
				{2, fixtures.T(t, "2020-01-05")},
				{2, fixtures.T(t, "2020-01-04")},
				{5, fixtures.T(t, "2020-01-03")},
				{1, fixtures.T(t, "2020-01-02")},
			},
			{
				{4, fixtures.T(t, "2020-01-05")},
				{2, fixtures.T(t, "2020-01-04")},
				{1, fixtures.T(t, "2020-01-03")},
				{3, fixtures.T(t, "2020-01-02")},
			},
		}).Lists())
	})

	t.Run("an asset has only one return", func(t *testing.T) {
		r1 := returns.List{
			{2, fixtures.T(t, "2020-01-05")},
		}

		r2 := returns.List{
			{4, fixtures.T(t, "2020-01-05")},
			{2, fixtures.T(t, "2020-01-04")},
			{1, fixtures.T(t, "2020-01-03")},
			{3, fixtures.T(t, "2020-01-02")},
		}

		table := returns.NewTable([]returns.List{r1, r2})

		assert.Equal(t, table.Lists(), returns.NewTable([]returns.List{
			{
				{2, fixtures.T(t, "2020-01-05")},
			},
			{
				{4, fixtures.T(t, "2020-01-05")},
			},
		}).Lists())
	})

	t.Run("one asset has no data", func(t *testing.T) {
		r1 := returns.List{
			{2, fixtures.T(t, "2020-01-04")},
			{5, fixtures.T(t, "2020-01-03")},
			{1, fixtures.T(t, "2020-01-02")},
			{1, fixtures.T(t, "2020-01-01")},
		}

		r2 := returns.List{}

		table := returns.NewTable([]returns.List{r1, r2})

		assert.Equal(t, table.NumberOfColumns(), 2)
		assert.Len(t, table.List(0), 0)
		assert.Len(t, table.List(1), 0)
	})

	t.Run("one asset has missing internal data", func(t *testing.T) {
		r1 := returns.List{
			{2, fixtures.T(t, "2020-01-04")},
			{5, fixtures.T(t, "2020-01-03")},
			{1, fixtures.T(t, "2020-01-02")},
			{1, fixtures.T(t, "2020-01-01")},
		}

		r2 := returns.List{
			{2, fixtures.T(t, "2020-01-04")},
			// {5, fixtures.T(t, "2020-01-03")},
			{1, fixtures.T(t, "2020-01-02")},
			{1, fixtures.T(t, "2020-01-01")},
		}

		table := returns.NewTable([]returns.List{r1, r2})

		assert.Equal(t, table.NumberOfColumns(), 2)
		assert.Equal(t, table.List(0), r1)
		assert.Equal(t, table.List(1), returns.List{
			{2, fixtures.T(t, "2020-01-04")},
			{0, fixtures.T(t, "2020-01-03")},
			{1, fixtures.T(t, "2020-01-02")},
			{1, fixtures.T(t, "2020-01-01")},
		})
	})
}

func TestTable_WriteCSV(t *testing.T) {
	t.Run("exponent number", func(t *testing.T) {
		expNum := -9.1e-05

		table := returns.NewTable([]returns.List{{{Time: fixtures.T(t, fixtures.Day1), Value: expNum}}})

		var buf bytes.Buffer
		err := table.WriteCSV(&buf, []string{"mango"})

		require.NoError(t, err)
		assert.NotContains(t, buf.String(), "e-")
	})
}
