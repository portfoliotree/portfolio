package portfolio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	"github.com/portfoliotree/portfolio/backtest"
	"github.com/portfoliotree/portfolio/backtest/backtestconfig"
)

type Specification struct {
	Name      string      `yaml:"name"`
	Benchmark Component   `yaml:"benchmark"`
	Assets    []Component `yaml:"assets"`
	Policy    Policy      `yaml:"policy"`

	Filepath  string `yaml:"-"`
	FileIndex int    `yaml:"-"`
}

type TypedSpecificationFile[S interface {
	Specification
}] struct {
	ID   string `yaml:"id"`
	Type string `yaml:"type"`
	Spec S      `yaml:"spec"`
}

// ParseOneSpecification decodes the contents of in to a Specification
// It supports a string containing YAML.
// The resulting Specification may have default values for unset fields.
func ParseOneSpecification(in string) (Specification, error) {
	result, err := ParseSpecifications(strings.NewReader(in))
	if err != nil {
		return Specification{}, err
	}
	if len(result) != 1 {
		return Specification{}, fmt.Errorf("expected input to have exactly one portfolio especified")
	}
	return result[0], nil
}

const portfolioTypeName = "Portfolio"

// ParseSpecifications decodes the contents of in to a list of Specifications
// The resulting Specification may have default values for unset fields.
func ParseSpecifications(r io.Reader) ([]Specification, error) {
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)
	var result []Specification
	for {
		var spec TypedSpecificationFile[Specification]
		if err := dec.Decode(&spec); err != nil {
			if err == io.EOF {
				return result, nil
			}
			return result, err
		}
		switch spec.Type {
		case portfolioTypeName:
		default:
			return result, fmt.Errorf("incorrect specification type got %q but expected %q", spec.Type, portfolioTypeName)
		}
		pf := spec.Spec
		pf.setDefaultPolicyWeightAlgorithm()
		if err := pf.ensureEqualNumberOfWeightsAndAssets(); err != nil {
			return result, err
		}
		result = append(result, pf)
	}
}

func (pf *Specification) RemoveAsset(index int) error {
	if index < 0 || index >= len(pf.Assets) {
		return fmt.Errorf("asset index %d out of range the portfolio has %d asssets", index, len(pf.Assets))
	}
	pf.Assets = slices.Delete(pf.Assets, index, index+1)
	return nil
}

func (pf *Specification) Backtest(ctx context.Context, weightsAlgorithm backtestconfig.PolicyWeightCalculatorFunc) (backtest.Result, error) {
	return pf.BacktestWithStartAndEndTime(ctx, time.Time{}, time.Time{}, weightsAlgorithm)
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

func (pf *Specification) BacktestWithStartAndEndTime(ctx context.Context, start, end time.Time, weightsFn backtestconfig.PolicyWeightCalculatorFunc) (backtest.Result, error) {
	if err := pf.ensureEqualNumberOfWeightsAndAssets(); err != nil {
		return backtest.Result{}, err
	}
	var err error
	weightsFn, err = pf.policyWeightFunction(weightsFn)
	if err != nil {
		return backtest.Result{}, err
	}

	assets, err := pf.AssetReturns(ctx)
	if err != nil {
		return backtest.Result{}, err
	}

	if start.IsZero() {
		start = assets.FirstTime()
	}
	if end.IsZero() {
		end = assets.LastTime()
	}

	return backtest.Run(ctx, end, start, assets, weightsFn,
		pf.Policy.WeightsAlgorithmLookBack.Function,
		pf.Policy.WeightsUpdatingInterval.CheckFunction(),
		pf.Policy.RebalancingInterval.CheckFunction(),
	)
}

type Policy struct {
	RebalancingInterval backtestconfig.Interval `yaml:"rebalancing_interval,omitempty"`

	Weights                  []float64               `yaml:"weights,omitempty"`
	WeightsAlgorithm         string                  `yaml:"weights_algorithm,omitempty"`
	WeightsAlgorithmLookBack backtestconfig.Window   `yaml:"weights_algorithm_look_back_window,omitempty"`
	WeightsUpdatingInterval  backtestconfig.Interval `yaml:"weights_updating_interval,omitempty"`
}

func (pf *Specification) ParseValues(q url.Values) error {
	if q.Has("asset-id") {
		pf.Assets = pf.Assets[:0]
		for _, assetID := range q["asset-id"] {
			pf.Assets = append(pf.Assets, Component{ID: assetID})
		}
	}
	if q.Has("benchmark-id") {
		pf.Benchmark.ID = q.Get("benchmark-id")
	}
	if q.Has("name") {
		pf.Name = q.Get("name")
	}
	if q.Has("filepath") {
		pf.Filepath = q.Get("filepath")
	}
	if q.Has("policy-rebalance") {
		pf.Policy.RebalancingInterval = backtestconfig.Interval(q.Get("policy-rebalance"))
	}
	if q.Has("policy-weights-algorithm") {
		pf.Policy.WeightsAlgorithm = q.Get("policy-weights-algorithm")
	}
	if q.Has("policy-weight") {
		pf.Policy.Weights = pf.Policy.Weights[:0]
		for i, weight := range q["policy-weight"] {
			f, err := strconv.ParseFloat(weight, 64)
			if err != nil {
				return fmt.Errorf("failed to parse policy weight at indx %d: %w", i, err)
			}
			pf.Policy.Weights = append(pf.Policy.Weights, f)
		}
	}
	if q.Has("policy-update-weights") {
		pf.Policy.WeightsUpdatingInterval = backtestconfig.Interval(q.Get("policy-update-weights"))
	}
	if q.Has("policy-weight-algorithm-look-back") {
		pf.Policy.WeightsAlgorithmLookBack = backtestconfig.Window(q.Get("policy-weight-algorithm-look-back"))
	}
	pf.filterEmptyAssetIDs()
	return pf.Validate()
}

func (pf *Specification) Values() url.Values {
	q := make(url.Values)
	if pf.Name != "" {
		q.Set("name", pf.Name)
	}
	if pf.Benchmark.ID != "" {
		q.Set("benchmark-id", pf.Benchmark.ID)
	}
	if pf.Filepath != "" {
		q.Set("filepath", pf.Filepath)
	}
	if pf.Assets != nil {
		for _, asset := range pf.Assets {
			q.Add("asset-id", asset.ID)
		}
	}
	if pf.Policy.RebalancingInterval != "" {
		q.Set("policy-rebalance", pf.Policy.RebalancingInterval.String())
	}
	if pf.Policy.WeightsAlgorithm != "" {
		q.Set("policy-weights-algorithm", pf.Policy.WeightsAlgorithm)
	}
	if pf.Policy.Weights != nil {
		for _, w := range pf.Policy.Weights {
			q.Add("policy-weight", strconv.FormatFloat(w, 'f', 4, 64))
		}
	}
	if pf.Policy.WeightsUpdatingInterval != "" {
		q.Set("policy-update-weights", string(pf.Policy.WeightsUpdatingInterval))
	}
	if pf.Policy.WeightsAlgorithmLookBack != "" {
		q.Set("policy-weight-algorithm-look-back", pf.Policy.WeightsAlgorithmLookBack.String())
	}
	return q
}

func (pf *Specification) Validate() error {
	var list []error
	for _, asset := range pf.Assets {
		list = append(list, asset.Validate())
	}
	if pf.Benchmark.ID != "" {
		if err := pf.Benchmark.Validate(); err != nil {
			list = append(list, err)
		}
	}
	return errors.Join(list...)
}

func (pf *Specification) filterEmptyAssetIDs() {
	filtered := pf.Assets[:0]
	for _, asset := range pf.Assets {
		if asset.ID != "" {
			filtered = append(filtered, asset)
		}
	}
	pf.Assets = filtered
}
