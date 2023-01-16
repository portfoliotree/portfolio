package returns

import (
	"math"
	"time"
)

const DateLayout = "2006-01-02"

type Return struct {
	Value float64   `json:"value" bson:"value"`
	Time  time.Time `json:"time" bson:"time"`
}

func New(t time.Time, v float64) Return {
	if math.IsNaN(v) {
		panic("return value must be a number")
	}
	return Return{Time: t, Value: v}
}
