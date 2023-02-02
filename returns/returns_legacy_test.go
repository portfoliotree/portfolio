package returns_test

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"

	"github.com/portfoliotree/portfolio/internal/fixtures"
	"github.com/portfoliotree/portfolio/returns"
)

func TestReturns(t *testing.T) {
	t.Run("Risk", func(t *testing.T) {
		table := returns.List{{Value: 1}, {Value: .5}, {Value: 1.5}}
		assert.Equal(t, table.Risk(), 0.5)
	})

	t.Run("Len", func(t *testing.T) {
		table := returns.List{{Value: 1}, {Value: 1}, {Value: 1}, {Value: 1}}
		assert.Equal(t, table.Len(), 4)
	})

	t.Run("returns.Table", func(t *testing.T) {
		table := returns.List{{Value: 1}, {Value: 1}, {Value: 1}, {Value: 1}}
		assert.Equal(t, table.Returns().Values(), []float64{1, 1, 1, 1})
	})
}

func TestReturns_FirstAndLastPeriod(t *testing.T) {

	rs := returns.List{
		{2, date("2020-01-04")},
		{1, date("2020-01-03")},
		{3, date("2020-01-02")},
	}

	end, start, err := rs.EndAndStartDate()
	assert.NoError(t, err)
	assert.Equal(t, start, date("2020-01-02"))
	assert.Equal(t, end, date("2020-01-04"))
}

func toSlices(dense *mat.Dense) [][]float64 {
	iL, jL := dense.Dims()
	result := make([][]float64, iL)
	d := dense.RawMatrix().Data
	for i := range result {
		result[i] = d[i*jL : (i+1)*jL]
	}
	return result
}

func TestComposite_correlationMatrix(t *testing.T) {
	t.Run("perfectly positively correlated", func(t *testing.T) {

		assert.Equal(t,
			toSlices(returns.NewTable([]returns.List{
				{{Time: fixtures.T(t, fixtures.Day3), Value: 10}, {Time: fixtures.T(t, fixtures.Day2), Value: 20}, {Time: fixtures.T(t, fixtures.Day1), Value: 10}, {Time: fixtures.T(t, fixtures.Day0), Value: 20}},
				{{Time: fixtures.T(t, fixtures.Day3), Value: 10}, {Time: fixtures.T(t, fixtures.Day2), Value: 20}, {Time: fixtures.T(t, fixtures.Day1), Value: 10}, {Time: fixtures.T(t, fixtures.Day0), Value: 20}},
			}).CorrelationMatrix()),
			[][]float64{
				{1, 1},
				{1, 1},
			},
		)
	})

	t.Run("perfectly negatively correlated", func(t *testing.T) {

		assert.Equal(t,
			toSlices(returns.NewTable([]returns.List{
				{{Time: fixtures.T(t, fixtures.Day3), Value: 10}, {Time: fixtures.T(t, fixtures.Day2), Value: 20}, {Time: fixtures.T(t, fixtures.Day1), Value: 10}, {Time: fixtures.T(t, fixtures.Day0), Value: 20}},
				{{Time: fixtures.T(t, fixtures.Day3), Value: 20}, {Time: fixtures.T(t, fixtures.Day2), Value: 10}, {Time: fixtures.T(t, fixtures.Day1), Value: 20}, {Time: fixtures.T(t, fixtures.Day0), Value: 10}},
			}).CorrelationMatrix()),
			[][]float64{
				{1, -1},
				{-1, 1},
			},
		)
	})

	t.Run("about halfish correlated", func(t *testing.T) {
		c := returns.NewTable([]returns.List{
			{{Time: fixtures.T(t, fixtures.Day2), Value: 0.1}, {Time: fixtures.T(t, fixtures.Day1), Value: -0.1}, {Time: fixtures.T(t, fixtures.Day0), Value: .1}},
			{{Time: fixtures.T(t, fixtures.Day2), Value: -0.1}, {Time: fixtures.T(t, fixtures.Day1), Value: -0.1}, {Time: fixtures.T(t, fixtures.Day0), Value: .1}},
		}).CorrelationMatrix()

		roughlyEqual(t, c.At(0, 1), 0.5)
		roughlyEqual(t, c.At(1, 0), 0.5)
	})
}

func TestDateAlignedReturnsList_ExpectedRisk(t *testing.T) {
	t.Run("with one asset", func(t *testing.T) {

		list := returns.NewTable([]returns.List{{
			{Time: fixtures.T(t, fixtures.Day2), Value: 1},
			{Time: fixtures.T(t, fixtures.Day1), Value: -2 / 3},
			{Time: fixtures.T(t, fixtures.Day0), Value: .5},
		}})

		assert.Equal(t, list.ExpectedRisk([]float64{1}), list.List(0).Risk())
	})

	t.Run("with two assets and one has no weight", func(t *testing.T) {

		list := returns.NewTable([]returns.List{
			{{Value: 1}, {Value: -2 / 3}, {Value: .5}},
			{{Value: .5}, {Value: .3}, {Value: .5}},
		})

		assert.Equal(t, list.ExpectedRisk([]float64{1, 0}), list.List(0).Risk())
	})

	t.Run("with two completely correlated assets", func(t *testing.T) {

		list := returns.NewTable([]returns.List{
			{{Value: 1}, {Value: -2 / 3}, {Value: .5}},
			{{Value: 1}, {Value: -2 / 3}, {Value: .5}},
		})

		compositeRisk := list.ExpectedRisk([]float64{0.2, 0.8})

		assert.Equal(t, compositeRisk, list.List(0).Risk())
		assert.Equal(t, compositeRisk, list.List(1).Risk())
	})

	t.Run("with equal risk contribution", func(t *testing.T) {
		list := returns.NewTable([]returns.List{
			{{Value: 1.5}, {Value: -0.25}, {Value: 1}},
			{{Value: 1.5}, {Value: -0.25}, {Value: 1}},
		})

		compositeRisk := list.ExpectedRisk([]float64{0.5, 0.5})

		roughlyEqual(t, compositeRisk, list.List(0).Risk())
		roughlyEqual(t, compositeRisk, list.List(1).Risk())
	})

	t.Run("with scaled but correlated assets", func(t *testing.T) {
		list := returns.NewTable([]returns.List{
			{{Value: 1.5}, {Value: -0.25}, {Value: 1}},
			{{Value: 3}, {Value: -0.5}, {Value: 2}},
		})

		weights := []float64{.5, .5}
		weightedAverage := list.List(0).Risk()*weights[0] + list.List(1).Risk()*weights[1]

		roughlyEqual(t, list.ExpectedRisk(weights), weightedAverage)
	})
}

func TestDateAlignedReturnsList_FirstAndLastSharedPeriod(t *testing.T) {

	list := returns.NewTable([]returns.List{
		{
			{2, date("2020-01-05")},
			{2, date("2020-01-04")},
			{5, date("2020-01-03")},
			{1, date("2020-01-02")},
		},
		{
			{2, date("2020-01-04")},
			{1, date("2020-01-03")},
			{3, date("2020-01-02")},
			{3, date("2020-01-01")},
		},
	})

	end, start, err := list.EndAndStartDates()
	assert.NoError(t, err)
	assert.Equal(t, start, date("2020-01-02"))
	assert.Equal(t, end, date("2020-01-04"))
}

func TestReturnsList_FirstAndLastSharedPeriod_no_overlap(t *testing.T) {

	list := returns.NewTable([]returns.List{
		{
			{2, date("2020-01-05")},
			{2, date("2020-01-04")},
		},
		{
			{5, date("2020-01-03")},
			{1, date("2020-01-02")},
			{1, date("2020-01-01")},
		},
	})

	_, _, err := list.EndAndStartDates()
	assert.NoError(t, err)
}

func TestReturns_TruncateToDateRange(t *testing.T) {

	table := returns.List{
		{Time: date("2021-06-25"), Value: 1.0},
		{Time: date("2021-06-24"), Value: 1.0},
		{Time: date("2021-06-23"), Value: 1.0},
		{Time: date("2021-06-22"), Value: 1.0},
		{Time: date("2021-06-21"), Value: 1.0},
		{Time: date("2021-06-18"), Value: 1.0},
		{Time: date("2021-06-17"), Value: 1.0},
	}

	rs := table.Between(date("2021-06-22"), date("2021-06-18"))

	assert.Equal(t, rs, returns.List{
		{Time: date("2021-06-22"), Value: 1.0},
		{Time: date("2021-06-21"), Value: 1.0},
		{Time: date("2021-06-18"), Value: 1.0},
	})
}

func TestReturns_TruncateToDateRange_end_is_after_final_return(t *testing.T) {

	table := returns.List{
		{Time: date("2021-06-25"), Value: 1.0},
		{Time: date("2021-06-24"), Value: 1.0},
		{Time: date("2021-06-23"), Value: 1.0},
		{Time: date("2021-06-22"), Value: 1.0},
		{Time: date("2021-06-21"), Value: 1.0},
		{Time: date("2021-06-18"), Value: 1.0},
		{Time: date("2021-06-17"), Value: 1.0},
	}

	rs := table.Between(date("2021-06-30"), date("2021-06-24"))

	assert.Equal(t, rs, returns.List{
		{Time: date("2021-06-25"), Value: 1.0},
		{Time: date("2021-06-24"), Value: 1.0},
	})
}

func TestReturns_TruncateToDateRange_start_is_before_initial_return(t *testing.T) {

	table := returns.List{
		{Time: date("2021-06-25"), Value: 1.0},
		{Time: date("2021-06-24"), Value: 1.0},
		{Time: date("2021-06-23"), Value: 1.0},
		{Time: date("2021-06-22"), Value: 1.0},
		{Time: date("2021-06-21"), Value: 1.0},
		{Time: date("2021-06-18"), Value: 1.0},
		{Time: date("2021-06-17"), Value: 1.0},
	}

	rs := table.Between(date("2021-06-21"), date("2021-06-01"))

	assert.Equal(t, rs, returns.List{
		{Time: date("2021-06-21"), Value: 1.0},
		{Time: date("2021-06-18"), Value: 1.0},
		{Time: date("2021-06-17"), Value: 1.0},
	})
}

func TestReturns_TruncateToDateRange_start_and_end_are_beyond_return_range(t *testing.T) {

	table := returns.List{
		{Time: date("2021-06-25"), Value: 1.0},
		{Time: date("2021-06-24"), Value: 1.0},
		{Time: date("2021-06-23"), Value: 1.0},
		{Time: date("2021-06-22"), Value: 1.0},
		{Time: date("2021-06-21"), Value: 1.0},
		{Time: date("2021-06-18"), Value: 1.0},
		{Time: date("2021-06-17"), Value: 1.0},
	}

	rs := table.Between(date("2021-06-30"), date("2021-06-01"))

	assert.Equal(t, rs, returns.List{
		{Time: date("2021-06-25"), Value: 1.0},
		{Time: date("2021-06-24"), Value: 1.0},
		{Time: date("2021-06-23"), Value: 1.0},
		{Time: date("2021-06-22"), Value: 1.0},
		{Time: date("2021-06-21"), Value: 1.0},
		{Time: date("2021-06-18"), Value: 1.0},
		{Time: date("2021-06-17"), Value: 1.0},
	})
}

func TestTruncateToOverlappingPeriods(t *testing.T) {

	lists := []returns.Table{
		returns.NewTable([]returns.List{
			{
				{4, date("2020-01-09")},
				{4, date("2020-01-08")},
				{4, date("2020-01-05")},
				{2, date("2020-01-04")},
				{1, date("2020-01-03")},
				{3, date("2020-01-02")},
				{1, date("2020-01-01")},
			},
			{
				{4, date("2020-01-08")},
				{4, date("2020-01-05")},
				{2, date("2020-01-04")},
				{1, date("2020-01-03")},
				{3, date("2020-01-02")},
				{1, date("2020-01-01")},
			},
		}),
		returns.NewTable([]returns.List{
			{
				{4, date("2020-01-10")},
				{4, date("2020-01-09")},
				{4, date("2020-01-08")},
				{4, date("2020-01-05")},
				{2, date("2020-01-04")},
				{1, date("2020-01-03")},
			},
			{
				{4, date("2020-01-10")},
				{4, date("2020-01-09")},
				{4, date("2020-01-08")},
				{4, date("2020-01-05")},
				{2, date("2020-01-04")},
				{1, date("2020-01-03")},
				{3, date("2020-01-02")},
			},
		}),
	}

	truncated, end, start, err := returns.AlignTables(lists...)
	assert.NoError(t, err)
	assert.Equal(t, end, date("2020-01-08"))
	assert.Equal(t, start, date("2020-01-03"))

	for _, list := range truncated {
		for _, table := range list.Lists() {
			e, s, err := table.EndAndStartDate()
			assert.NoError(t, err)
			assert.Equal(t, e, date("2020-01-08"))
			assert.Equal(t, s, date("2020-01-03"))
		}
	}
}

func date(str string) time.Time {
	d, err := time.Parse(time.DateOnly, str)
	if err != nil {
		panic(fmt.Errorf("failed to parse date %q: %s", str, err))
	}
	return d
}

func datef(format string, a ...interface{}) time.Time {
	return date(fmt.Sprintf(format, a...))
}

func roughlyEqual(t *testing.T, n1, n2 float64) {
	t.Helper()

	t.Logf("%f should be really damn close to %f", n1, n2)

	if math.IsNaN(n1) {
		t.Logf("%f should not be NaN", n1)
		t.Fail()
	}

	if math.IsNaN(n2) {
		t.Logf("%f should not be NaN", n2)
		t.Fail()
	}

	if math.Abs(n1-n2) > 0.0001 {
		t.Fail()
	}
}

func TestReturns_Within(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var list returns.List
		assert.Len(t, list.Between(date("2050-05-07"), date("2000-05-07")), 0)
	})

	t.Run("out of order", func(t *testing.T) {
		list := returns.List{{Time: date("2020-02-01")}}
		assert.Len(t, list.Between(date("2000-05-07"), date("2050-05-07")), 0)
	})

	t.Run("one value in range", func(t *testing.T) {
		list := returns.List{{Time: date("2020-02-01")}}
		assert.Len(t, list.Between(date("2050-05-07"), date("2000-05-07")), 1)
	})

	t.Run("dates are before returns", func(t *testing.T) {
		list := returns.List{{Time: date("2020-02-01")}}
		assert.Len(t, list.Between(date("2000-12-07"), date("2000-05-07")), 0)
	})

	t.Run("dates are after returns", func(t *testing.T) {
		list := returns.List{{Time: date("2020-02-01")}}
		assert.Len(t, list.Between(date("2050-12-07"), date("2050-05-07")), 0)
	})

	t.Run("exactly one value for day", func(t *testing.T) {
		list := returns.List{
			{
				Time: date("2020-04-01"),
			},
			{
				Time: date("2020-03-15"),
			},
			{
				Time: date("2020-02-01"),
			},
		}
		sort.Sort(list)
		assert.Len(t, list.Between(date("2020-03-15"), date("2020-03-15")), 1)
	})

	t.Run("just the middle", func(t *testing.T) {
		list := returns.List{
			{
				Time: date("2020-04-01"),
			},
			{
				Time: date("2020-03-15"),
			},
			{
				Time: date("2020-02-01"),
			},
		}
		sort.Sort(list)
		assert.Len(t, list.Between(date("2020-03-20"), date("2020-03-10")), 1)
	})

	t.Run("times out of order", func(t *testing.T) {
		list := returns.List{
			{
				Time: date("2020-04-01"),
			},
			{
				Time: date("2020-03-01"),
			},
			{
				Time: date("2020-02-01"),
			},
		}
		sort.Sort(list)
		assert.Len(t, list.Between(date("1999-03-20"), date("2050-03-10")), 0)
	})

	t.Run("range around a single return", func(t *testing.T) {
		list := returns.List{
			{
				Time: date("2020-05-15"),
			},
			{
				Time: date("2020-04-15"),
			},
			{
				Time: date("2020-03-15"),
			},
			{
				Time: date("2020-02-15"),
			},
			{
				Time: date("2020-01-15"),
			},
		}

		for index, ret := range list {
			t.Run(strconv.Itoa(index), func(t *testing.T) {
				month := len(list) - index
				t.Logf("month: %02d", month)
				subsection := list.Between(datef("2020-%02d-20", month), datef("2020-%02d-01", month))
				assert.Len(t, subsection, 1)
				assert.Equal(t, subsection[0].Time, ret.Time)
			})
		}
	})
}

func TestReturns_ExcessReturns(t *testing.T) {

	list := returns.List{
		{Value: 420, Time: date("2021-10-07")}, // out of time range of o

		{Value: 0.01, Time: date("2021-10-06")},
		{Value: 0.02, Time: date("2021-10-05")},
		{Value: 0.04, Time: date("2021-10-04")},
	}
	other := returns.List{
		{Value: -0.01, Time: date("2021-10-06")},
		{Value: 0.02, Time: date("2021-10-05")},
		{Value: 0.02, Time: date("2021-10-04")},

		{Value: 90000, Time: date("2021-10-01")}, // out of time range of r
	}

	got := list.Excess(other)

	assert.Equal(t, got, returns.List{
		{Value: 0.02, Time: date("2021-10-06")},
		{Value: 0.00, Time: date("2021-10-05")},
		{Value: 0.02, Time: date("2021-10-04")},
	})
}

func FuzzReturnsAnnualizedRisk(f *testing.F) {
	f.Add(returnsToJSON(f, nil))
	f.Add(returnsToJSON(f, returns.List{}))
	f.Add(returnsToJSON(f, returns.List{{Time: date("2021-05-26"), Value: .1}}))
	f.Add(returnsToJSON(f, makeReturnsFromFloats(randomFloats(500, 11))))

	f.Fuzz(func(t *testing.T, buf []byte) {
		rs := returnsFromJSON(t, buf)
		risk := rs.AnnualizedRisk()
		if math.IsNaN(risk) {
			t.Errorf("got %f", risk)
		}
	})
}

func FuzzAnnualizedTimeWeightedReturn(f *testing.F) {
	f.Add(returnsToJSON(f, nil))
	f.Add(returnsToJSON(f, returns.List{}))
	f.Add(returnsToJSON(f, returns.List{{Time: date("2021-05-26"), Value: .1}}))
	f.Add(returnsToJSON(f, makeReturnsFromFloats(randomFloats(500, 11))))

	f.Fuzz(func(t *testing.T, buf []byte) {
		rs := returnsFromJSON(t, buf)
		risk := rs.AnnualizedTimeWeightedReturn()
		if math.IsNaN(risk) {
			t.Errorf("got %f", risk)
		}
	})
}

func FuzzAnnualizedArithmeticReturn(f *testing.F) {
	f.Add(returnsToJSON(f, nil))
	f.Add(returnsToJSON(f, returns.List{}))
	f.Add(returnsToJSON(f, returns.List{{Time: date("2021-05-26"), Value: .1}}))
	f.Add(returnsToJSON(f, makeReturnsFromFloats(randomFloats(500, 11))))

	f.Fuzz(func(t *testing.T, buf []byte) {
		rs := returnsFromJSON(t, buf)
		risk := rs.AnnualizedArithmeticReturn()
		if math.IsNaN(risk) {
			t.Errorf("got %f", risk)
		}
	})
}

func randomFloats(n int, seed int64) []float64 {
	r := rand.New(rand.NewSource(seed))
	rs := make([]float64, 0, n)
	for len(rs) < cap(rs) {
		rs = append(rs, (r.Float64()-.48)/100)
	}
	return rs
}

func makeReturnsFromFloats(numbers []float64) returns.List {
	t := date("2022-06-16")
	rs := make(returns.List, 0, len(numbers))
	for _, n := range numbers {
		t = t.AddDate(0, 0, -1)
		if wd := t.Weekday(); wd == time.Sunday || wd == time.Saturday {
			continue
		}
		rs = append(rs, returns.Return{
			Time:  t,
			Value: n,
		})
	}
	return rs
}

func returnsToJSON(f *testing.F, rs returns.List) []byte {
	f.Helper()
	buf, err := json.Marshal(rs)
	if err != nil {
		f.Fatal(err)
	}
	return buf
}

func returnsFromJSON(t *testing.T, buf []byte) returns.List {
	t.Helper()
	var rs returns.List
	err := json.Unmarshal(buf, &rs)
	if err != nil {
		t.Skip()
	}
	return rs
}
