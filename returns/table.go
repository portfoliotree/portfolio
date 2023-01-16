package returns

import (
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/portfoliotree/round"
	"go.mongodb.org/mongo-driver/bson"
	"gonum.org/v1/gonum/mat"

	"github.com/portfoliotree/portfolio/calculations"
)

type Table struct {
	times  []time.Time
	values [][]float64
	// isRoot bool
}

func NewTable(list []List) Table {
	if len(list) == 0 {
		return Table{}
	}
	table := Table{
		// isRoot: true,
		values: make([][]float64, 0, len(list)),
	}
	for _, slice := range list {
		table = table.AddColumn(slice)
	}
	return table
}

func (table *Table) UnmarshalBSON(buf []byte) error {
	var enc encodedTable
	err := bson.Unmarshal(buf, &enc)
	table.times = enc.Times
	table.values = enc.Values
	return err
}

func (table Table) MarshalBSON() ([]byte, error) {
	return bson.Marshal(encodedTable{
		Times:  table.times,
		Values: table.values,
	})
}

type encodedTable struct {
	Times  []time.Time `json:"times" bson:"times"`
	Values [][]float64 `json:"values" bson:"values"`
}

func (table *Table) UnmarshalJSON(buf []byte) error {
	var enc encodedTable
	err := json.Unmarshal(buf, &enc)
	table.times = enc.Times
	table.values = enc.Values
	return err
}

func (table Table) MarshalJSON() ([]byte, error) {
	t := encodedTable{
		Times:  table.times,
		Values: table.values,
	}
	err := round.Recursive(t.Values, 6)
	if err != nil {
		return nil, err
	}
	return json.Marshal(t)
}

func (table Table) Join(other Table) Table {
	updated := table
	for _, slice := range other.Lists() {
		updated = updated.AddColumn(slice)
	}
	return updated
}

func (table Table) Times() []time.Time { return table.times }
func (table Table) List(columnIndex int) List {
	list := make(List, len(table.times))
	for i := range table.times {
		list[i].Time = table.times[i]
		list[i].Value = table.values[columnIndex][i]
	}
	return list
}

func (table Table) ColumnValues() [][]float64 { return table.values }
func (table Table) NumberOfColumns() int      { return len(table.values) }
func (table Table) NumberOfRows() int         { return len(table.times) }

func (table Table) Row(tm time.Time) ([]float64, bool) {
	if table.NumberOfRows() == 0 {
		return nil, false
	}
	index, found := sort.Find(len(table.times), func(i int) int {
		return compareTimes(table.times[i], tm)
	})
	result := make([]float64, len(table.values))
	if !found {
		return result, false
	}
	for i := range table.values {
		result[i] = table.values[i][index]
	}
	return result, true
}

func (table Table) HasRow(tm time.Time) bool {
	_, found := sort.Find(len(table.times), func(i int) int {
		return compareTimes(table.times[i], tm)
	})
	return found
}

func (table Table) AddColumn(list List) Table {
	//if !table.isRoot {
	//	panic("modifying a sliced Table is prohibited")
	//}
	sort.Sort(list)
	if len(table.values) == 0 {
		return table.addInitialColumn(list)
	}
	return table.addAdditionalColumn(list)
}

func (table Table) AddColumns(lists []List) Table {
	updated := table
	for _, list := range lists {
		updated = updated.AddColumn(list)
	}
	return updated
}

func (table Table) Between(last, first time.Time) Table {
	lastIdx, firstIdx := table.RangeIndexes(last, first)
	values := make([][]float64, len(table.values))
	for i := range table.values {
		values[i] = table.values[i][lastIdx:firstIdx:firstIdx]
	}
	return Table{
		times:  table.times[lastIdx:firstIdx:firstIdx],
		values: values,
	}
}

func (table Table) addInitialColumn(s List) Table {
	newValues := make([]float64, 0, maxInt(len(s), len(table.times)))
	for _, r := range s {
		table.times = append(table.times, r.Time)
		newValues = append(newValues, r.Value)
	}
	table.values = append(table.values, newValues)
	return table
}

func (table Table) addAdditionalColumn(list List) Table {
	list = list.Between(table.LastTime(), table.FirstTime())
	updated := table.Between(list.LastTime(), list.FirstTime())

	for _, r := range list {
		_, updated = updated.ensureRowForTime(r.Time)
	}

	newValues := make([]float64, len(updated.times))
	for i, tm := range updated.times {
		value, _ := list.Value(tm)
		newValues[i] = value
	}
	updated.values = append(updated.values, newValues)
	return updated
}

func (table Table) ensureRowForTime(tm time.Time) (index int, updated Table) {
	for i, et := range table.times {
		if et.Equal(tm) {
			return i, table
		}
		if tm.After(et) {
			index, updated = i, table
			updated.times = append(updated.times[:i], append([]time.Time{tm}, updated.times[i:]...)...)
			for j, values := range updated.values {
				updated.values[j] = append(values[:i], append([]float64{0}, values[i:]...)...)
			}
			break

			//// an early return makes the coverage dip below 100% because the
			//// empty block outside the loop would never execute. This break
			//// is essentially like the following line
			// return index, updated
		}
	}
	return index, updated
}

func (table Table) FirstTime() time.Time { return indexOrEmpty(table.times, firstIndex(table.times)) }
func (table Table) LastTime() time.Time  { return indexOrEmpty(table.times, lastIndex(table.times)) }
func (table Table) TimeAfter(tm time.Time) (time.Time, bool) {
	if tm.Before(table.FirstTime()) {
		return table.FirstTime(), true
	}
	if tm.After(table.LastTime()) {
		return time.Time{}, false
	}
	index := indexOfClosest(table.times, func(t time.Time) time.Time {
		return t
	}, tm)
	next := indexOrEmpty(table.times, index-1)
	return next, !next.IsZero()
}

func (table Table) TimeBefore(tm time.Time) (time.Time, bool) {
	if tm.After(table.LastTime()) {
		return table.LastTime(), true
	}
	if tm.Before(table.FirstTime()) {
		return time.Time{}, false
	}
	index := indexOfClosest(table.times, func(t time.Time) time.Time {
		return t
	}, tm)
	next := indexOrEmpty(table.times, index+1)
	return next, !next.IsZero()
}

func (table Table) Lists() []List {
	result := make([]List, len(table.values))
	for i := range table.values {
		result[i] = table.List(i)
	}
	return result
}

func (table Table) CorrelationMatrix() *mat.Dense {
	return calculations.CorrelationMatrix(table.values)
}

func (table Table) CorrelationMatrixValues() [][]float64 {
	dense := table.CorrelationMatrix()
	if dense == nil {
		return nil
	}
	iL, jL := dense.Dims()
	result := make([][]float64, iL)
	d := dense.RawMatrix().Data
	for i := range result {
		result[i] = d[i*jL : (i+1)*jL]
	}
	return result
}

func (table Table) ExpectedRisk(weights []float64) float64 {
	risks := table.Risks()
	r, _, _ := calculations.RiskFromRiskContribution(risks, weights, table.CorrelationMatrix())
	return r
}

func (table Table) Risks() []float64 {
	result := make([]float64, table.NumberOfColumns())
	for i, column := range table.values {
		result[i] = calculations.RiskFromStdDev(column)
	}
	return result
}

func (table Table) AnnualizedRisks() []float64 {
	result := table.Risks()
	for i, risk := range result {
		result[i] = calculations.AnnualizeRisk(risk, calculations.PeriodsPerYear)
	}
	return result
}

func (table Table) TimeWeightedReturns() []float64 {
	result := make([]float64, table.NumberOfColumns())
	for i, column := range table.values {
		result[i] = calculations.AnnualizedTimeWeightedReturn(column, calculations.PeriodsPerYear)
	}
	return result
}

func (table Table) AnnualizedArithmeticReturns() []float64 {
	result := make([]float64, table.NumberOfColumns())
	for i, column := range table.values {
		result[i] = calculations.AnnualizedArithmeticReturn(column)
	}
	return result
}

func (table Table) EndAndStartDates() (end, start time.Time, _ error) {
	if table.NumberOfColumns() == 0 {
		return time.Time{}, time.Time{}, ErrorNoReturns{}
	}
	if end.Before(start) {
		return time.Time{}, time.Time{}, errors.New("no overlap")
	}
	return table.LastTime(), table.FirstTime(), nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type ColumnGroup struct {
	index, length int
}

func (group ColumnGroup) Length() int { return group.length }

func (table Table) AddColumnGroup(lists []List) (ColumnGroup, Table) {
	updated := table
	startingColumnIndex := table.NumberOfColumns()
	for _, list := range lists {
		updated = updated.AddColumn(list)
	}
	return ColumnGroup{
		index:  startingColumnIndex,
		length: len(lists),
	}, updated
}

func (table Table) ColumnGroupColumnIndex(group ColumnGroup, groupIndex int) (columnIndex int) {
	columnIndex = group.index + groupIndex

	if columnIndex >= len(table.values) {
		panic("column index out of bounds")
	}

	return columnIndex
}

func (table Table) ColumnGroupAsTable(group ColumnGroup) Table {
	return Table{
		times:  table.times,
		values: table.values[group.index : group.index+group.length : group.index+group.length],
	}
}

func (table Table) ColumnGroupLists(group ColumnGroup) []List {
	result := make([]List, group.length)
	for i, list := range table.values[group.index : group.index+group.length] {
		result[i] = make(List, len(table.times))
		for j, v := range list {
			result[i][j].Value = v
			result[i][j].Time = table.times[j]
		}
	}
	return result
}

func (table Table) ColumnGroupValues(group ColumnGroup) [][]float64 {
	return table.values[group.index : group.index+group.length : group.index+group.length]
}

func (table Table) MostRecentValues() []float64 {
	result := make([]float64, len(table.values))
	for i, column := range table.values {
		result[i] = column[0]
	}
	return result
}

// AlignTables may be used to ensure multiple tables are date-aligned.
func AlignTables(tables ...Table) (_ []Table, end, start time.Time, _ error) {
	var (
		table  Table
		groups []ColumnGroup
	)
	for _, rl := range tables {
		var group ColumnGroup
		group, table = table.AddColumnGroup(rl.Lists())
		groups = append(groups, group)
	}
	result := make([]Table, len(tables))
	for i, g := range groups {
		result[i] = table.ColumnGroupAsTable(g)
	}
	return result, table.LastTime(), table.FirstTime(), nil
}

// RangeIndexes is a helper to be used to align additional non-return columns to a table
func (table Table) RangeIndexes(last, first time.Time) (end int, start int) {
	tmFn := func(t time.Time) time.Time { return t }
	lastIdx, firstIdx := lowAndHighIndexesWithinTimes(table.times, last, first, tmFn)
	return lastIdx, firstIdx
}
