package returns

import (
	"fmt"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/portfolio/internal/fixtures"
)

func Test_betweenIndexes(t *testing.T) {
	t.Run("range around a single return", func(t *testing.T) {
		list := List{
			{
				Time: fixtures.T(t, "2020-07-15"),
			},
			{
				Time: fixtures.T(t, "2020-06-15"),
			},
			{
				Time: fixtures.T(t, "2020-05-15"),
			},
			{
				Time: fixtures.T(t, "2020-04-15"),
			},
			{
				Time: fixtures.T(t, "2020-03-15"),
			},
			{
				Time: fixtures.T(t, "2020-02-15"),
			},
			{
				Time: fixtures.T(t, "2020-01-15"),
			},
		}
		sort.Sort(list)

		for i, ret := range list[1:] {
			index := i + 1

			t.Run(strconv.Itoa(index), func(i int, r Return) func(t *testing.T) {
				return func(t *testing.T) {
					month := len(list) - index
					t1 := fixtures.T(t, fmt.Sprintf("2020-%02d-20", month+1))
					t0 := fixtures.T(t, fmt.Sprintf("2020-%02d-20", month))
					t.Logf("%s\t%s", t1.Format(time.DateOnly), t0.Format(time.DateOnly))

					i1, i0 := lowAndHighIndexesWithinTimes(list, t1, t0, func(r Return) time.Time {
						return r.Time
					})

					t.Logf("t1: %d, t0: %d", i1, i0)
					assert.Equal(t, i0-i1, 1)
					assert.Equal(t, i0, index)
					assert.Equal(t, i1, index-1)
				}
			}(index, ret))
		}
	})
}
