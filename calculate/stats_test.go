package calculate

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/floats"
)

func TestMetrics(t *testing.T) {
	data, err := loadTestdataReturns(filepath.FromSlash("testdata/metrics.tsv"), 3)
	require.NoError(t, err)
	cash, portfolio, benchmark := data[0], data[1], data[2]
	assert.Len(t, cash, 10_000)
	assert.Len(t, portfolio, 10_000)
	assert.Len(t, benchmark, 10_000)
	for i := range data {
		slices.Reverse(data[i])
	}

	// Downside Volatility	13.03%	=SQRT(AVERAGE(K2:K10001))*SQRT($N$1)
	t.Run("Downside Volatility", func(t *testing.T) {
		result := DownsideVolatility(portfolio, PeriodsPerYear)
		assert.InDelta(t, 0.1303, result, 0.0001)
	})

	// Sortino Ratio	 0.17 	=(N6-N5)/N12
	t.Run("Sortino Ratio", func(t *testing.T) {
		downsideVol := DownsideVolatility(portfolio, PeriodsPerYear)
		result := SortinoRatio(portfolio, cash, downsideVol, PeriodsPerYear)
		assert.InDelta(t, 0.17, result, 0.01)
	})

	// Calmar Ratio	 0.06 	=N7/N17
	t.Run("Calmar Ratio", func(t *testing.T) {
		maxDrawdown, _ := MaxDrawdown(portfolio)
		result := CalmarRatio(portfolio, cash, maxDrawdown, PeriodsPerYear)
		assert.InDelta(t, 0.059, result, 0.001)
	})

	// Ulcer Index	13.82	=SQRT(AVERAGE(I2:I10001))*SQRT(N1)
	t.Run("Ulcer Index", func(t *testing.T) {
		result := UlcerIndex(portfolio, PeriodsPerYear)
		assert.InDelta(t, 13.82, result, 0.01)
	})

	// Max Drawdown	37.38%	=1-MIN(H2:H10001)
	t.Run("Max Drawdown", func(t *testing.T) {
		result, _ := MaxDrawdown(portfolio)
		assert.InDelta(t, 0.3738, result, 0.0001)
	})

	// Tracking Error	6.27%	=STDEV.P(F2:F10001)*SQRT(N1)
	t.Run("Tracking Error", func(t *testing.T) {
		excess := make([]float64, len(portfolio))
		floats.SubTo(excess, portfolio, benchmark)
		result := TrackingError(excess, PeriodsPerYear)
		assert.InDelta(t, 0.0627, result, 0.0001)
	})

	// Information Ratio	 0.06 	=N9/N10
	t.Run("Information Ratio", func(t *testing.T) {
		result := InformationRatio(portfolio, benchmark, PeriodsPerYear)
		assert.InDelta(t, 0.06, result, 0.01)
	})

	// Beta to Benchmark	0.75	=SLOPE(D2:D10001,E2:E10001)
	t.Run("Beta to Benchmark", func(t *testing.T) {
		result := BetaToBenchmark(portfolio, benchmark)
		assert.InDelta(t, 0.748, result, 0.001)
	})

	// VaR (5% Confidence)	 (2,148.94)	=N18*N11*N2
	t.Run("Value at Risk aka VaR", func(t *testing.T) {
		const (
			portfolioValue  = 10_000.00
			confidenceLevel = 0.95
		)
		result := ValueAtRisk(portfolio, portfolioValue, confidenceLevel, PeriodsPerYear)
		assert.InDelta(t, 2148.94, result, 0.01)
	})
}

func loadTestdataReturns(fileName string, columns int) ([][]float64, error) {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(bytes.NewBuffer(buf))
	r.FieldsPerRecord = columns
	r.Comma = '\t'
	r.ReuseRecord = true
	r.TrimLeadingSpace = true
	data := make([][]float64, columns)
	for i := 0; ; i++ {
		line, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		if i < 1 || len(line) != 3 {
			continue
		}
		for j := range line {
			f, err := strconv.ParseFloat(line[j], 64)
			if err != nil {
				return nil, err
			}
			data[j] = append(data[j], f)
		}
	}
	return data, nil
}
