package apiclient

import (
	"fmt"
	"strings"

	"github.com/Enclave-Markets/enclave-go/models"
)

func (client *ApiClient) AddPerpsOrder(req models.AddOrderReq) (*models.GenericResponse[models.ApiOrder], error) {
	path := models.V1PerpsOrdersPath

	res, err := NewHttpJsonClient[models.AddOrderReq, models.GenericResponse[models.ApiOrder]](client.ApiEndpoint + path).
		WithTimeout(client.timeout).
		SetHeaders(client.getHeaders("POST", path, req)).
		Post(req)
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

	res, err := NewHttpJsonClient[models.BatchAddOrderReq, models.GenericResponse[models.BatchAddOrderRes]](client.ApiEndpoint + path).
		WithTimeout(client.timeout).
		SetHeaders(client.getHeaders("POST", path, req)).
		Post(req)
	if err != nil {
		return res, fmt.Errorf("error with http req in perps batch order: %w", err)
	}
	if !res.Success {
		return res, fmt.Errorf("error in perps batch order %v: %v", req, res.Error)
	}

	return res, err
}

func (client *ApiClient) GetPerpOrder(orderId models.OrderID) (*models.GenericResponse[models.ApiOrder], error) {
	path := models.V1PerpsOrdersPath + "/" + string(orderId)

	res, err := NewHttpJsonClient[any, models.GenericResponse[models.ApiOrder]](client.ApiEndpoint + path).
		WithTimeout(client.timeout).
		SetHeaders(client.getHeaders("GET", path, nil)).
		Get(nil)
	if err != nil {
		return nil, fmt.Errorf("error in http req spot get order: %w", err)
	}
	if !res.Success {
		return res, fmt.Errorf("bad request perps get order %s: %v", orderId, res.Error)
	}

	return res, nil
}

func (client *ApiClient) CancelPerpsOrdersByClientId(clientIds []models.ClientOrderID) (*models.GenericResponse[models.BatchCancelRes], error) {

	ids := make([]string, 0, len(clientIds))
	for _, id := range clientIds {
		ids = append(ids, "client:"+string(id))
	}

	path := models.V1PerpsBatchOrdersPath + "?orderIDs=" + strings.Join(ids, ",")
	res, err := NewHttpJsonClient[any, models.GenericResponse[models.BatchCancelRes]](client.ApiEndpoint + path).
		WithTimeout(client.timeout).
		SetHeaders(client.getHeaders("DELETE", path, nil)).
		Delete(nil)
	if err != nil {
		return res, fmt.Errorf("error in http req perps delete batch: %w", err)
	}
	if !res.Success {
		return res, fmt.Errorf("bad request perps delete batch: %v", res.Error)
	}

	return res, nil
}

func (client *ApiClient) CancelAllPerpsOrders() error {
	path := models.V1PerpsOrdersPath

	res, err := NewHttpJsonClient[any, models.GenericResponse[any]](client.ApiEndpoint + path).
		WithTimeout(client.timeout).
		SetHeaders(client.getHeaders("DELETE", path, nil)).
		Delete(nil)
	if err != nil {
		return fmt.Errorf("error in http req perps delete all orders: %w", err)
	}
	if !res.Success {
		return fmt.Errorf("bad request spot delete all orders: %v", res.Error)
	}

	return nil
}

func (client *ApiClient) CancelAllPerpsOrdersOnMarket(market models.Market) error {
	path := models.V1PerpsOrdersPath + "?market=" + string(market)

	res, err := NewHttpJsonClient[any, models.GenericResponse[any]](client.ApiEndpoint + path).
		WithTimeout(client.timeout).
		SetHeaders(client.getHeaders("DELETE", path, nil)).
		Delete(nil)
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

	res, err := NewHttpJsonClient[any, models.GenericResponse[[]models.PerpsContract]](client.ApiEndpoint + path).
		WithTimeout(client.timeout).
		SetHeaders(client.getHeaders("GET", path, nil)).
		Get(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute perps contracts request: %w", err)
	}

	if !res.Success {
		return res, fmt.Errorf("failed to get perps contracts: %v", res.Error)
	}

	return res, nil
}

func (client *ApiClient) GetPerpsDepthBook(market models.Market) (*models.GenericResponse[models.BookSnapshot], error) {
	path := models.V1PerpsDepthPath + "?market=" + string(market)

	res, err := NewHttpJsonClient[any, models.GenericResponse[models.BookSnapshot]](client.ApiEndpoint + path).
		WithTimeout(client.timeout).
		SetHeaders(client.getHeaders("GET", path, nil)).
		Get(nil)
	if err != nil {
		return nil, fmt.Errorf("error in http req perps get depth book: %w", err)
	}
	if !res.Success {
		return res, fmt.Errorf("bad request Spot perps depth book: %v", res.Error)
	}

	return res, nil
}

func (client *ApiClient) GetPerpsFillsByOrderID(orderID models.OrderID) (*models.GenericResponse[[]models.ApiFill], error) {
	path := models.V1PerpsOrdersPath + "/" + string(orderID) + "/fills"

	res, err := NewHttpJsonClient[any, models.GenericResponse[[]models.ApiFill]](client.ApiEndpoint + path).
		WithTimeout(client.timeout).
		SetHeaders(client.getHeaders("GET", path, nil)).
		Get(nil)
	if err != nil {
		return nil, fmt.Errorf("error in http req spot get fill by order ID: %w", err)
	}

	if !res.Success {
		return res, fmt.Errorf("bad request spot fill by order id %s: %v", orderID, res.Error)
	}

	return res, err
}

func (client *ApiClient) GetPerpsFillsByClientOrderID(orderID models.OrderID) (*models.GenericResponse[[]models.ApiFill], error) {
	path := models.V1PerpsOrdersPath + "/" + models.V1SpotClientOrderIDPrefix + string(orderID) + "/fills"

	res, err := NewHttpJsonClient[any, models.GenericResponse[[]models.ApiFill]](client.ApiEndpoint + path).
		WithTimeout(client.timeout).
		SetHeaders(client.getHeaders("GET", path, nil)).
		Get(nil)
	if err != nil {
		return nil, fmt.Errorf("error in http req spot get fill by client order ID: %w", err)
	}

	if !res.Success {
		return res, fmt.Errorf("bad request spot fill by client order id %s: %v", orderID, res.Error)
	}

	return res, err
}

func (client *ApiClient) GetPerpsFills(params models.FillParams) (*models.V1PageRes[models.ApiFill], error) {
	path := models.V1PerpsFillsPath
	path += params.GetFillPathParams()

	res, err := NewHttpJsonClient[any, models.V1PageRes[models.ApiFill]](client.ApiEndpoint + path).
		WithTimeout(client.timeout).
		SetHeaders(client.getHeaders("GET", path, nil)).
		Get(nil)
	if err != nil {
		return nil, fmt.Errorf("error in http req perps get fills: %w", err)
	}

	return res, err
}
