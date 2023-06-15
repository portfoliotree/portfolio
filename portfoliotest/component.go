package portfoliotest

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"path"

	"github.com/portfoliotree/portfolio"
	"github.com/portfoliotree/portfolio/returns"
)

// ComponentReturnsProvider decodes the contents of data to be used in tests
// note the data time range may be updated periodically.
func ComponentReturnsProvider() portfolio.ComponentReturnsProvider {
	return crp{}
}

type crp struct{}

//go:embed testdata
var data embed.FS

// ComponentReturnsList implements portfolio.ComponentReturnsProvider
func (crp) ComponentReturnsList(_ context.Context, component portfolio.Component) (returns.List, error) {
	buf, err := fs.ReadFile(data, path.Join("testdata", "returns", component.ID+".json"))
	if err != nil {
		return nil, err
	}
	var list returns.List
	err = json.Unmarshal(buf, &list)
	return list, err
}

// ComponentReturnsTable implements portfolio.ComponentReturnsProvider
func (di crp) ComponentReturnsTable(ctx context.Context, components ...portfolio.Component) (returns.Table, error) {
	var table returns.Table
	for _, component := range components {
		list, err := di.ComponentReturnsList(ctx, component)
		if err != nil {
			return returns.Table{}, err
		}
		table = table.AddColumn(list)
	}
	return table, nil
}
