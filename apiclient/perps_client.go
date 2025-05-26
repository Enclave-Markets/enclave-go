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
		return res, fmt.Errorf("error with http req in perp add order: %w", err)
	}
	if !res.Success {
		return res, fmt.Errorf("error in perps add order %v: %v", req, res.Error)
	}

	return res, err
}

func (client *ApiClient) AddPerpsBatchOrders(req models.BatchAddOrderReq) (*models.GenericResponse[models.BatchAddOrderRes], error) {
	path := models.V1PerpsBatchOrdersPath

	res, err := NewHttpJsonClient[models.BatchAddOrderReq, models.GenericResponse[models.BatchAddOrderRes]](
		client.ApiEndpoint + path).SetHeaders(client.getHeaders("POST", path, req)).Post(req)
	if err != nil {
		return res, fmt.Errorf("error with http req in perps batch order: %w", err)
	}
	if !res.Success {
		return res, fmt.Errorf("error in perps batch order %v: %v", req, res.Error)
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

// GetPerpsContracts retrieves all perpetual futures contracts from the /v1/perps/contracts endpoint
func (client *ApiClient) GetPerpsContracts() (*models.GenericResponse[[]models.PerpsContract], error) {
	path := models.V1PerpsContractsPath

	res, err := NewHttpJsonClient[any, models.GenericResponse[[]models.PerpsContract]](
		client.ApiEndpoint + path).SetHeaders(client.getHeaders("GET", path, nil)).Get(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute perps contracts request: %w", err)
	}

	if !res.Success {
		return res, fmt.Errorf("failed to get perps contracts: %v", res.Error)
	}

	return res, nil
}
