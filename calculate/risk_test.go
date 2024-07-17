package calculate

import (
	"fmt"
	"testing"

	"github.com/portfoliotree/round"
	"github.com/stretchr/testify/assert"
)

func TestRisk(t *testing.T) {

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

func TestPortfolioVolatility(t *testing.T) {
	for _, tt := range []struct {
		Name         string
		Weights      []float64
		Risks        []float64
		Correlations [][]float64

		ExpectedPortfolioVolatility float64
		ExpectRiskContributions     []float64
	}{
		{
			Name:    "legacy example",
			Weights: []float64{-.1, .8, .3},
			Risks:   []float64{.2, .05, .25},
			Correlations: [][]float64{
				{1.00, 0.15, 0.60},
				{0.15, 1.00, 0.05},
				{0.60, 0.05, 1.00},
			},
			ExpectRiskContributions:     []float64{-0.0081, 0.0212, 0.0635},
			ExpectedPortfolioVolatility: .0767,
		},
		{
			Name:    "two",
			Weights: []float64{0.5, 0.5},
			Risks:   []float64{0.1668, 0.0428},
			Correlations: [][]float64{
				{1.00, 0.25},
				{0.25, 1.00},
			},
			ExpectRiskContributions:     []float64{0.0812, 0.0099},
			ExpectedPortfolioVolatility: .0911,
		},
		{
			Name:    "three",
			Weights: []float64{0.2500, 0.5000, 0.2500},
			Risks:   []float64{0.0428, 0.1668, 0.1693},
			Correlations: [][]float64{
				{1.0000, 0.2465, 0.4077},
				{0.2465, 1.0000, 0.1091},
				{0.4077, 0.1091, 1.0000},
			},
			ExpectRiskContributions:     []float64{0.0051, 0.0740, 0.0231},
			ExpectedPortfolioVolatility: .1022,
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			portfolioVol, riskContributions := PortfolioVolatility(tt.Weights, tt.Risks, tt.Correlations)

			const precision = 4
			portfolioVol = round.Decimal(portfolioVol, precision)
			_ = round.Recursive(riskContributions, precision)

			assert.Equal(t, tt.ExpectRiskContributions, riskContributions)
			assert.Equal(t, tt.ExpectedPortfolioVolatility, portfolioVol)
		})
	}
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
