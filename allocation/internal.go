package allocation

import (
	"errors"

	"github.com/portfoliotree/portfolio/returns"
)

func ensureEnoughReturns(assetReturns returns.Table) error {
	if assetReturns.NumberOfColumns() == 0 || assetReturns.NumberOfRows() < 2 {
		return errors.New("not enough data")
	}
	return nil
}

func isOnlyZeros(a []float64) bool {
	for _, v := range a {
		if v != 0 {
			return false
		}
	}
	return true
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
