package allocation

import (
	"context"
	"time"

	"github.com/portfoliotree/timetable"

	"github.com/portfoliotree/portfolio/calculate"
)

type EqualInverseVariance struct{}

func (cw *EqualInverseVariance) Name() string { return "Equal Inverse Variance" }

func (*EqualInverseVariance) PolicyWeights(_ context.Context, _ time.Time, assetReturns timetable.Compact[float64], ws []float64) ([]float64, error) {
	if isOnlyZeros(ws) {
		for i := range ws {
			ws[i] = 1.0
		}
		scaleToUnitRange(ws)
	}

	err := ensureEnoughReturns(assetReturns)
	if err != nil {
		return ws, err
	}

	vols := make([]float64, 0, assetReturns.NumberOfColumns())
	for _, vs := range assetReturns.UnderlyingValues() {
		vols = append(vols, calculate.RiskFromStdDev(vs))
	}
	calculate.InverseVarianceWeights(ws, vols)
	return ws, nil
}
