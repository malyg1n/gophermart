package http

import (
	"encoding/json"
	"gophermart/model"
	"gophermart/pkg/errs"
	model2 "gophermart/provider/accrual/model"
	"net/http"
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
	client := &http.Client{}
	resp, err := client.Get(p.addr + "/api/orders/" + orderID)
	defer func() {
		resp.Body.Close()
	}()

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusTooManyRequests {
			return nil, errs.ErrToManyRequests
		}

		return nil, errs.ErrAccrualResponse
	}

	dec := json.NewDecoder(resp.Body)
	var o model2.Order
	err = dec.Decode(&o)
	if err != nil {
		return nil, err
	}

	return o.ToCanonical(), err
}
