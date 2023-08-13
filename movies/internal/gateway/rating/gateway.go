package rating

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"movix/movies/internal/gateway"
	"movix/pkg/discovery"
	"movix/rating/pkg/model"
)

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

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	var addrs []string
	var err error

	if len(g.addrs) == 0 {
		addrs, err = g.registry.ServiceAddresses(ctx, "metadata")
		if err != nil {
			return 0, err
		}
	} else {
		log.Println("no service discovery, use fixed addresses", g.addrs)
		addrs = g.addrs
	}

	endpointUrl := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"

	log.Println("send request to rating service", endpointUrl)

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
		log.Println("rating not found for", recordType, "with id", recordID)
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		log.Println("rating service error: ", resp.StatusCode)
		return 0, fmt.Errorf("non-2xx response: %v", resp)
	}

	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}

	return v, nil
}

func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, userID model.UserID, ratingValue model.RatingValue) error {
	var addrs []string
	var err error

	if len(g.addrs) == 0 {
		addrs, err = g.registry.ServiceAddresses(ctx, "metadata")
		if err != nil {
			return err
		}
	} else {
		log.Println("no service discovery, use fixed addresses", g.addrs)
		addrs = g.addrs
	}

	endpointUrl := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"

	req, err := http.NewRequestWithContext(ctx, "PUT", endpointUrl, nil)
	if err != nil {
		return err
	}

	values := req.URL.Query()

	values.Set("id", string(recordID))
	values.Set("type", fmt.Sprintf("%v", recordType))
	values.Set("value", fmt.Sprintf("%v", ratingValue))
	values.Set("userID", string(userID))

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
