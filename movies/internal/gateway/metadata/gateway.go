package metadata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"movix/metadata/pkg/model"
	"movix/movies/internal/gateway"
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
	var addrs []string
	var err error

	if len(g.addrs) == 0 {
		addrs, err = g.registry.ServiceAddresses(ctx, "metadata")
		if err != nil {
			return nil, err
		}
	} else {
		addrs = g.addrs
	}

	endpointUrl := "http://" + addrs[rand.Intn(len(addrs))] + "/metadata"

	log.Println("send request to metadata service", endpointUrl)

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
		return nil, gateway.ErrNotFound
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
