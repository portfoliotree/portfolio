package fixtures

import (
	"testing"
	"time"
)

const dateLayout = time.DateOnly

const (
	Day0      = "2022-10-20" // Thursday
	Day1      = "2022-10-21" // Friday
	Day2      = "2022-10-24" // Monday
	Day3      = "2022-10-25" // Tuesday
	FirstDay  = Day0
	LastDay   = Day3
	DayBefore = "2022-10-19"
	DayAfter  = "2022-10-26"
)

func T(t *testing.T, s string) time.Time {
	t.Helper()
	tm, err := time.Parse(dateLayout, s)
	if err != nil {
		t.Fatal(err)
	}
	return tm
}

// EveryFriday sets up to ten things on
func EveryFriday[V any](t *testing.T, list []V, set func(v V, t time.Time) V) []V {
	fridays := []time.Time{
		T(t, "2022-10-28"),
		T(t, "2022-10-21"),
		T(t, "2022-10-14"),
		T(t, "2022-10-07"),
		T(t, "2022-09-30"),
		T(t, "2022-09-23"),
		T(t, "2022-09-16"),
		T(t, "2022-09-09"),
		T(t, "2022-08-26"),
		T(t, "2022-09-02"),
	}
	for i := range list[:min(len(list), len(fridays))] {
		list[i] = set(list[i], fridays[i])
	}
	return list
}
