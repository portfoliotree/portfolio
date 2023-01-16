package backtest

type ErrorNotEnoughData struct{}

func (err ErrorNotEnoughData) Error() string { return "not enough data" }
