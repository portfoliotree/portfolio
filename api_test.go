package portfolio_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/portfolio"
	"github.com/portfoliotree/portfolio/portfoliotest"
)

func TestMain(m *testing.M) {
	server := httptest.NewServer(testdataAssetReturns(portfoliotest.ComponentReturnsProvider()))
	http.DefaultClient = server.Client()
	_ = os.Setenv(portfolio.ServerURLEnvironmentVariableName, server.URL)
	os.Exit(func() int {
		defer server.Close()
		return m.Run()
	}())
}

func testdataAssetReturns(crp portfolio.ComponentReturnsProvider) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		assets, err := portfolio.ParseComponentsFromURL(req.URL.Query(), "asset")
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		table, err := crp.ComponentReturnsTable(req.Context(), assets...)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSONResponse(res, table)
	}
}

func writeJSONResponse(res http.ResponseWriter, data any) {
	buf, err := json.Marshal(data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(buf)
}

func Test_APIEndpoints(t *testing.T) {
	if value, found := os.LookupEnv("CI"); !found || value != "true" {
		t.Skip("Skipping test in CI environment")
	}

	t.Run("returns", func(t *testing.T) {
		pf := portfolio.Specification{
			Assets: []portfolio.Component{
				{ID: "AAPL"},
				{ID: "GOOG"},
			},
		}
		table, err := pf.AssetReturns(context.Background())
		assert.NoError(t, err)
		if table.NumberOfColumns() != 2 {
			t.Errorf("Expected 2 columns, got %d", table.NumberOfColumns())
		}
		if table.NumberOfRows() < 10 {
			t.Errorf("Expected at least 10 rows, got %d", table.NumberOfRows())
		}
	})
}

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
			_, err := tt.pf.AssetReturns(tt.ctx)
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
	_, err := pf.AssetReturns(context.Background())
	assert.ErrorContains(t, err, "lemon")
}
