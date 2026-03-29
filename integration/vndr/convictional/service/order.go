package service

import (
	"encoding/json"
	"fmt"

	dto "github.com/harishash/dotshop-be/integration/vndr/convictional/buyer/dto"
	api "github.com/harishash/dotshop-be/integration/vndr/convictional/client"
)

type IOrderService interface {
	GetAllOrders() (dto.Order, error)
	// GetOrderByID(id string) (dto.Order, error)
	// DeleteOrder(id string) (dto.Order, error)
	//UpdateOrder(order dto.Order) (dto.Order, error)
}

type OrderService struct {
	client *api.APIClient
}

func NewOrderService() *OrderService {
	return &OrderService{
		client: api.Client,
	}
}

func (o *OrderService) GetAllOrders() (dto.Order, error) {
	allOrdersResponse, err := o.client.GetAll("orders")
	if err != nil {
		return dto.Order{}, fmt.Errorf("error getting all Orders: %v", err)
	}

	var allOrders dto.Order
	err = json.Unmarshal(allOrdersResponse, &allOrders)
	if err != nil {
		return dto.Order{}, fmt.Errorf("error parsing all allOrders response: %v", err)
	}

	return allOrders, nil
}
