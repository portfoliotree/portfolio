package calculate

import (
	"gonum.org/v1/gonum/stat"
)

func CorrelationMatrix(values [][]float64) [][]float64 {
	n := len(values)
	if len(values) == 0 {
		return nil
	}
	m := make([][]float64, n)
	for i := range m {
		m[i] = make([]float64, n)
	}
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			if i == j {
				m[i][j] = 1.0
				continue
			}
			vi, vj := values[i], values[j]
			c := stat.Correlation(vi, vj, nil)
			m[i][j] = c
			m[j][i] = c
		}
	}
	return m
}
