package accrual

import (
	"encoding/json"
	"gophermart/model"
	"gophermart/pkg/errs"
	"net/http"
)

// HttpProvider struct.
type HttpProvider struct {
	addr string
}

// NewAccrualHttpProvider struct.
func NewAccrualHttpProvider(addr string) HttpProvider {
	return HttpProvider{
		addr: addr,
	}
}

// CheckOrder in accrual system.
func (p HttpProvider) CheckOrder(orderID string) (*model.Order, error) {
	client := &http.Client{}
	resp, err := client.Get(p.addr + "/api/order/" + orderID)
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
	var o Order
	err = dec.Decode(&o)
	if err != nil {
		return nil, err
	}

	return o.ToCanonical(), err
}
