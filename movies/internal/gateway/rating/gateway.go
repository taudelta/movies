package rating

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"movix/rating/pkg/model"
)

var ErrNotFound = errors.New("not found")

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr: addr}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", g.addr+"/rating", nil)
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
	req, err := http.NewRequestWithContext(ctx, "PUT", g.addr+"/rating", nil)
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
