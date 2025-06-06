package order

import (
	"context"
	"time"
)

type Service interface {
	PostOrder(ctx context.Context, order Order) error
	GetOrderForAccount(ctx context.Context, accountId string) ([]Order, error)
}
type Order struct {
	ID         string           `json:"id"`
	CreatedAt  time.Time        `json:"created_at"`
	AccountID  string           `json:"account_id"`
	TotalPrice float64          `json:"total_price"`
	Products   []OrderedProduct `json:"products"`
}

type orderService struct {
	repo Repository
}

func NewOrderService(repo Repository) Service {
	return &orderService{repo: repo}
}

func (s *orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {

}

func (s *orderService) GetOrderForAccount(ctx context.Context, accountId string) ([]Order, error) {

}

type OrderedProduct struct {
	ID          string  `json:"id"`
	Quantity    int     `json:"quantity"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}
