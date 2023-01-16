package calculations_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/portfoliotree/portfolio/calculations"
)

func TestDiscreteReturns(t *testing.T) {
	please := NewWithT(t)
	please.Expect(calculations.HoldingPeriodReturns([]float64{100, 50})).To(Equal([]float64{1}))
	please.Expect(calculations.HoldingPeriodReturns([]float64{50, 50})).To(Equal([]float64{0}))
	please.Expect(calculations.HoldingPeriodReturns([]float64{50, 100})).To(Equal([]float64{-0.5}))
	please.Expect(calculations.HoldingPeriodReturns([]float64{50, 100, 100})).To(Equal([]float64{-0.5, 0}))
	please.Expect(calculations.HoldingPeriodReturns(nil)).To(HaveLen(0))
	please.Expect(calculations.HoldingPeriodReturns([]float64{50, 100, 100})).To(Equal([]float64{-0.5, 0}))
}
