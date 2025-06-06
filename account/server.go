// protoc --go_out=./proto --go-grpc_out=./proto account.proto
package account

import (
	"context"
	"fmt"
	"net"

	"github.com/NuZard84/go_microservices/account/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service Service
}

func ListenGRPC(service Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	serv := grpc.NewServer()

	pb.RegisterAccountServiceServer(serv, &grpcServer{service: service})
	reflection.Register(serv)

	if err := serv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

func (s *grpcServer) PostAccount(ctx context.Context, req *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	a, err := s.service.PostAccount(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return &pb.PostAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	a, err := s.service.GetAccount(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	res, err := s.service.ListAccounts(ctx, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}

	accounts := []*pb.Account{}
	for _, a := range res {
		accounts = append(accounts, &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		})
	}
	return &pb.GetAccountsResponse{
		Accounts: accounts,
	}, nil
}
