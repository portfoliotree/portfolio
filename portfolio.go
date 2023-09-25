package portfolio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	"github.com/portfoliotree/portfolio/allocation"
	"github.com/portfoliotree/portfolio/backtest"
	"github.com/portfoliotree/portfolio/backtest/backtestconfig"
	"github.com/portfoliotree/portfolio/returns"
)

type Document struct {
	ID       primitive.ObjectID `json:"_id"      yaml:"_id"      bson:"_id"`
	Type     string             `json:"type"     yaml:"type"     bson:"type"`
	Metadata Metadata           `json:"metadata" yaml:"metadata" bson:"metadata"`
	Spec     Specification      `json:"spec"     yaml:"spec"     bson:"spec"`
}

type Metadata struct {
	Name        string    `json:"name,omitempty"        yaml:"name,omitempty"        bson:"name,omitempty"`
	Benchmark   Component `json:"benchmark,omitempty"   yaml:"benchmark,omitempty"   bson:"benchmark,omitempty"`
	Description string    `json:"description,omitempty" yaml:"description,omitempty" bson:"description,omitempty"`
	Privacy     string    `json:"privacy,omitempty"     yaml:"privacy,omitempty"     bson:"privacy,omitempty"`
}

// Specification models a portfolio.
type Specification struct {
	Assets []Component `json:"assets" yaml:"assets" bson:"assets"`
	Policy Policy      `json:"policy" yaml:"policy" bson:"policy"`
}

// ParseOneDocument decodes the contents of in to a Specification
// It supports a string containing YAML.
// The resulting Specification may have default values for unset fields.
func ParseOneDocument(in string) (Document, error) {
	result, err := ParseDocuments(strings.NewReader(in))
	if err != nil {
		return Document{}, err
	}
	if len(result) != 1 {
		return Document{}, fmt.Errorf("expected input to have exactly one portfolio especified")
	}
	return result[0], nil
}

const portfolioTypeName = "Portfolio"

// ParseDocuments decodes the contents of in to a list of Specifications
// The resulting Specification may have default values for unset fields.
func ParseDocuments(r io.Reader) ([]Document, error) {
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)
	var result []Document
	for index := 0; ; index++ {
		var document Document
		if err := dec.Decode(&document); err != nil {
			if err == io.EOF {
				return result, nil
			}
			return result, err
		}
		switch document.Type {
		case portfolioTypeName:
		default:
			return result, fmt.Errorf("incorrect specification type got %q but expected %q", document.Type, portfolioTypeName)
		}
		if len(document.Spec.Policy.Weights) > 0 && len(document.Spec.Policy.Weights) != len(document.Spec.Assets) {
			return result, errAssetAndWeightsLenMismatch(&document.Spec)
		}
		document.Spec.setDefaultPolicyWeightAlgorithm()
		result = append(result, document)
	}
}

func (pf *Specification) RemoveAsset(index int) error {
	if index < 0 || index >= len(pf.Assets) {
		return fmt.Errorf("asset index %d out of range the portfolio has %d asssets", index, len(pf.Assets))
	}
	pf.Assets = slices.Delete(pf.Assets, index, index+1)
	return nil
}

func (pf *Specification) Backtest(ctx context.Context, assets returns.Table, alg allocation.Algorithm) (backtest.Result, error) {
	return pf.BacktestWithStartAndEndTime(ctx, time.Time{}, time.Time{}, assets, alg)
}

func (pf *Specification) setDefaultPolicyWeightAlgorithm() {
	if len(pf.Policy.Weights) > 0 {
		pf.Policy.WeightsAlgorithm = (*allocation.ConstantWeights)(nil).Name()
	} else {
		pf.Policy.WeightsAlgorithm = (*allocation.EqualWeights)(nil).Name()
	}
}

func (pf *Specification) BacktestWithStartAndEndTime(ctx context.Context, start, end time.Time, assets returns.Table, alg allocation.Algorithm) (backtest.Result, error) {
	if alg == nil {
		var err error
		alg, err = pf.Algorithm(nil)
		if err != nil {
			return backtest.Result{}, err
		}
	}
	return backtest.Run(ctx, end, start, assets, alg,
		pf.Policy.WeightsAlgorithmLookBack.Function,
		pf.Policy.WeightsUpdatingInterval.CheckFunction(),
		pf.Policy.RebalancingInterval.CheckFunction(),
	)
}

type Policy struct {
	RebalancingInterval      backtestconfig.Interval `json:"rebalancing_interval,omitempty"               yaml:"rebalancing_interval,omitempty"               bson:"rebalancing_interval"`
	Weights                  []float64               `json:"weights,omitempty"                            yaml:"weights,omitempty"                            bson:"weights"`
	WeightsAlgorithm         string                  `json:"weights_algorithm,omitempty"                  yaml:"weights_algorithm,omitempty"                  bson:"weights_algorithm"`
	WeightsAlgorithmLookBack backtestconfig.Window   `json:"weights_algorithm_look_back_window,omitempty" yaml:"weights_algorithm_look_back_window,omitempty" bson:"weights_algorithm_look_back_window"`
	WeightsUpdatingInterval  backtestconfig.Interval `json:"weights_updating_interval,omitempty"          yaml:"weights_updating_interval,omitempty"          bson:"weights_updating_interval"`
}

func (pf *Specification) Algorithm(algorithmOptions []allocation.Algorithm) (allocation.Algorithm, error) {
	if len(algorithmOptions) == 0 {
		algorithmOptions = allocation.NewDefaultAlgorithmsList()
	}

	for _, alg := range algorithmOptions {
		if alg.Name() != pf.Policy.WeightsAlgorithm {
			continue
		}
		if se, ok := alg.(allocation.WeightSetter); ok {
			if len(pf.Policy.Weights) != len(pf.Assets) {
				return nil, errAssetAndWeightsLenMismatch(pf)
			}
			se.SetWeights(slices.Clone(pf.Policy.Weights))
		}
		return alg, nil // algorithm is known
	}

	return nil, errors.New("unknown algorithm")
}

func errAssetAndWeightsLenMismatch(spec *Specification) error {
	return fmt.Errorf("expected the number of policy weights to be the same as the number of assets got %d but expected %d", len(spec.Policy.Weights), len(spec.Assets))
}
