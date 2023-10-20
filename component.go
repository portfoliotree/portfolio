package portfolio

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Component struct {
	Type  string `yaml:"type,omitempty"  json:"type,omitempty"  bson:"type"`
	ID    string `yaml:"id,omitempty"    json:"id,omitempty"    bson:"id"`
	Label string `yaml:"label,omitempty" json:"label,omitempty" bson:"label"`
}

var componentExpression = regexp.MustCompile(`^[a-zA-Z0-9.:s]{1,24}$`)

func (component Component) Validate() error {
	switch component.ID {
	case "":
		return fmt.Errorf(`component ID must be set`)
	case "undefined":
		return fmt.Errorf(`component ID must not be "undefined"`)
	}
	if !componentExpression.MatchString(component.ID) {
		return fmt.Errorf("component id %q does not match the component ID pattern %q", component.ID, componentExpression.String())
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
