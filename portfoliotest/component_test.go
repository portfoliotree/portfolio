package portfoliotest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/portfolio"
)

func TestComponentReturnsProvider(t *testing.T) {
	crp := ComponentReturnsProvider()

	t.Run("List", func(t *testing.T) {
		t.Run("returns not in testdata", func(t *testing.T) {
			_, err := crp.ComponentReturnsList(context.Background(), portfolio.Component{
				ID: "BANANA",
			})
			assert.Error(t, err)
		})
		t.Run("returns not in testdata", func(t *testing.T) {
			list, err := crp.ComponentReturnsList(context.Background(), portfolio.Component{
				ID: "GOOG",
			})
			assert.NoError(t, err)
			assert.NotZero(t, list)
		})
	})
	t.Run("Table", func(t *testing.T) {
		t.Run("returns not in testdata", func(t *testing.T) {
			_, err := crp.ComponentReturnsTable(context.Background(), portfolio.Component{
				ID: "BANANA",
			})
			assert.Error(t, err)
		})
		t.Run("returns not in testdata", func(t *testing.T) {
			tab, err := crp.ComponentReturnsTable(context.Background(), portfolio.Component{
				ID: "GOOG",
			}, portfolio.Component{
				ID: "AAPL",
			})
			assert.NoError(t, err)
			assert.Equal(t, 2, tab.NumberOfColumns())
		})
	})
}
