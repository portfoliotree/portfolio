package calculations_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/portfolio/calculations"
)

func TestDiscreteReturns(t *testing.T) {
	assert.Equal(t, calculations.HoldingPeriodReturns([]float64{100, 50}), []float64{1})
	assert.Equal(t, calculations.HoldingPeriodReturns([]float64{50, 50}), []float64{0})
	assert.Equal(t, calculations.HoldingPeriodReturns([]float64{50, 100}), []float64{-0.5})
	assert.Equal(t, calculations.HoldingPeriodReturns([]float64{50, 100, 100}), []float64{-0.5, 0})
	assert.Len(t, calculations.HoldingPeriodReturns(nil), 0)
	assert.Equal(t, calculations.HoldingPeriodReturns([]float64{50, 100, 100}), []float64{-0.5, 0})
}
