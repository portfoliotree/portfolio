package portfolio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponent_Validate(t *testing.T) {
	for _, tt := range []struct {
		Name           string
		Component      Component
		ErrorSubstring string
	}{
		{
			Name: "with undefined ID",
			Component: Component{
				ID: "undefined",
			},
			ErrorSubstring: "undefined",
		},
		{
			Name: "with ID not set",
			Component: Component{
				ID: "",
			},
			ErrorSubstring: "component ID must be set",
		},
		{
			Name: "with disallowed character in ID",
			Component: Component{
				ID: "#banana",
			},
			ErrorSubstring: "pattern",
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.Component.Validate()
			if tt.ErrorSubstring != "" {
				assert.ErrorContains(t, err, tt.ErrorSubstring)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
