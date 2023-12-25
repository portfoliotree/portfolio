package backtestconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/portfolio/backtest/backtestconfig"
)

func TestDurations_Validate(t *testing.T) {
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
