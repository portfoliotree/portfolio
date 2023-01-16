package backtestconfig_test

import (
	"github.com/portfoliotree/portfolio/backtest"
	"github.com/portfoliotree/portfolio/backtest/backtestconfig"
)

var (
	_ backtest.PolicyWeightCalculator = backtestconfig.ConstantWeights{}
	_ backtest.PolicyWeightCalculator = backtestconfig.EqualWeights{}
	_ backtest.PolicyWeightCalculator = backtestconfig.PolicyWeightCalculatorFunc(nil)
)
