package catalog

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/NuZard84/go_microservices/catalog/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedCatalogServiceServer
	service Service
}

func ListenGRPC(service Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	serv := grpc.NewServer()

	pb.RegisterCatalogServiceServer(serv, &grpcServer{service: service})
	reflection.Register(serv)

	if err := serv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

func (s *grpcServer) PostProduct(ctx context.Context, req *pb.PostProductRequest) (*pb.PostProductResponse, error) {

	p, err := s.service.PostProduct(ctx, req.Name, req.Description, req.Price)
	if err != nil {
		log.Printf("failed to post product: %v", err)
		return nil, err
	}

	return &pb.PostProductResponse{
		Product: &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.service.GetProduct(ctx, req.Id)
	if err != nil {
		log.Printf("failed to get product: %v", err)
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {

	var res []Product
	var err error

	if req.Query != "" {
		res, err = s.service.SearchProducts(ctx, req.Query, req.Skip, req.Take)
	} else if len(req.Ids) > 0 {
		res, err = s.service.GetProductByIDs(ctx, req.Ids)
	} else {
		res, err = s.service.GetProducts(ctx, req.Skip, req.Take)
	}

	if err != nil {
		log.Printf("failed to get products: %v", err)
		return nil, err
	}

	productsRes := []*pb.Product{}

	for _, p := range res {
		productsRes = append(productsRes, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	return &pb.GetProductsResponse{
		Products: productsRes,
	}, nil
}
