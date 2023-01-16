package returns_test

import (
	"sort"
	"testing"
	"time"

	"github.com/portfoliotree/round"
	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/portfolio/internal/fixtures"
	"github.com/portfoliotree/portfolio/returns"
)

// var _ portfolio.Asset = returns.List{}

func TestList_Returns(t *testing.T) {
	list := returns.List{
		returns.New(fixtures.T(t, fixtures.Day0), 400),
	}

	assert.Equal(t, list.Returns(), returns.List{
		returns.New(fixtures.T(t, fixtures.Day0), 400),
	})
}

func TestList_Value(t *testing.T) {
	d := fixtures.T(t, "2021-12-31")

	t.Run("empty", func(t *testing.T) {
		var slice returns.List
		slice.Sort()
		updated, found := slice.Value(d)
		assert.False(t, found)
		assert.Zero(t, updated)
	})

	t.Run("not found", func(t *testing.T) {
		t.Run("times are not equal", func(t *testing.T) {
			var (
				existingReturn = returns.New(d.AddDate(1, 0, 0), .024)
				slice          = returns.List{
					existingReturn,
				}
			)
			slice.Sort()

			// when
			updated, found := slice.Value(d)

			// then
			assert.False(t, found)
			assert.Zero(t, updated)
		})
	})

	t.Run("found", func(t *testing.T) {
		t.Run("exact match one return", func(t *testing.T) {
			var (
				existingReturn = returns.New(d, .024)
				slice          = returns.List{
					existingReturn,
				}
			)
			slice.Sort()

			updated, found := slice.Value(d)

			assert.True(t, found)
			assert.Equal(t, updated, .024)
		})

		t.Run("in the middle", func(t *testing.T) {
			// given
			var (
				returnBefore = returns.New(d.AddDate(-1, 0, 0), .1)
				returnAfter  = returns.New(d.AddDate(1, 0, 0), .2)
				targetReturn = returns.New(d, .666)
				slice        = returns.List{
					returnAfter, targetReturn, returnBefore,
				}
			)
			slice.Sort()

			// when
			updated, found := slice.Value(d)

			// then
			assert.True(t, found)
			assert.Equal(t, updated, .666)
		})

		t.Run("at the beginning", func(t *testing.T) {
			// given
			var (
				returnBefore        = returns.New(d.AddDate(-1, 0, 0), .1)
				returnFurtherBefore = returns.New(d.AddDate(-2, 0, 0), .2)
				targetReturn        = returns.New(d, .666)
				slice               = returns.List{
					targetReturn, returnBefore, returnFurtherBefore,
				}
			)
			slice.Sort()

			// when
			updated, found := slice.Value(d)

			// then
			assert.True(t, found)
			assert.Equal(t, updated, .666)
		})

		t.Run("at the end", func(t *testing.T) {
			// given
			var (
				returnAfter        = returns.New(d.AddDate(1, 0, 0), .1)
				returnFurtherAfter = returns.New(d.AddDate(2, 0, 0), .2)
				targetReturn       = returns.New(d, .666)
				slice              = returns.List{
					returnFurtherAfter, returnAfter, targetReturn,
				}
			)
			slice.Sort()

			// when
			updated, found := slice.Value(d)

			// then
			assert.True(t, found)
			assert.Equal(t, updated, .666)
		})
	})
}

func TestList_Values(t *testing.T) {
	rs := returns.List{
		returns.New(fixtures.T(t, "2022-03-01"), 0.3),
		returns.New(fixtures.T(t, "2022-02-01"), 0.2),
		returns.New(fixtures.T(t, "2022-01-01"), 0.1),
	}

	assert.Equal(t, rs.Values(), []float64{0.3, 0.2, 0.1})
}

func TestList_Sort(t *testing.T) {
	d := fixtures.T(t, "2022-01-01")
	t.Run("out of order", func(t *testing.T) {
		// given
		d0, d1, d2 := d.AddDate(-2, 0, 0), d.AddDate(-1, 0, 0), d
		slice := returns.List{
			returns.New(d1, 0.2),
			returns.New(d0, 0.1),
			returns.New(d2, 0.3),
		}

		// when
		slice.Sort()

		// then
		assert.Equal(t, slice, returns.List{
			returns.New(d2, 0.3),
			returns.New(d1, 0.2),
			returns.New(d0, 0.1),
		})
	})
}

func TestList_Insert(t *testing.T) {
	d := fixtures.T(t, "2021-12-31")

	t.Run("empty", func(t *testing.T) {
		var (
			newReturn = returns.New(d, .001)
			slice     returns.List
		)

		updated := slice.Insert(newReturn)

		assert.Len(t, updated, 1)
		assert.Contains(t, updated, newReturn)
		assert.True(t, sort.IsSorted(updated))
	})

	t.Run("before", func(t *testing.T) {
		// given
		var (
			newReturn      = returns.New(d, 0)
			existingReturn = returns.New(d.AddDate(1, 0, 0), 0)
			slice          = returns.List{
				existingReturn,
			}
		)

		// when
		updated := slice.Insert(newReturn)

		// then
		assert.Len(t, updated, 2)
		assert.Contains(t, updated, existingReturn)
		assert.Contains(t, updated, newReturn)
		assert.True(t, sort.IsSorted(updated))
	})

	t.Run("after", func(t *testing.T) {
		// given
		var (
			newReturn      = returns.New(d, 0)
			existingReturn = returns.New(d.AddDate(-1, 0, 0), 0)
			slice          = returns.List{
				existingReturn,
			}
		)

		// when
		updated := slice.Insert(newReturn)

		// then
		assert.Equal(t, updated, returns.List{
			newReturn,
			existingReturn,
		})
	})

	t.Run("between", func(t *testing.T) {
		// given
		var (
			newReturn           = returns.New(d, 0)
			existingReturnOlder = returns.New(d.AddDate(-1, 0, 0), 0)
			existingReturnNewer = returns.New(d.AddDate(1, 0, 0), 0)
			slice               = returns.List{
				existingReturnNewer, existingReturnOlder,
			}
		)

		// when
		updated := slice.Insert(newReturn)

		// then
		assert.Equal(t, updated, returns.List{
			existingReturnNewer, newReturn, existingReturnOlder,
		})
	})

	t.Run("update", func(t *testing.T) {
		// given
		var (
			existingReturnNewer = returns.New(d, 5)
			newReturn           = returns.New(d, 7)
			slice               = returns.List{
				existingReturnNewer,
			}
		)

		// when
		updated := slice.Insert(newReturn)

		// then
		assert.Equal(t, updated, returns.List{newReturn})
	})
}

func TestList_First(t *testing.T) {
	d := fixtures.T(t, "2021-12-31")

	t.Run("empty", func(t *testing.T) {
		var list returns.List

		assert.NotPanics(t, func() {
			list.First()
		})
	})

	t.Run("two returns", func(t *testing.T) {
		// given
		list := returns.List{
			returns.New(d.AddDate(1, 0, 0), 0.2),
			returns.New(d, .1),
		}

		// when
		r := list.First()
		rt := list.FirstTime()

		// then
		assert.Equal(t, rt, d)
		assert.Equal(t, r, returns.New(d, .1))
	})
}

func TestList_Last(t *testing.T) {
	d := fixtures.T(t, "2021-12-31")

	t.Run("empty", func(t *testing.T) {
		var list returns.List
		assert.NotPanics(t, func() {
			list.Last()
		})
	})

	t.Run("two returns", func(t *testing.T) {
		// given
		list := returns.List{
			returns.New(d, 0.2),
			returns.New(d.AddDate(-1, 0, 0), .1),
		}

		// when
		r := list.Last()
		rt := list.LastTime()

		// then
		assert.Equal(t, r, returns.New(d, .2))
		assert.Equal(t, rt, d)
	})
}

func TestList_Excess(t *testing.T) {
	result := returns.List{
		rtn(t, fixtures.Day3, 0.4),
		rtn(t, fixtures.Day2, 0.2),
		rtn(t, fixtures.Day1, 0.1),
		rtn(t, fixtures.Day0, 0.6),
	}.Excess(returns.List{
		rtn(t, fixtures.Day2, 0.1),
		rtn(t, fixtures.Day1, 0.1),
	})
	_ = round.Recursive(&result, 1)
	assert.Equal(t, result, returns.List{
		rtn(t, fixtures.Day2, 0.1),
		rtn(t, fixtures.Day1, 0.0),
	})
}

func TestList_Between(t *testing.T) {
	t.Run("out of bounds", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day3, 0.4),
			rtn(t, fixtures.Day2, 0.2),
			rtn(t, fixtures.Day1, 0.1),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.DayAfter), fixtures.T(t, fixtures.DayBefore))
		assert.Equal(t, result, list)
	})

	t.Run("both before", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day3, 0.4),
			rtn(t, fixtures.Day2, 0.2),
			rtn(t, fixtures.Day1, 0.1),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.DayBefore), fixtures.T(t, fixtures.DayBefore).AddDate(-1, 0, 0))
		assert.Len(t, result, 0)
	})

	t.Run("both after", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day3, 0.4),
			rtn(t, fixtures.Day2, 0.2),
			rtn(t, fixtures.Day1, 0.1),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.DayAfter).AddDate(1, 0, 0), fixtures.T(t, fixtures.DayAfter))
		assert.Len(t, result, 0)
	})

	t.Run("same day", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day3, 0.4),
			rtn(t, fixtures.Day2, 0.2),
			rtn(t, fixtures.Day1, 0.1),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.Day1), fixtures.T(t, fixtures.Day1))
		assert.Equal(t, result, returns.List{
			rtn(t, fixtures.Day1, 0.1),
		})
	})

	t.Run("days between", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day3, 0.4),
			rtn(t, fixtures.Day2, 0.2),
			rtn(t, fixtures.Day1, 0.1),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.Day3), fixtures.T(t, fixtures.Day1))
		assert.Equal(t, result, returns.List{
			rtn(t, fixtures.Day3, 0.4),
			rtn(t, fixtures.Day2, 0.2),
			rtn(t, fixtures.Day1, 0.1),
		})
	})

	t.Run("a single return", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.Day0), fixtures.T(t, fixtures.Day0))
		assert.Equal(t, result, returns.List{
			rtn(t, fixtures.Day0, 0.6),
		})
	})

	t.Run("two returns", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day1, 0.6),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.Day1), fixtures.T(t, fixtures.Day0))
		assert.Equal(t, result, returns.List{
			rtn(t, fixtures.Day1, 0.6),
			rtn(t, fixtures.Day0, 0.6),
		})
	})

	t.Run("two returns", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day2, 0.6),
			rtn(t, fixtures.Day1, 0.6),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.Day2), fixtures.T(t, fixtures.Day0))
		assert.Equal(t, result, returns.List{
			rtn(t, fixtures.Day2, 0.6),
			rtn(t, fixtures.Day1, 0.6),
			rtn(t, fixtures.Day0, 0.6),
		})
	})

	t.Run("truncate the first return", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day2, 0.6),
			rtn(t, fixtures.Day1, 0.6),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.Day2), fixtures.T(t, fixtures.Day1))
		assert.Equal(t, result, returns.List{
			rtn(t, fixtures.Day2, 0.6),
			rtn(t, fixtures.Day1, 0.6),
		})
	})

	t.Run("truncate the last return", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day2, 0.6),
			rtn(t, fixtures.Day1, 0.6),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.Day1), fixtures.T(t, fixtures.Day0))
		assert.Equal(t, result, returns.List{
			rtn(t, fixtures.Day1, 0.6),
			rtn(t, fixtures.Day0, 0.6),
		})
	})

	t.Run("end during weekend", func(t *testing.T) {
		t.Run("on sunday", func(t *testing.T) {

			twoWeeksOfReturns := returns.List{
				{Time: fixtures.T(t, "2021-04-23")},
				{Time: fixtures.T(t, "2021-04-22")},
				{Time: fixtures.T(t, "2021-04-21")},
				{Time: fixtures.T(t, "2021-04-20")},
				{Time: fixtures.T(t, "2021-04-19")},
				{Time: fixtures.T(t, "2021-04-16")},
				{Time: fixtures.T(t, "2021-04-15")},
				{Time: fixtures.T(t, "2021-04-14")},
				{Time: fixtures.T(t, "2021-04-13")},
				{Time: fixtures.T(t, "2021-04-12")},
			}

			end := fixtures.T(t, "2021-04-18")
			assert.Equal(t, end.Weekday(), time.Sunday)
			start := fixtures.T(t, "2021-04-12")

			result := twoWeeksOfReturns.Between(end, start)

			assert.Equal(t, result, returns.List{
				{Time: fixtures.T(t, "2021-04-16")},
				{Time: fixtures.T(t, "2021-04-15")},
				{Time: fixtures.T(t, "2021-04-14")},
				{Time: fixtures.T(t, "2021-04-13")},
				{Time: fixtures.T(t, "2021-04-12")},
			})
		})

		t.Run("on saturday", func(t *testing.T) {

			twoWeeksOfReturns := returns.List{
				{Time: fixtures.T(t, "2021-04-23")},
				{Time: fixtures.T(t, "2021-04-22")},
				{Time: fixtures.T(t, "2021-04-21")},
				{Time: fixtures.T(t, "2021-04-20")},
				{Time: fixtures.T(t, "2021-04-19")},
				{Time: fixtures.T(t, "2021-04-16")},
				{Time: fixtures.T(t, "2021-04-15")},
				{Time: fixtures.T(t, "2021-04-14")},
				{Time: fixtures.T(t, "2021-04-13")},
				{Time: fixtures.T(t, "2021-04-12")},
			}

			end := fixtures.T(t, "2021-04-17")
			assert.Equal(t, end.Weekday(), time.Saturday)
			start := fixtures.T(t, "2021-04-12")

			result := twoWeeksOfReturns.Between(end, start)

			assert.Equal(t, result, returns.List{
				{Time: fixtures.T(t, "2021-04-16")},
				{Time: fixtures.T(t, "2021-04-15")},
				{Time: fixtures.T(t, "2021-04-14")},
				{Time: fixtures.T(t, "2021-04-13")},
				{Time: fixtures.T(t, "2021-04-12")},
			})
		})
	})
}

func rtn(t *testing.T, s string, v float64) returns.Return {
	t.Helper()
	return returns.New(fixtures.T(t, s), v)
}
