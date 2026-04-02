package auth

import (
	"context"

	Autorization_servise "github.com/voznikaetnepriyazn/Autorization-proto/generated"
	"google.golang.org/grpc"
)

type serverAPI struct {
	Autorization_servise.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) { //регистрирует обработчик
	Autorization_servise.RegisterAuthServer(gRPC, &serverAPI{})
}

// заглушки
func (s *serverAPI) Login(ctx context.Context, req *Autorization_servise.LoginRequest) (*Autorization_servise.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(ctx context.Context, req *Autorization_servise.RegisterRequest) (*Autorization_servise.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *Autorization_servise.IsAdminRequest) (*Autorization_servise.IsAdminResponse, error) {
	panic("implement me")
}
