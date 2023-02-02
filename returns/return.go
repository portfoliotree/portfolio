package returns

import (
	"math"
	"time"
)

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
