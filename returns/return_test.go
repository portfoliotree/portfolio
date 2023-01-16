package returns_test

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/portfoliotree/portfolio/internal/fixtures"
	"github.com/portfoliotree/portfolio/returns"
)

func TestNewReturn(t *testing.T) {
	d := fixtures.T(t, "2022-01-01")
	t.Run("okay", func(t *testing.T) {
		r := returns.New(d, 0.4)
		o := NewWithT(t)
		o.Expect(r.Time).To(Equal(d))
		o.Expect(r.Value).To(Equal(0.4))
	})
	t.Run("nan", func(t *testing.T) {
		o := NewWithT(t)
		o.Expect(func() {
			returns.New(d, math.NaN())
		}).To(Panic())
	})
}
