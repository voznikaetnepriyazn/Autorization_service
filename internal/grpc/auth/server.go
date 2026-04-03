package auth

import (
	"context"

	"github.com/google/uuid"
	Autorization_servise "github.com/voznikaetnepriyazn/Autorization-proto/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue = 0
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID uuid.UUID, eer error)
	isAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
}

type serverAPI struct {
	Autorization_servise.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) { //регистрирует обработчик
	Autorization_servise.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *Autorization_servise.LoginRequest) (*Autorization_servise.LoginResponse, error) {
	//переделать с пакетом и в мидлвейр валидацию
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is requred")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is requred")
	}

	if req.GetAppId() == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "app_id is requred")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &Autorization_servise.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) RegisterUser(ctx context.Context, req *Autorization_servise.RegisterRequest) (*Autorization_servise.RegisterResponse, error) {
	//переделать с пакетом и в мидлвейр валидацию
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is requred")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is requred")
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		//todo - user exist
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &Autorization_servise.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *Autorization_servise.IsAdminRequest) (*Autorization_servise.IsAdminResponse, error) {
	panic("implement me")
}
