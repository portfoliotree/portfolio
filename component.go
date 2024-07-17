package portfolio

import (
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	ComponentTypeSecurity   = "Security"
	ComponentTypePortfolio  = "Portfolio"
	ComponentTypeEquity     = "Equity"
	ComponentTypeETF        = "ETF"
	ComponentTypeFactor     = "Factor"
	ComponentTypeMutualFund = "Mutual Fund"
)

func ComponentTypes() []string {
	return []string{
		ComponentTypeSecurity,
		ComponentTypePortfolio,
		ComponentTypeEquity,
		ComponentTypeETF,
		ComponentTypeFactor,
		ComponentTypeMutualFund,
	}
}

type Component struct {
	Type  string `yaml:"type,omitempty"  json:"type,omitempty"  bson:"type"`
	ID    string `yaml:"id,omitempty"    json:"id,omitempty"    bson:"id"`
	Label string `yaml:"label,omitempty" json:"label,omitempty" bson:"label"`
}

var componentExpression = regexp.MustCompile(`^[a-zA-Z0-9.:]{1,24}$`)

func (component *Component) Validate() error {
	if component.ID == "" {
		return fmt.Errorf(`component ID must be set`)
	}
	if component.ID == "undefined" {
		return fmt.Errorf(`component ID must not be "undefined"`)
	}
	if !componentExpression.MatchString(component.ID) {
		return fmt.Errorf("component id %q does not match the component ID pattern %q", component.ID, componentExpression.String())
	}
	if component.Type != "" && !slices.Contains(ComponentTypes(), component.Type) {
		return fmt.Errorf("component type %q is not a known component type", component.Type)
	}
	return nil
}

func (component *Component) marshalURLValues(q url.Values, prefix string) {
	q.Add(strings.Join([]string{prefix, "id"}, "-"), component.ID)
}

func (component *Component) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		return value.Decode(&component.ID)
	case yaml.MappingNode:
		type C Component
		var c C
		if err := value.Decode(&c); err != nil {
			return err
		}
		*component = Component(c)
		return nil
	default:
		return fmt.Errorf("wrong YAML type: expected either a component identifier (string) or a Component")
	}
}
