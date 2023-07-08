package metadata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"movix/metadata/pkg/model"
)

var ErrNotFound = errors.New("not found")

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr: addr}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", g.addr+"/metadata", nil)
	if err != nil {
		return nil, err
	}

	values := req.URL.Query()

	values.Add("id", id)

	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}

	var metadata model.Metadata

	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}
