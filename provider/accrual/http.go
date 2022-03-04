package accrual

import (
	"encoding/json"
	"gophermart/model"
	"gophermart/pkg/errs"
	"gophermart/pkg/logger"
	"net/http"
)

// HTTPProvider struct.
type HTTPProvider struct {
	addr string
}

// NewAccrualHTTPProvider struct.
func NewAccrualHTTPProvider(addr string) HTTPProvider {
	return HTTPProvider{
		addr: addr,
	}
}

// CheckOrder in accrual system.
func (p HTTPProvider) CheckOrder(orderID string) (*model.Order, error) {
	client := &http.Client{}
	resp, err := client.Get(p.addr + "/api/orders/" + orderID)
	defer func() {
		resp.Body.Close()
	}()

	if err != nil {
		return nil, err
	}
	logger.GetLogger().Info(p.addr + "/api/orders/" + orderID)
	logger.GetLogger().Info(resp.StatusCode)
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
