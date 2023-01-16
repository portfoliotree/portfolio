// Package backtest calculates portfolio returns and asset weights from historic asset returns.
// It re-balances asset weights and updates policies based on provided functions. See Run.
//
// DailyRebalancedWithStaticWeights is a simplified "back-tester" for calculating daily rebalanced returns of a portfolio
// given static policy asset weights.
//
//	Please remember, investing carries inherent risks including but not limited to the potential loss of principal. Past performance is no guarantee of future results. The data, equations, and calculations in these docs and code are for informational purposes only and should not be considered financial advice. It is important to carefully consider your own financial situation before making any investment decisions. You should seek the advice of a licensed financial professional before making any investment decisions. You should seek code review of an experienced software developer before consulting this library (or any library that imports it) to make investment decisions.
package backtest
