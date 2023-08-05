package metadata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"movix/metadata/pkg/model"
	"movix/pkg/discovery"
)

var ErrNotFound = errors.New("not found")

type Gateway struct {
	addrs    []string
	registry discovery.Registry
}

func New(registry discovery.Registry, addrs []string) *Gateway {
	return &Gateway{
		addrs:    addrs,
		registry: registry,
	}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}

	endpointUrl := "http://" + addrs[rand.Intn(len(addrs))] + "/metadata"

	req, err := http.NewRequestWithContext(ctx, "GET", endpointUrl, nil)
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
