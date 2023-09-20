package portfolio_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/portfoliotree/portfolio"
)

func TestParseSpecificationFile(t *testing.T) {
	tmp := t.TempDir()
	badPortfolioSpecFilepath := filepath.Join(tmp, "invalid_portfolio.yml")
	// language=yaml
	require.NoError(t, os.WriteFile(badPortfolioSpecFilepath, []byte(`type: Banana`), 0o666))

	for _, tt := range []struct {
		Name                string
		FilePath            string
		ErrorStringContains string
		Documents           []portfolio.Document
	}{
		{
			Name:     "asset ids with policy weights",
			FilePath: filepath.Join("examples", "60-40_portfolio.yml"),
			Documents: []portfolio.Document{
				{
					Type: "Portfolio",
					Metadata: portfolio.Metadata{
						Name:      "60/40",
						Benchmark: portfolio.Component{ID: "BIGPX"},
					},
					Spec: portfolio.Specification{
						Assets: []portfolio.Component{
							{ID: "ACWI"},
							{ID: "AGG"},
						},
						Policy: portfolio.Policy{
							Weights:             []float64{60, 40},
							WeightsAlgorithm:    portfolio.PolicyAlgorithmConstantWeights,
							RebalancingInterval: "Quarterly",
						},
					},
					Filepath: "examples/60-40_portfolio.yml",
				},
			},
		},
		{
			Name:     "mixed asset spec node type and weight algorithm",
			FilePath: filepath.Join("examples", "maang_portfolio.yml"),
			Documents: []portfolio.Document{
				{
					Type: "Portfolio",
					Metadata: portfolio.Metadata{
						Name:      "MAANG",
						Benchmark: portfolio.Component{ID: "SPY"},
					},
					Spec: portfolio.Specification{
						Assets: []portfolio.Component{
							{ID: "META"},
							{ID: "AMZN"},
							{ID: "AAPL"},
							{ID: "NFLX"},
							{ID: "GOOG"},
						},
						Policy: portfolio.Policy{
							RebalancingInterval: "Quarterly",
							WeightsAlgorithm:    portfolio.PolicyAlgorithmEqualWeights,
						},
					},
					Filepath: "examples/maang_portfolio.yml",
				},
			},
		},
		{
			Name:                "file does not exist",
			FilePath:            "missing_portfolio.yml",
			ErrorStringContains: "file",
		},
		{
			Name:                "not a yaml file",
			FilePath:            "lemon.png",
			ErrorStringContains: "it must have a _portfolio.yml file name suffix",
		},
		{
			Name:                "not a valid portfolio specification",
			FilePath:            badPortfolioSpecFilepath,
			ErrorStringContains: "incorrect specification type",
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			documents, err := portfolio.ParseDocumentFile(tt.FilePath)
			if tt.ErrorStringContains == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.ErrorStringContains)
			}
			assert.Equal(t, tt.Documents, documents)
		})
	}
}
