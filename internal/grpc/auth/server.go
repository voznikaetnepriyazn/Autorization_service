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
	Login(ctx context.Context, email string, password string, appID uuid.UUID) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID uuid.UUID, err error)
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

	if req.GetAppId() == nil {
		return nil, status.Error(codes.InvalidArgument, "app_id is requred")
	}

	appId, err := uuid.ParseBytes(req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse app id")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), appId)
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
		UserId: userID[:],
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *Autorization_servise.IsAdminRequest) (*Autorization_servise.IsAdminResponse, error) {
	if len(req.GetUserId()) == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "user_id is requred")
	}

	userId, err := uuid.ParseBytes(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse user id")
	}

	isAdmin, err := s.auth.isAdmin(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &Autorization_servise.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
