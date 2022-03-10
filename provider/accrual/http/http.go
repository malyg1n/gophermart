package http

import (
	"encoding/json"
	"github.com/hashicorp/go-retryablehttp"
	"gophermart/model"
	baseModel "gophermart/provider/accrual/model"
)

// Provider struct.
type Provider struct {
	addr string
}

// NewAccrualHTTPProvider struct.
func NewAccrualHTTPProvider(addr string) Provider {
	return Provider{
		addr: addr,
	}
}

// CheckOrder in accrual system.
func (p Provider) CheckOrder(orderID string) (*model.Order, error) {
	client := retryablehttp.NewClient()
	client.RetryMax = 5
	resp, err := client.Get(p.addr + "/api/orders/" + orderID)

	defer func() {
		resp.Body.Close()
	}()

	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)
	var o baseModel.Order
	err = dec.Decode(&o)
	if err != nil {
		return nil, err
	}

	return o.ToCanonical(), err
}
