package catalog

import (
	"context"
	"log"

	"github.com/NuZard84/go_microservices/catalog/proto/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Printf("failed to connect to catalog service: %v", err)
		return nil, err
	}

	return &Client{
		conn:    conn,
		service: pb.NewCatalogServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) PostProduct(ctx context.Context, name string, description string, price float64) (*Product, error) {
	req := &pb.PostProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
	}

	res, err := c.service.PostProduct(ctx, req)
	if err != nil {
		log.Printf("failed to post product: %v", err)
		return nil, err
	}

	return &Product{
		ID:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	req := &pb.GetProductRequest{
		Id: id,
	}

	res, err := c.service.GetProduct(ctx, req)
	if err != nil {
		log.Printf("failed to get product: %v", err)
		return nil, err
	}

	return &Product{
		ID:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, skip uint64, take uint64, ids []string, query string) ([]*Product, error) {
	req := &pb.GetProductsRequest{
		Skip:  skip,
		Take:  take,
		Ids:   ids,
		Query: query,
	}

	res, err := c.service.GetProducts(ctx, req)
	if err != nil {
		log.Printf("failed to get products: %v", err)
		return nil, err
	}

	products := []*Product{}
	for _, p := range res.Products {
		products = append(products, &Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	return products, nil
}
