package backtestconfig_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/portfolio/backtest/backtestconfig"
	"github.com/portfoliotree/portfolio/internal/fixtures"
	"github.com/portfoliotree/portfolio/returns"
)

func TestWindows_Validate(t *testing.T) {
	for _, d := range backtestconfig.Windows() {
		t.Run(d.String(), func(t *testing.T) {
			err := d.Validate()
			assert.NoError(t, err)
		})
	}

	t.Run("not set", func(t *testing.T) {
		err := backtestconfig.Window("").Validate()
		assert.NoError(t, err)
	})

	t.Run("an animal", func(t *testing.T) {
		err := backtestconfig.Window("Cat").Validate()
		assert.Error(t, err)
	})
}

func TestWindow_Function(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		var zero backtestconfig.Window

		today := fixtures.T(t, fixtures.Day2)
		table := returns.NewTable([]returns.List{{
			returns.New(fixtures.T(t, fixtures.Day3), .1),
			returns.New(today, .1),
			returns.New(fixtures.T(t, fixtures.Day1), .1),
			returns.New(fixtures.T(t, fixtures.Day0), .1),
		}})

		result := zero.Function(today, table)
		assert.Equal(t, result.FirstTime().Format(time.DateOnly), fixtures.Day0)
	})
}
