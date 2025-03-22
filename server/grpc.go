package server

import (
	"context"
	"log"
	"net"
	"user-srv/domain"
	"user-srv/proto"
	"user-srv/services"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	proto.UnimplementedUserServiceServer
	service services.UserService
}

func NewGRPCServer(service services.UserService) *GRPCServer {
	return &GRPCServer{service: service}
}

func (s *GRPCServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.UserResponse, error) {
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := s.service.Create(ctx, user); err != nil {
		return nil, err
	}
	return &proto.UserResponse{
		Id:        int32(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *GRPCServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.UserResponse, error) {
	user, err := s.service.GetByID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &proto.UserResponse{
		Id:        int32(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *GRPCServer) GetAllUsers(ctx context.Context, req *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error) {
	users, err := s.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var resp []*proto.UserResponse
	for _, user := range users {
		resp = append(resp, &proto.UserResponse{
			Id:        int32(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}
	return &proto.GetAllUsersResponse{Users: resp}, nil
}

func (s *GRPCServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UserResponse, error) {
	user := &domain.User{
		ID:       int(req.Id),
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := s.service.Update(ctx, user); err != nil {
		return nil, err
	}
	return &proto.UserResponse{
		Id:        int32(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *GRPCServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	if err := s.service.Delete(ctx, int(req.Id)); err != nil {
		return nil, err
	}
	return &proto.DeleteUserResponse{}, nil
}

func (s *GRPCServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, err := s.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.LoginResponse{Token: token}, nil
}

func (s *GRPCServer) GetCurrentUser(ctx context.Context, req *proto.GetCurrentUserRequest) (*proto.UserResponse, error) {
	// Пока заглушка, авторизацию по токену позже добавим
	user, err := s.service.GetByID(ctx, 1)
	if err != nil {
		return nil, err
	}
	return &proto.UserResponse{
		Id:        int32(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func StartGRPCServer(service services.UserService, addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, NewGRPCServer(service))

	log.Printf("Starting gRPC server on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}
