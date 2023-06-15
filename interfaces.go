package portfolio

import (
	"context"

	"github.com/portfoliotree/portfolio/returns"
)

type ComponentReturnsProvider interface {
	ComponentReturnsList(ctx context.Context, component Component) (returns.List, error)
	ComponentReturnsTable(ctx context.Context, component ...Component) (returns.Table, error)
}
