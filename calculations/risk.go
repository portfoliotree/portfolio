package calculations

import (
	"errors"
	"fmt"
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

func DenseMatrixToFloatSlices(dense *mat.Dense) [][]float64 {
	if dense == nil {
		return make([][]float64, 0)
	}
	iLength, jLength := dense.Dims()
	result := make([][]float64, iLength)
	d := dense.RawMatrix().Data
	for i := range result {
		result[i] = d[i*jLength : (i+1)*jLength]
	}
	return result
}

func Risk(risks, weights []float64, correlations [][]float64) (float64, error) {
	numberOfAssets := len(risks)

	for i, r := range risks {
		risks[i] = r / 100.0
	}
	if sum := floats.Sum(weights); math.Abs(100.0-sum) > 0.0001 {
		return 0, fmt.Errorf(`weights must add up to 100 but got %.2f`, sum)
	}
	for i, w := range weights {
		weights[i] = w / 100.0
	}

	if len(weights) != numberOfAssets {
		return 0, errors.New("the number of risks must be the same as the number of weights")
	}

	if len(correlations) != numberOfAssets {
		return 0, fmt.Errorf("correlations must be n by n where n is the length of risks; the number of rows is should be %d", numberOfAssets)
	}

	cor := mat.NewDense(len(risks), len(risks), nil)
	for i, row := range correlations {
		if len(row) != numberOfAssets {
			return 0, fmt.Errorf("correlations must be n by n where n is the length of risks; row %d has %d it should have %d", i, len(row), numberOfAssets)
		}
		for j, v := range correlations[i] {
			if v < -1 || v > 1 {
				return 0, fmt.Errorf("correlation values must be within [-1, 1] the value %.2f at row %d column %d is not in this range", v, i+1, j+1)
			}
			cor.Set(i, j, v)
		}
	}

	r, _, _ := RiskFromRiskContribution(risks, weights, cor)

	return r, nil
}

func RiskFromRiskContribution(risks, weights []float64, correlations *mat.Dense) (float64, []float64, []float64) {
	n := len(risks)
	nSquared := n * n

	// START memory allocation stuff
	bufferSize := 3*nSquared + 3*n
	b := make([]float64, bufferSize)
	b1 := b[:nSquared:nSquared]
	b = b[len(b1):]
	b2 := b[:nSquared:nSquared]
	b = b[len(b2):]
	b3 := b[:n:n]
	b = b[len(b3):]
	b4 := b[:nSquared:nSquared]
	b = b[len(b4):]
	v := b[:n:n]
	b = b[len(v):]
	rw := b[:n:n]
	// END memory allocation stuff

	d := mat.NewDense(n, n, b1)

	for i := 0; i < len(risks); i++ {
		d.Set(i, i, risks[i])
	}

	V := mat.NewDense(n, n, b2)
	V.Product(d, correlations, d)

	mx1 := mat.NewDense(n, 1, b3)
	mx1.Product(V, mat.NewDense(n, 1, weights))
	mx2 := mat.NewDense(n, n, b4)
	mx2.Product(mx1, mat.NewDense(1, n, weights))

	for i := 0; i < n; i++ {
		v[i] = mx2.At(i, i)
	}
	pv := floats.Sum(v)
	copy(rw, v)
	for i := range rw {
		rw[i] /= pv
	}

	return math.Sqrt(pv), v, rw
}

func ExpectedRisk(risks, weights []float64, correlations *mat.Dense) float64 {
	r, _, _ := RiskFromRiskContribution(risks, weights, correlations)
	return r
}

func NumberOfBets(weightedAverageRisk float64, portfolioRisk float64) (float64, error) {
	if portfolioRisk == 0 {
		return 0, errors.New("can't divide by 0")
	}

	abRatio := weightedAverageRisk / portfolioRisk
	bets := abRatio * abRatio

	return bets, nil
}

func RiskFromStdDev(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	// PopStdDev may be the correct method to call here for backwards looking stuff.
	// For now, we use sample standard deviation.
	return stat.StdDev(values, nil)
}

func WeightedAverageRisk(weights, risks []float64) float64 {
	return stat.Mean(risks, weights)
}
