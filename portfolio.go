package portfolio

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

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
		document.Spec.setDefaultPolicyWeightAlgorithm()
		if err := document.Spec.ensureEqualNumberOfWeightsAndAssets(); err != nil {
			return result, err
		}
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

func (pf *Specification) Backtest(ctx context.Context, assets returns.Table, weightsAlgorithm backtestconfig.PolicyWeightCalculatorFunc) (backtest.Result, error) {
	return pf.BacktestWithStartAndEndTime(ctx, time.Time{}, time.Time{}, assets, weightsAlgorithm)
}

const (
	PolicyAlgorithmEqualWeights    = "EqualWeights"
	PolicyAlgorithmConstantWeights = "ConstantWeights"
)

func (pf *Specification) setDefaultPolicyWeightAlgorithm() {
	if pf.Policy.WeightsAlgorithm != "" {
		return
	}
	if len(pf.Policy.Weights) > 0 {
		pf.Policy.WeightsAlgorithm = PolicyAlgorithmConstantWeights
	} else {
		pf.Policy.WeightsAlgorithm = PolicyAlgorithmEqualWeights
	}
}

func (pf *Specification) ensureEqualNumberOfWeightsAndAssets() error {
	switch pf.Policy.WeightsAlgorithm {
	case PolicyAlgorithmConstantWeights:
		if len(pf.Policy.Weights) != len(pf.Assets) {
			return fmt.Errorf("the number of assets and number of weights must be equal: len(assets) is %d and len(weights) is %d", len(pf.Assets), len(pf.Policy.Weights))
		}
	}
	return nil
}

func (pf *Specification) policyWeightFunction(weights backtestconfig.PolicyWeightCalculatorFunc) (backtestconfig.PolicyWeightCalculatorFunc, error) {
	switch pf.Policy.WeightsAlgorithm {
	case PolicyAlgorithmEqualWeights:
		return backtestconfig.EqualWeights{}.PolicyWeights, nil
	case PolicyAlgorithmConstantWeights:
		return backtestconfig.ConstantWeights(pf.Policy.Weights).PolicyWeights, nil
	default:
		if weights == nil {
			return nil, fmt.Errorf("policy %q not supported by the backtest runner", pf.Policy.WeightsAlgorithm)
		}
		return weights, nil
	}
}

func (pf *Specification) BacktestWithStartAndEndTime(ctx context.Context, start, end time.Time, assets returns.Table, weightsFn backtestconfig.PolicyWeightCalculatorFunc) (backtest.Result, error) {
	if err := pf.ensureEqualNumberOfWeightsAndAssets(); err != nil {
		return backtest.Result{}, err
	}
	var err error
	weightsFn, err = pf.policyWeightFunction(weightsFn)
	if err != nil {
		return backtest.Result{}, err
	}

	return backtest.Run(ctx, end, start, assets, weightsFn,
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
