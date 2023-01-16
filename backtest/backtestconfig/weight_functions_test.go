package backtestconfig_test

import (
	"github.com/portfoliotree/portfolio/backtest"
	"github.com/portfoliotree/portfolio/backtest/backtestconfig"
)

var (
	_ backtest.PolicyWeightsFunc = backtestconfig.ConstantWeights{}.PolicyWeights
	_ backtest.PolicyWeightsFunc = backtestconfig.EqualWeights{}.PolicyWeights
)
