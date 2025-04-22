package calculate_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/portfolio/calculate"
)

func TestDiscreteReturns(t *testing.T) {
	assert.Equal(t, calculate.HoldingPeriodReturns([]float64{100, 50}), []float64{1})
	assert.Equal(t, calculate.HoldingPeriodReturns([]float64{50, 50}), []float64{0})
	assert.Equal(t, calculate.HoldingPeriodReturns([]float64{50, 100}), []float64{-0.5})
	assert.Equal(t, calculate.HoldingPeriodReturns([]float64{50, 100, 100}), []float64{-0.5, 0})
	assert.Len(t, calculate.HoldingPeriodReturns(nil), 0)
	assert.Equal(t, calculate.HoldingPeriodReturns([]float64{50, 100, 100}), []float64{-0.5, 0})
}

func TestCompoundReturn(t *testing.T) {
	t.Run("when vol is zero", func(t *testing.T) {
		result := calculate.CAGRFromArithmeticReturn(0.01230, 0)
		assert.InDelta(t, 0.01230, result, 0.00001)
	})
}

// TODO: add shorthand test to see where this would work
// shorthand := arithmeticReturn - ((volatility * volatility) / 2)

func FuzzCompoundReturn(f *testing.F) {
	const minVol, maxVol = 0.0, 2000 / 100
	f.Add(0.01, 0.0)
	f.Add(0.01, 0.01)
	f.Fuzz(func(t *testing.T, arithmeticReturn, v float64) {
		v = math.Abs(v)
		volatility := minVol + v*(maxVol-minVol)

		geometricReturn := calculate.CAGRFromArithmeticReturn(arithmeticReturn, volatility)
		arithmeticReturnResult := calculate.ArithmeticReturnFromCAGR(geometricReturn, volatility)

		assert.False(t, math.IsNaN(geometricReturn))
		assert.False(t, math.IsInf(geometricReturn, -1))
		assert.False(t, math.IsInf(geometricReturn, 1))
		assert.False(t, math.IsNaN(arithmeticReturnResult))
		assert.False(t, math.IsInf(arithmeticReturnResult, -1))
		assert.False(t, math.IsInf(arithmeticReturnResult, 1))

		if arithmeticReturn > calculate.MinArithmeticReturn {
			assert.InDelta(t, arithmeticReturn, arithmeticReturnResult, 0.000001, "arithmeticReturn: %v, volatility: %v", arithmeticReturn, volatility)
		}
	})
}
