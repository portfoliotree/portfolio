package allocation

import (
	"slices"

	"github.com/portfoliotree/portfolio/backtest"
)

type Algorithm interface {
	backtest.PolicyWeightCalculator
	Name() string
}

func NewDefaultAlgorithmsList() []Algorithm {
	return []Algorithm{
		new(ConstantWeights),
		new(EqualWeights),
		new(EqualInverseVariance),
		new(EqualRiskContribution),
		new(EqualVolatility),
		new(EqualInverseVolatility),
	}
}

func AlgorithmNames(algorithmOptions []Algorithm) []string {
	names := make([]string, 0, len(algorithmOptions))
	for _, alg := range algorithmOptions {
		names = append(names, alg.Name())
	}
	slices.Sort(names)
	names = slices.Compact(names)
	return names
}

type WeightSetter interface {
	SetWeights([]float64)
}

func AlgorithmRequiresWeights(alg Algorithm) bool {
	_, ok := alg.(WeightSetter)
	return ok
}
