package backtest

import (
	"errors"
	"time"
)

type WeightSnapshot struct {
	Time              time.Time `json:"time"`
	Weights           []float64 `json:"weights"`
	IsRebalanceDay    bool      `json:"r" bson:"r,omitempty"`
	IsPolicyUpdateDay bool      `json:"p" bson:"p,omitempty"`
}

// WeightSnapshotList should be ordered where the WeightSnapshot at index
// 0 is the most recent and at len(list)-1 least recent.
type WeightSnapshotList []WeightSnapshot

func (list WeightSnapshotList) reverse() {
	for i := 0; i < len(list)/2; i++ {
		t := list[i]
		list[i] = list[len(list)-1-i]
		list[len(list)-1-i] = t
	}
}

func (list WeightSnapshotList) Times() []time.Time {
	var result []time.Time
	for _, t := range list {
		result = append(result, t.Time)
	}
	return result
}

func (list WeightSnapshotList) AverageWeightForIndex(index int) (float64, error) {
	sum := 0.0

	for _, snapshot := range list {
		if index >= len(snapshot.Weights) || index < 0 {
			return 0.0, errors.New("index out of range")
		}
		sum += snapshot.Weights[index]
	}

	return sum / float64(len(list)), nil
}
