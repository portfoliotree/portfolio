package portfolio_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/portfolio"
)

func TestSpecification_AssetReturns(t *testing.T) {
	for _, tt := range []struct {
		Name string
		ctx  context.Context
		pf   portfolio.Specification

		ErrorStringContains string
	}{
		{
			Name:                "nil context",
			pf:                  portfolio.Specification{Assets: []portfolio.Component{{ID: "AAPL"}}},
			ErrorStringContains: "Context",
		},
		{
			Name: "no assets",
			pf:   portfolio.Specification{Assets: []portfolio.Component{}},
		},
		{
			Name: "nil assets",
			pf:   portfolio.Specification{Assets: nil},
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := portfolio.AssetReturnsTable(tt.ctx, tt.pf.Assets)
			if tt.ErrorStringContains == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.ErrorStringContains)
			}
		})
	}
}

func Test_Specification_AssetReturns_bad_URL(t *testing.T) {
	t.Setenv(portfolio.ServerURLEnvironmentVariableName, ":lemon:")
	pf := portfolio.Specification{Assets: []portfolio.Component{{ID: "AAPL"}}}
	_, err := portfolio.AssetReturnsTable(context.Background(), pf.Assets)
	assert.ErrorContains(t, err, "lemon")
}
