package returns

import (
	"sort"
	"time"

	"github.com/portfoliotree/portfolio/calculations"
)

type List []Return

func (list List) Returns() List      { return list }
func (list List) Sort()              { sort.Sort(list) }
func (list List) Less(i, j int) bool { return list[i].Time.After(list[j].Time) }
func (list List) Len() int           { return len(list) }
func (list List) Swap(i, j int)      { list[i], list[j] = list[j], list[i] }

func (list List) First() Return        { return indexOrEmpty(list, firstIndex(list)) }
func (list List) Last() Return         { return indexOrEmpty(list, lastIndex(list)) }
func (list List) FirstTime() time.Time { return indexOrEmpty(list, firstIndex(list)).Time }
func (list List) LastTime() time.Time  { return indexOrEmpty(list, lastIndex(list)).Time }

func (list List) Values() []float64 {
	result := make([]float64, len(list))
	for i := range list {
		result[i] = list[i].Value
	}
	return result
}

func (list List) Times() []time.Time {
	result := make([]time.Time, len(list))
	for i := range list {
		result[i] = list[i].Time
	}
	return result
}

// Between returns a slice of a list. You may want pass result into slices.Clone()
// before using any mutating functions (such as Insert).
func (list List) Between(t1, t0 time.Time) List {
	tmFn := func(r Return) time.Time { return r.Time }
	last, first := lowAndHighIndexesWithinTimes(list, t1, t0, tmFn)
	return list[last:first:first]
}

// Insert correctly places newReturn in List.
// If newReturn.Time is the same as an existing return in s, the value of the existing return is overwritten.
func (list List) Insert(newReturn Return) List {
	for i, er := range list {
		if newReturn.Time.Equal(er.Time) {
			list[i].Value = newReturn.Value
			return list
		}
		if !newReturn.Time.After(er.Time) {
			continue
		}
		return append(list[:i], append(List{newReturn}, list[i:]...)...)
	}
	return append(list, newReturn)
}

func (list List) Value(t time.Time) (float64, bool) {
	index, found := sort.Find(len(list), func(i int) int {
		return compareTimes(list[i].Time, t)
	})
	if !found {
		return 0, found
	}
	return list[index].Value, found
}

func (list List) Excess(other List) List {
	table := NewTable([]List{list, other})
	result := make(List, len(table.times))
	for i, t := range table.times {
		result[i].Time = t
		result[i].Value = table.values[0][i] - table.values[1][i]
	}
	return result
}

func (list List) TimeWeightedReturn() float64 {
	return calculations.TimeWeightedReturn(list.Values())
}

func (list List) EndAndStartDate() (end, start time.Time, _ error) {
	if list.Len() == 0 {
		return time.Time{}, time.Time{}, ErrorNoReturns{}
	}
	return list.LastTime(), list.FirstTime(), nil
}

// Risk calls calculations.RiskFromStdDev
func (list List) Risk() float64 {
	return calculations.RiskFromStdDev(list.Values())
}

// AnnualizedRisk must receive at least 2 returns otherwise it returns 0
func (list List) AnnualizedRisk() float64 {
	return calculations.AnnualizeRisk(list.Risk(), calculations.PeriodsPerYear)
}

// AnnualizedTimeWeightedReturn must receive at least 2 returns otherwise it returns 0
func (list List) AnnualizedTimeWeightedReturn() float64 {
	return calculations.AnnualizedTimeWeightedReturn(list.Values(), calculations.PeriodsPerYear)
}

func (list List) AnnualizedArithmeticReturn() float64 {
	return calculations.AnnualizedArithmeticReturn(list.Values())
}

func compareTimes(x, y time.Time) int {
	if x.Equal(y) {
		return 0
	}
	if x.Before(y) {
		return -1
	}
	return 1
}

func indexOrEmpty[T any](list []T, i int) T {
	if i < 0 || i >= len(list) {
		var zero T
		return zero
	}
	return list[i]
}

func firstIndex[T any](list []T) int {
	return len(list) - 1
}

func lastIndex[T any](_ []T) int {
	return 0
}

func indexOfClosest[T any](list []T, time func(T) time.Time, t time.Time) int {
	if len(list) == 0 {
		return 0
	}
	index := len(list) / 2

	if t.Equal(time(list[index])) {
		return index
	}

	//// Not required by tests. I am keeping it until I run the full test suite
	//if len(list) == 2 && t.Before(time(list[0])) && (t.After(time(list[1])) || t.Equal(time(list[1]))) {
	//	return 1
	//}

	if t.Before(time(list[index])) {
		if len(list) == 1 {
			return 1
		}
		return index + indexOfClosest(list[index:], time, t)
	}
	return indexOfClosest(list[:index], time, t)
}

func inBounds[T any](list []T, index int) bool {
	return index > 0 && index < len(list)
}

// lowAndHighIndexesWithinTimes returns indexFinal and indexInitial where
// list[indexFinal].Time is before or equal to t1 and
// list[indexInitial].Time is after or equal to t0
// It returns zeros when the slice is empty or when t1 is after t0.
func lowAndHighIndexesWithinTimes[T any](list []T, t1, t0 time.Time, time func(T) time.Time) (int, int) {
	if len(list) == 1 && t0.Equal(time(list[0])) && t1.Equal(time(list[0])) {
		return 0, 1
	}
	if len(list) == 0 ||
		t1.Before(t0) ||
		t1.Before(time(list[firstIndex(list)])) ||
		t0.After(time(list[lastIndex(list)])) {
		return 0, 0
	}
	indexFinal := indexOfClosest(list, time, t1)
	indexInitial := indexOfClosest(list, time, t0)
	if inBounds(list, indexInitial) && time(list[indexInitial]).Before(t0) {
		indexInitial--
	}
	return indexFinal, minInt(indexInitial+1, len(list))
}
