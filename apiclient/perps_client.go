package apiclient

import (
	"fmt"

	"github.com/Enclave-Markets/enclave-go/models"
)

func (client *ApiClient) AddPerpsOrder(req models.AddOrderReq) (*models.GenericResponse[models.ApiOrder], error) {
	path := models.V1PerpsOrdersPath

	res, err := NewHttpJsonClient[models.AddOrderReq, models.GenericResponse[models.ApiOrder]](
		client.ApiEndpoint + path).SetHeaders(client.getHeaders("POST", path, req)).Post(req)
	if err != nil {
		return res, fmt.Errorf("error with http req in spot add order: %w", err)
	}
	if !res.Success {
		return res, fmt.Errorf("error in spot add order %v: %v", req, res.Error)
	}

	return res, err
}

func (client *ApiClient) CancelAllPerpsOrdersOnMarket(market models.Market) error {
	path := models.V1PerpsOrdersPath + "?market=" + string(market)

	res, err := NewHttpJsonClient[any, models.GenericResponse[any]](
		client.ApiEndpoint + path).SetHeaders(client.getHeaders("DELETE", path, nil)).Delete(nil)
	if err != nil {
		return fmt.Errorf("error in http req perps delete all orders: %w", err)
	}
	if !res.Success {
		return fmt.Errorf("bad request perps delete all orders: %v", res.Error)
	}

	return nil
}
