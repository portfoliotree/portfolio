package calculations

import (
	"fmt"
	"math"
	"testing"

	"github.com/portfoliotree/round"
	"github.com/stretchr/testify/assert"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

func TestRisk(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		risk, err := Risk(
			[]float64{10, 10},
			[]float64{100, 0.0},
			[][]float64{
				{1, 1},
				{1, 1},
			},
		)
		if err != nil {
			t.Errorf("it should not return an error: %s", err)
		}

		if math.IsNaN(risk) {
			t.Errorf("risk should be a number got: %f", risk)
		}
	})
}

func ExampleRiskFromStdDev() {
	risk := RiskFromStdDev([]float64{-0.1, 0.1, -0.1, 0.1, -0.1, 0.1, -0.1, 0.1})
	fmt.Printf("%.2f", risk)
	// Output: 0.11
}

func TestClosingReturns(t *testing.T) {
	// test return of one dollar per day

	// t.Run("one dollar per day", func(t *testing.T) {
	// 	please := NewGomegaWithT(t)
	//
	// 	quotes := []alphavantage.Quote{
	// 		{Time: parseTime("2020-08-06"), Close: 24000},
	// 		{Time: parseTime("2020-08-05"), Close: 16000},
	// 		{Time: parseTime("2020-08-04"), Close: 8000},
	// 		{Time: parseTime("2020-08-03"), Close: 8000},
	// 		{Time: parseTime("2020-08-02"), Close: 2000},
	// 		{Time: parseTime("2020-08-01"), Close: 200},
	// 	}
	//
	// 	returns := ClosingReturns(quotes)
	//
	// 	please.Expect(returns).To(HaveLen(5))
	// 	please.Expect(returns[0]).To(Equal(0.5))
	// 	please.Expect(returns[1]).To(Equal(1.0))
	// 	please.Expect(returns[2]).To(Equal(0.0))
	// 	please.Expect(returns[3]).To(Equal(3.0))
	// 	please.Expect(returns[4]).To(Equal(9.0))
	// })
}

/*
2020-08-03,1002,1003,1002,1003,1
2020-08-02,1001,1002,1001,1002,1
2020-08-01,1000,1001,1000,1001,1
*/

func TestNumberOfBets(t *testing.T) {
	a := 6.0
	b := 3.0

	betsval, err := NumberOfBets(a, b)
	if err != nil {
		t.Errorf("it should not return an error, got: %s", err)
	}

	if betsval != 4.0 {
		t.Errorf("%v should have equaled 4", betsval)
	}
}

func TestNumberOfBets_Perfect_Correlation(t *testing.T) {
	a := 6.0
	b := 6.0

	betsval, err := NumberOfBets(a, b)
	if err != nil {
		t.Errorf("it should not return an error, got: %s", err)
	}

	if betsval != 1.0 {
		t.Errorf("%v should have equaled 1", betsval)
	}
}

func TestNumberOfBets_Zero_Portfolio_Risk(t *testing.T) {
	a := 6.0

	_, err := NumberOfBets(a, 0)

	if err == nil {
		t.Errorf("it should return an error")
	}
}

func TestRiskFromRiskContribution(t *testing.T) {
	ws := []float64{-.1, .8, .3}
	rs := []float64{.2, .05, .25}
	cs := mat.NewDense(3, 3, []float64{
		1.00, 0.15, 0.60,
		0.15, 1.00, 0.05,
		0.60, 0.05, 1.00,
	})

	const precision = 4
	tr, rcs, rw := RiskFromRiskContribution(rs, ws, cs)
	tr = round.Decimal(tr, precision)
	_ = round.Recursive(rcs, precision)
	_ = round.Recursive(rw, precision)
	sum := round.Decimal(floats.Sum(rw), precision)

	assert.Equal(t, rcs, []float64{-0.0006, 0.0016, 0.0049})
	assert.Equal(t, tr, .0767)
	assert.Equal(t, rw, []float64{-0.1054, 0.2770, 0.8284})
	assert.Equal(t, sum, 1.0)
}

// func TestLegacyRiskFromRiskContribution(t *testing.T) {
// 	please := NewWithT(t)
//
// 	ws := []float64{-.1, .8, .3}
// 	rs := []float64{.2, .05, .25}
// 	cs := [][]float64{
// 		{1.00, 0.15, 0.60},
// 		{0.15, 1.00, 0.05},
// 		{0.60, 0.05, 1.00},
// 	}
//
// 	tr, rcs, rw := legacyRiskFromRiskContribution(rs, ws, cs)
//
// 	please.Expect(rcs).To(floattest.EqualSlice(4, []float64{-0.0006, 0.0016, 0.0049}))
// 	please.Expect(tr).To(floattest.Equal(4, .076714))
// 	please.Expect(rw).To(floattest.EqualSlice(4, []float64{-0.10535, 0.27698, 0.82838}))
// }
