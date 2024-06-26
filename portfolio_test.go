package portfolio_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/portfoliotree/portfolio"
	"github.com/portfoliotree/portfolio/returns"
)

func ExampleParse() {
	// language=yaml
	specYAML := `
---
type: Portfolio
metadata:
  name: 60/40
  benchmark: BIGPX
spec:
  assets: [ACWI, AGG]
  policy:
    weights: [60, 40]
    weights_algorithm: Constant Weights
    rebalancing_interval: Quarterly
`

	pf, err := portfolio.ParseOneDocument(specYAML)
	if err != nil {
		panic(err)
	}
	fmt.Println("Name:", pf.Metadata.Name)
	fmt.Println("Alg:", pf.Spec.Policy.WeightsAlgorithm)

	// Output:
	// Name: 60/40
	// Alg: Constant Weights
}

func ExampleOpen() {
	portfolios, err := portfolio.ParseSpecificationFile(filepath.Join("examples", "60-40_portfolio.yml"))
	if err != nil {
		panic(err)
	}
	pf := portfolios[0]
	fmt.Println("Name:", pf.Metadata.Name)
	fmt.Println("Alg:", pf.Spec.Policy.WeightsAlgorithm)

	// Output:
	// Name: 60/40
	// Alg: Constant Weights
}

func TestParse(t *testing.T) {
	for _, tt := range []struct {
		Name                string
		SpecYAML            string
		ErrorStringContains string
		Portfolio           portfolio.Specification
	}{
		{
			Name:                "invalid yaml",
			SpecYAML:            "---}",
			ErrorStringContains: "yaml: unmarshal errors",
		},
		{
			Name: "wrong type",
			// language=yaml
			SpecYAML:            `type: Banana`,
			ErrorStringContains: "incorrect specification type",
		},
		{
			Name: "the number of assets and policy weights do not match",
			// language=yaml
			SpecYAML:            `{type: Portfolio, spec: {assets: ["a"], policy: {weights: [1, 2]}}}`,
			ErrorStringContains: "expected the number of policy weights to be the same as the number of assets",
		},
		{
			Name: "component field is invalid",
			// language=yaml
			SpecYAML:            `{type: Portfolio, spec: {benchmark: {id: []}}}`,
			ErrorStringContains: "yaml: unmarshal errors:",
		},
		{
			Name: "empty input",
			// language=yaml
			SpecYAML:            ``,
			ErrorStringContains: "exactly one portfolio",
		},
		{
			Name: "empty input",
			// language=yaml
			SpecYAML: `
{type: Portfolio}
---
{type: Portfolio}`,
			ErrorStringContains: "exactly one portfolio",
		},
		{
			Name: "component kind is not correct",
			// language=yaml
			SpecYAML:            `{type: Portfolio, metadata: {benchmark: []}}`,
			ErrorStringContains: "wrong YAML type:",
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			p, err := portfolio.ParseOneDocument(tt.SpecYAML)
			if tt.ErrorStringContains == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.Portfolio, p)
			} else {
				assert.ErrorContains(t, err, tt.ErrorStringContains)
			}
		})
	}
}

func ExampleSpecification_Backtest() {
	// language=yaml
	portfolioSpecYAML := `---
type: Portfolio
metadata:
  name: 60/40
  benchmark: BIGPX
spec:
  assets: [ACWI, AGG]
  policy:
    weights: [60, 40]
    weights_algorithm: ConstantWeights
    rebalancing_interval: Quarterly
`

	pf, err := portfolio.ParseOneDocument(portfolioSpecYAML)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	assets, err := pf.Spec.AssetReturns(ctx)
	if err != nil {
		panic(err)
	}
	result, err := pf.Spec.Backtest(ctx, assets, nil)
	if err != nil {
		panic(err)
	}

	portfolioReturns := result.Returns()
	fmt.Printf("Annualized Risk: %.2f%%\n", portfolioReturns.AnnualizedRisk()*100)
	fmt.Printf("Annualized Return: %.2f%%\n", portfolioReturns.AnnualizedArithmeticReturn()*100)
	fmt.Printf("Backtest start date: %s\n", result.ReturnsTable.FirstTime().Format(time.DateOnly))
	fmt.Printf("Backtest end date: %s\n", result.ReturnsTable.LastTime().Format(time.DateOnly))

	// Output:
	// Annualized Risk: 11.46%
	// Annualized Return: 5.10%
	// Backtest start date: 2008-03-31
	// Backtest end date: 2023-06-14
}

func TestPortfolio_Backtest(t *testing.T) {
	for _, tt := range []struct {
		Name                  string
		PortfolioSpecFilePath string
		Portfolio             portfolio.Specification
		ctx                   context.Context

		ErrorSubstring string
	}{
		{
			Name: "wrong number of weights",
			Portfolio: portfolio.Specification{
				Assets: []portfolio.Component{{ID: "AAPL"}},
				Policy: portfolio.Policy{
					Weights:          []float64{50, 50},
					WeightsAlgorithm: "Constant Weights",
				},
			},
			ctx:            context.Background(),
			ErrorSubstring: "expected the number of policy weights to be the same as the number of assets",
		},
		{
			Name: "unknown policy algorithm",
			Portfolio: portfolio.Specification{
				Assets: []portfolio.Component{{ID: "AAPL"}},
				Policy: portfolio.Policy{
					Weights:          []float64{50, 50},
					WeightsAlgorithm: "unknown",
				},
			},
			ctx:            context.Background(),
			ErrorSubstring: `unknown algorithm`,
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			pf := tt.Portfolio
			_, err := pf.Backtest(tt.ctx, returns.NewTable([]returns.List{{}}), nil)
			if tt.ErrorSubstring == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.ErrorSubstring)
			}
		})
	}
}

func TestPortfolio_Backtest_custom_function(t *testing.T) {
	_, err := (&portfolio.Specification{
		Assets: []portfolio.Component{
			{ID: "AAPL"},
			{ID: "GOOG"},
		},
	}).Backtest(context.Background(), returns.NewTable([]returns.List{{}}), ErrorAlg{})
	assert.EqualError(t, err, "lemon")
}

func Test_Portfolio_Validate(t *testing.T) {
	for _, tt := range []struct {
		Name      string
		Portfolio portfolio.Document
		ExpectErr bool
	}{
		{
			Name: "okay", Portfolio: portfolio.Document{
				Type: "Portfolio",
			}, ExpectErr: false,
		},
		{
			Name: "bad asset",
			Portfolio: portfolio.Document{
				Spec: portfolio.Specification{
					Assets: []portfolio.Component{{ID: "_"}},
				},
			},
			ExpectErr: true,
		},
		{
			Name: "benchmark",
			Portfolio: portfolio.Document{
				Metadata: portfolio.Metadata{
					Benchmark: portfolio.Component{ID: "()"},
				},
			},
			ExpectErr: true,
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.Portfolio.Validate()
			if tt.ExpectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPortfolio_RemoveAsset(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var zero portfolio.Specification
		require.Error(t, zero.RemoveAsset(0))
	})

	t.Run("empty", func(t *testing.T) {
		pf := portfolio.Specification{
			Assets: []portfolio.Component{},
		}
		require.Error(t, pf.RemoveAsset(0))
	})

	t.Run("remove one", func(t *testing.T) {
		pf := portfolio.Specification{
			Assets: []portfolio.Component{
				{ID: "banana"},
			},
		}
		require.NoError(t, pf.RemoveAsset(0))
		require.Len(t, pf.Assets, 0)
	})

	t.Run("remove one keep first", func(t *testing.T) {
		pf := portfolio.Specification{
			Assets: []portfolio.Component{
				{ID: "orange"},
				{ID: "banana"},
			},
		}
		require.NoError(t, pf.RemoveAsset(1))
		require.Equal(t, []portfolio.Component{{ID: "orange"}}, pf.Assets)
	})

	t.Run("remove one keep last", func(t *testing.T) {
		pf := portfolio.Specification{
			Assets: []portfolio.Component{
				{ID: "banana"},
				{ID: "orange"},
			},
		}
		require.NoError(t, pf.RemoveAsset(0))
		require.Equal(t, []portfolio.Component{{ID: "orange"}}, pf.Assets)
	})

	t.Run("out of bounds", func(t *testing.T) {
		pf := portfolio.Specification{
			Assets: []portfolio.Component{
				{ID: "banana"},
				{ID: "orange"},
			},
		}
		require.Error(t, pf.RemoveAsset(3))
		require.Equal(t, []portfolio.Component{{ID: "banana"}, {ID: "orange"}}, pf.Assets)
	})

	t.Run("negative index of bounds", func(t *testing.T) {
		pf := portfolio.Specification{
			Assets: []portfolio.Component{
				{ID: "banana"},
				{ID: "orange"},
			},
		}
		require.Error(t, pf.RemoveAsset(-1))
	})
}

type ErrorAlg struct{}

func (ErrorAlg) Name() string { return "" }

func (ErrorAlg) PolicyWeights(ctx context.Context, today time.Time, assets returns.Table, currentWeights []float64) ([]float64, error) {
	return nil, fmt.Errorf("lemon")
}
