package backtest_test

import (
	"context"
	"testing"

	"github.com/portfoliotree/portfolio"
	"github.com/portfoliotree/portfolio/backtest"
	"github.com/portfoliotree/portfolio/backtest/backtestconfig"
	"github.com/portfoliotree/portfolio/portfoliotest"
	"github.com/portfoliotree/portfolio/returns"
)

func BenchmarkRun1Q(b *testing.B) {
	rs := benchmarkRunReturns(b)
	rs = rs.Between(rs.LastTime(), rs.LastTime().AddDate(0, -3, 0))
	benchmarkRun(b, rs)
}

func BenchmarkRun1Y(b *testing.B) {
	rs := benchmarkRunReturns(b)
	rs = rs.Between(rs.LastTime(), rs.LastTime().AddDate(-1, 0, 0))
	benchmarkRun(b, rs)
}

func BenchmarkRun3Y(b *testing.B) {
	rs := benchmarkRunReturns(b)
	rs = rs.Between(rs.LastTime(), rs.LastTime().AddDate(-3, 0, 0))
	benchmarkRun(b, rs)
}

func BenchmarkRunMax(b *testing.B) {
	rs := benchmarkRunReturns(b)
	rs = rs.Between(rs.LastTime(), rs.FirstTime())
	benchmarkRun(b, rs)
}

func benchmarkRun(b *testing.B, table returns.Table) {
	b.Helper()
	end := table.LastTime()
	start := table.FirstTime()
	fn := backtestconfig.EqualWeights{}
	lookback := backtestconfig.OneQuarterWindow.Function
	rebalance := backtestconfig.Daily()
	updatePolicyWeights := backtestconfig.Monthly()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := backtest.Run(ctx, end, start, table, fn, lookback, rebalance, updatePolicyWeights)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkRunReturns(b *testing.B) returns.Table {
	p := portfoliotest.ComponentReturnsProvider()
	assets := []portfolio.Component{{ID: "ACWI"}, {ID: "AGG"}}
	table, err := p.ComponentReturnsTable(context.Background(), assets...)
	if err != nil {
		b.Fatal(err)
	}
	return table
}
