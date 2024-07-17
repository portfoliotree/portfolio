package calculate

import (
	"context"
	"errors"

	"gonum.org/v1/gonum/optimize"
)

const (
	maxTries              = 50_000
	skipContextCheckCount = 500
	preCancelCheckTries   = 10_000
)

func checkTries(ctx context.Context, try int) error {
	switch {
	case try > preCancelCheckTries && try%skipContextCheckCount == 0:
		return ctx.Err()
	case try > maxTries:
		return errors.New("reached max tries to calculate policy")
	default:
		return nil
	}
}

func optWeights(ctx context.Context, weights []float64, fn func(ws []float64) float64) error {
	var (
		try = 0
		m   = &optimize.NelderMead{}
		s   = &optimize.Settings{
			Converger: &optimize.FunctionConverge{
				Absolute:   1e-10,
				Relative:   1,
				Iterations: 1000,
			},
		}
		ws = make([]float64, len(weights))
		p  = optimize.Problem{
			Func: func(x []float64) float64 {
				copy(ws, x)
				scaleToUnitRange(ws)
				return fn(ws)
			},
			Status: func() (optimize.Status, error) {
				err := checkTries(ctx, try)
				if err != nil {
					return optimize.RuntimeLimit, err
				}
				try++
				return optimize.NotTerminated, nil
			},
		}
	)
	optResult, err := optimize.Minimize(p, weights, s, m)
	if err != nil {
		return err
	}

	copy(weights, optResult.X)
	scaleToUnitRange(weights)

	return nil
}

func scaleToUnitRange(list []float64) {
	sum := 0.0
	for _, v := range list {
		sum += v
	}
	if sum == 0 {
		return
	}
	for i := range list {
		list[i] /= sum
	}
}
