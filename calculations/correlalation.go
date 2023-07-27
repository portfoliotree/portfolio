package calculations

import (
	"math"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

func CorrelationMatrix(values [][]float64) *mat.Dense {
	if len(values) == 0 {
		return nil
	}
	m := mat.NewDense(len(values), len(values), nil)
	mp := make(map[int][]float64, len(values))
	min := math.MaxInt
	for i := range values {
		mp[i] = values[i]
		if len(mp[i]) < min {
			min = len(mp[i])
		}
	}
	for i := 0; i < len(values); i++ {
		for j := i; j < len(values); j++ {
			if i == j {
				m.Set(i, j, 1)
				continue
			}
			vi, vj := mp[i], mp[j]
			c := stat.Correlation(vi[:min], vj[:min], nil)
			m.Set(i, j, c)
			m.Set(j, i, c)
		}
	}
	return m
}

func DenseToSlices(dense *mat.Dense) [][]float64 {
	iL, jL := dense.Dims()
	result := make([][]float64, iL)
	d := dense.RawMatrix().Data
	for i := range result {
		result[i] = d[i*jL : (i+1)*jL]
	}
	return result
}
