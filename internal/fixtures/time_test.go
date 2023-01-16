package fixtures_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/portfolio/internal/fixtures"
)

func Test_days(t *testing.T) {
	assert.Equal(t, fixtures.T(t, fixtures.Day0).Weekday(), time.Thursday)
	assert.Equal(t, fixtures.T(t, fixtures.Day1).Weekday(), time.Friday)
	assert.Equal(t, fixtures.T(t, fixtures.Day2).Weekday(), time.Monday)
	assert.Equal(t, fixtures.T(t, fixtures.Day3).Weekday(), time.Tuesday)
}
