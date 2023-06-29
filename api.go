package portfolio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/portfoliotree/portfolio/returns"
)

const (
	// ServerURLEnvironmentVariableName is used in testing to override the host
	//  and scheme for server calls.
	ServerURLEnvironmentVariableName = "PORTFOLIO_TREE_URL"

	// DefaultURL is the scheme and host for the API calls.
	DefaultURL = "https://portfoliotree.com"
)

func portfolioTreeURL() string {
	if val := os.Getenv(ServerURLEnvironmentVariableName); val != "" {
		return val
	}
	return DefaultURL
}

const (
	ReturnsURLPath = "/api/returns"
)

func (pf *Specification) AssetReturns(ctx context.Context) (returns.Table, error) {
	if len(pf.Assets) == 0 {
		return returns.Table{}, nil
	}
	u, err := url.Parse(portfolioTreeURL())
	if err != nil {
		return returns.Table{}, err
	}
	u.Path = ReturnsURLPath
	q := u.Query()
	for _, c := range pf.Assets {
		c.marshalURLValues(q, "asset")
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return returns.Table{}, err
	}

	return doJSONRequest[returns.Table](http.DefaultClient.Do, req)
}

func doJSONRequest[T any](do func(r *http.Request) (*http.Response, error), req *http.Request) (T, error) {
	var result T
	req.Header.Set("accept", "application/json")
	res, err := do(req)
	if err != nil {
		return result, err
	}
	defer closeAndIgnoreErrors(res.Body)
	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated:
	default:
		var message string
		if strings.HasPrefix(res.Header.Get("content-type"), "text/plain") {
			b, _ := io.ReadAll(res.Body)
			message = string(b)
		} else {
			message = fmt.Sprintf("request failed %s", res.Status)
		}
		return result, errors.New(message)
	}
	// TODO: do response header validation
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(buf, &result); err != nil {
		return result, err
	}
	return result, nil
}

func closeAndIgnoreErrors(closer io.Closer) {
	_ = closer.Close()
}

// ComponentReturnsProvider is currently used for tests.
type ComponentReturnsProvider interface {
	ComponentReturnsList(ctx context.Context, component Component) (returns.List, error)
	ComponentReturnsTable(ctx context.Context, component ...Component) (returns.Table, error)
}