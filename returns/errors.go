package returns

type ErrorNoReturns struct{}

func (err ErrorNoReturns) Error() string {
	return "no returns"
}
