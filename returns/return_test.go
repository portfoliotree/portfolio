package returns_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/portfolio/internal/fixtures"
	"github.com/portfoliotree/portfolio/returns"
)

func TestNewReturn(t *testing.T) {
	d := fixtures.T(t, "2022-01-01")
	t.Run("okay", func(t *testing.T) {
		r := returns.New(d, 0.4)
		assert.Equal(t, r.Time, d)
		assert.Equal(t, r.Value, 0.4)
	})
	t.Run("nan", func(t *testing.T) {
		assert.Panics(t, func() {
			returns.New(d, math.NaN())
		})
	})
}
