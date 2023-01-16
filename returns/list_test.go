package returns_test

import (
	"sort"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/portfoliotree/round"

	"github.com/portfoliotree/portfolio/internal/fixtures"
	"github.com/portfoliotree/portfolio/returns"
)

// var _ portfolio.Asset = returns.List{}

func TestList_Returns(t *testing.T) {
	list := returns.List{
		returns.New(fixtures.T(t, fixtures.Day0), 400),
	}

	o := NewWithT(t)
	o.Expect(list.Returns()).To(Equal(returns.List{
		returns.New(fixtures.T(t, fixtures.Day0), 400),
	}))
}

func TestList_Value(t *testing.T) {
	d := fixtures.T(t, "2021-12-31")

	t.Run("empty", func(t *testing.T) {
		o := NewWithT(t)
		var slice returns.List
		slice.Sort()
		updated, found := slice.Value(d)
		o.Expect(found).To(BeFalse())
		o.Expect(updated).To(BeZero())
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
			o := NewWithT(t)
			o.Expect(found).To(BeFalse())
			o.Expect(updated).To(Equal(0.0))
		})
	})

	t.Run("found", func(t *testing.T) {
		t.Run("exact match one return", func(t *testing.T) {
			o := NewWithT(t)
			var (
				existingReturn = returns.New(d, .024)
				slice          = returns.List{
					existingReturn,
				}
			)
			slice.Sort()

			updated, found := slice.Value(d)

			o.Expect(found).To(BeTrue())
			o.Expect(updated).To(Equal(.024))
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
			o := NewWithT(t)
			o.Expect(found).To(BeTrue())
			o.Expect(updated).To(Equal(.666))
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
			o := NewWithT(t)
			o.Expect(found).To(BeTrue())
			o.Expect(updated).To(Equal(.666))
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
			o := NewWithT(t)
			o.Expect(found).To(BeTrue())
			o.Expect(updated).To(Equal(.666))
		})
	})
}

func TestList_Values(t *testing.T) {
	rs := returns.List{
		returns.New(fixtures.T(t, "2022-03-01"), 0.3),
		returns.New(fixtures.T(t, "2022-02-01"), 0.2),
		returns.New(fixtures.T(t, "2022-01-01"), 0.1),
	}
	o := NewWithT(t)
	o.Expect(rs.Values()).To(Equal([]float64{0.3, 0.2, 0.1}))
}

func TestList_Sort(t *testing.T) {
	d := fixtures.T(t, "2022-01-01")
	t.Run("out of order", func(t *testing.T) {
		// given
		d0, d1, d2 := d.AddDate(-2, 0, 0), d.AddDate(-1, 0, 0), d
		slice := returns.List{
			returns.New(d1, 0.1),
			returns.New(d0, 0.1),
			returns.New(d2, 0.1),
		}

		// when
		slice.Sort()

		// then
		o := NewWithT(t)
		o.Expect(slice).To(Equal(returns.List{
			returns.New(d2, 0.1),
			returns.New(d1, 0.1),
			returns.New(d0, 0.1),
		}))
	})
}

func TestList_Insert(t *testing.T) {
	d := fixtures.T(t, "2021-12-31")

	t.Run("empty", func(t *testing.T) {
		o := NewWithT(t)
		var (
			newReturn = returns.New(d, .001)
			slice     returns.List
		)

		updated := slice.Insert(newReturn)

		o.Expect(updated).To(HaveLen(1))
		o.Expect(updated).To(ContainElement(newReturn))
		o.Expect(sort.IsSorted(updated)).To(BeTrue())
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
		o := NewWithT(t)
		o.Expect(updated).To(HaveLen(2))
		o.Expect(updated).To(Equal(returns.List{
			existingReturn,
			newReturn,
		}))
		o.Expect(sort.IsSorted(updated)).To(BeTrue())
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
		o := NewWithT(t)
		o.Expect(updated).To(HaveLen(2))
		o.Expect(updated).To(Equal(returns.List{
			newReturn,
			existingReturn,
		}))
		o.Expect(sort.IsSorted(updated)).To(BeTrue())
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
		o := NewWithT(t)
		o.Expect(updated).To(HaveLen(3))
		o.Expect(updated).To(Equal(returns.List{
			existingReturnNewer, newReturn, existingReturnOlder,
		}))
		o.Expect(sort.IsSorted(updated)).To(BeTrue())
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
		o := NewWithT(t)
		o.Expect(updated).To(HaveLen(1))
		o.Expect(updated).To(Equal(returns.List{
			newReturn,
		}))
		o.Expect(sort.IsSorted(updated)).To(BeTrue())
	})
}

func TestList_First(t *testing.T) {
	d := fixtures.T(t, "2021-12-31")

	t.Run("empty", func(t *testing.T) {
		var list returns.List
		o := NewWithT(t)
		o.Expect(func() {
			list.First()
		}).NotTo(Panic())
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
		o := NewWithT(t)
		o.Expect(rt).To(Equal(d))
		o.Expect(r).To(Equal(returns.New(d, .1)))
	})
}

func TestList_Last(t *testing.T) {
	d := fixtures.T(t, "2021-12-31")

	t.Run("empty", func(t *testing.T) {
		var list returns.List
		o := NewWithT(t)
		o.Expect(func() {
			list.Last()
		}).NotTo(Panic())
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
		o := NewWithT(t)
		o.Expect(r).To(Equal(returns.New(d, .2)))
		o.Expect(rt).To(Equal(d))
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
	o := NewWithT(t)
	_ = round.Recursive(&result, 1)
	o.Expect(result).To(Equal(returns.List{
		rtn(t, fixtures.Day2, 0.1),
		rtn(t, fixtures.Day1, 0.0),
	}))
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
		o := NewWithT(t)
		o.Expect(result).To(Equal(list))
	})

	t.Run("both before", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day3, 0.4),
			rtn(t, fixtures.Day2, 0.2),
			rtn(t, fixtures.Day1, 0.1),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.DayBefore), fixtures.T(t, fixtures.DayBefore).AddDate(-1, 0, 0))
		o := NewWithT(t)
		o.Expect(result).To(HaveLen(0))
	})

	t.Run("both after", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day3, 0.4),
			rtn(t, fixtures.Day2, 0.2),
			rtn(t, fixtures.Day1, 0.1),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.DayAfter).AddDate(1, 0, 0), fixtures.T(t, fixtures.DayAfter))
		o := NewWithT(t)
		o.Expect(result).To(HaveLen(0))
	})

	t.Run("same day", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day3, 0.4),
			rtn(t, fixtures.Day2, 0.2),
			rtn(t, fixtures.Day1, 0.1),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.Day1), fixtures.T(t, fixtures.Day1))
		o := NewWithT(t)
		o.Expect(result).To(Equal(returns.List{
			rtn(t, fixtures.Day1, 0.1),
		}))
	})

	t.Run("days between", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day3, 0.4),
			rtn(t, fixtures.Day2, 0.2),
			rtn(t, fixtures.Day1, 0.1),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.Day3), fixtures.T(t, fixtures.Day1))
		o := NewWithT(t)
		o.Expect(result).To(Equal(returns.List{
			rtn(t, fixtures.Day3, 0.4),
			rtn(t, fixtures.Day2, 0.2),
			rtn(t, fixtures.Day1, 0.1),
		}))
	})

	t.Run("a single return", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.Day0), fixtures.T(t, fixtures.Day0))
		o := NewWithT(t)
		o.Expect(result).To(Equal(returns.List{
			rtn(t, fixtures.Day0, 0.6),
		}))
	})

	t.Run("two returns", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day1, 0.6),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.Day1), fixtures.T(t, fixtures.Day0))
		o := NewWithT(t)
		o.Expect(result).To(Equal(returns.List{
			rtn(t, fixtures.Day1, 0.6),
			rtn(t, fixtures.Day0, 0.6),
		}))
	})

	t.Run("two returns", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day2, 0.6),
			rtn(t, fixtures.Day1, 0.6),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.Day2), fixtures.T(t, fixtures.Day0))
		o := NewWithT(t)
		o.Expect(result).To(Equal(returns.List{
			rtn(t, fixtures.Day2, 0.6),
			rtn(t, fixtures.Day1, 0.6),
			rtn(t, fixtures.Day0, 0.6),
		}))
	})

	t.Run("truncate the first return", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day2, 0.6),
			rtn(t, fixtures.Day1, 0.6),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.Day2), fixtures.T(t, fixtures.Day1))
		o := NewWithT(t)
		o.Expect(result).To(Equal(returns.List{
			rtn(t, fixtures.Day2, 0.6),
			rtn(t, fixtures.Day1, 0.6),
		}))
	})

	t.Run("truncate the last return", func(t *testing.T) {
		list := returns.List{
			rtn(t, fixtures.Day2, 0.6),
			rtn(t, fixtures.Day1, 0.6),
			rtn(t, fixtures.Day0, 0.6),
		}
		result := list.Between(fixtures.T(t, fixtures.Day1), fixtures.T(t, fixtures.Day0))
		o := NewWithT(t)
		o.Expect(result).To(Equal(returns.List{
			rtn(t, fixtures.Day1, 0.6),
			rtn(t, fixtures.Day0, 0.6),
		}))
	})

	t.Run("start during weekend", func(t *testing.T) {
		//o := NewWithT(t)
		//
		//twoWeeksOfReturns := returns.List{
		//	{Time: fixtures.T(t, "2021-04-23")},
		//	{Time: fixtures.T(t, "2021-04-22")},
		//	{Time: fixtures.T(t, "2021-04-21")},
		//	{Time: fixtures.T(t, "2021-04-20")},
		//	{Time: fixtures.T(t, "2021-04-19")},
		//	{Time: fixtures.T(t, "2021-04-16")},
		//	{Time: fixtures.T(t, "2021-04-15")},
		//	{Time: fixtures.T(t, "2021-04-14")},
		//	{Time: fixtures.T(t, "2021-04-13")},
		//	{Time: fixtures.T(t, "2021-04-12")},
		//}
		//
		//end := twoWeeksOfReturns.LastTime()
		//start := backtest.DurationWeek.SubtractFrom(twoWeeksOfReturns.LastTime())
		//
		//result := twoWeeksOfReturns.Between(end, start)
		//
		//o.Expect(result).To(Equal(returns.List{
		//	{Time: fixtures.T(t, "2021-04-23")},
		//	{Time: fixtures.T(t, "2021-04-22")},
		//	{Time: fixtures.T(t, "2021-04-21")},
		//	{Time: fixtures.T(t, "2021-04-20")},
		//	{Time: fixtures.T(t, "2021-04-19")},
		//}))
	})

	t.Run("end during weekend", func(t *testing.T) {
		t.Run("on sunday", func(t *testing.T) {
			o := NewWithT(t)

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
			o.Expect(end.Weekday()).To(Equal(time.Sunday))
			start := fixtures.T(t, "2021-04-12")

			result := twoWeeksOfReturns.Between(end, start)

			o.Expect(result).To(Equal(returns.List{
				{Time: fixtures.T(t, "2021-04-16")},
				{Time: fixtures.T(t, "2021-04-15")},
				{Time: fixtures.T(t, "2021-04-14")},
				{Time: fixtures.T(t, "2021-04-13")},
				{Time: fixtures.T(t, "2021-04-12")},
			}))
		})

		t.Run("on saturday", func(t *testing.T) {
			o := NewWithT(t)

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
			o.Expect(end.Weekday()).To(Equal(time.Saturday))
			start := fixtures.T(t, "2021-04-12")

			result := twoWeeksOfReturns.Between(end, start)

			o.Expect(result).To(Equal(returns.List{
				{Time: fixtures.T(t, "2021-04-16")},
				{Time: fixtures.T(t, "2021-04-15")},
				{Time: fixtures.T(t, "2021-04-14")},
				{Time: fixtures.T(t, "2021-04-13")},
				{Time: fixtures.T(t, "2021-04-12")},
			}))
		})
	})
}

func rtn(t *testing.T, s string, v float64) returns.Return {
	t.Helper()
	return returns.New(fixtures.T(t, s), v)
}
