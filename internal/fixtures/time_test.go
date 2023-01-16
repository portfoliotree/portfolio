package fixtures_test

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/portfoliotree/portfolio/internal/fixtures"
)

func Test_days(t *testing.T) {
	o := NewWithT(t)
	o.Expect(fixtures.T(t, fixtures.Day0).Weekday()).To(Equal(time.Thursday))
	o.Expect(fixtures.T(t, fixtures.Day1).Weekday()).To(Equal(time.Friday))
	o.Expect(fixtures.T(t, fixtures.Day2).Weekday()).To(Equal(time.Monday))
	o.Expect(fixtures.T(t, fixtures.Day3).Weekday()).To(Equal(time.Tuesday))
}
