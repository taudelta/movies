package rating

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"movix/pkg/discovery"
	"movix/rating/pkg/model"
)

var ErrNotFound = errors.New("not found")

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{
		registry: registry,
	}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return 0, err
	}

	endpointUrl := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"

	req, err := http.NewRequestWithContext(ctx, "GET", endpointUrl, nil)
	if err != nil {
		return 0, err
	}

	values := req.URL.Query()

	values.Set("id", string(recordID))
	values.Set("type", fmt.Sprintf("%v", recordType))

	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx response: %v", resp)
	}

	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}

	return v, nil
}

func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return err
	}

	endpointUrl := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"

	req, err := http.NewRequestWithContext(ctx, "PUT", endpointUrl, nil)
	if err != nil {
		return err
	}

	values := req.URL.Query()

	values.Set("id", string(recordID))
	values.Set("type", fmt.Sprintf("%v", recordType))
	values.Set("value", fmt.Sprintf("%v", rating.Value))
	values.Set("userID", string(rating.UserID))

	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", resp)
	}

	return nil
}
