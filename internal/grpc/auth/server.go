package auth

import (
	"context"

	//"context"
	"fmt"

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
	IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
}

type serverAPI struct {
	Autorization_servise.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) { //регистрирует обработчик
	Autorization_servise.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *Autorization_servise.LoginRequest) (*Autorization_servise.LoginResponse, error) {
	// 1. Валидация Email
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// 2. Валидация Пароля
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// 3. Валидация AppId (для int32 дефолтное значение — 0)
	if req.GetAppId() == nil {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	// Получаем appId напрямую как int32 (никаких uuid.Parse больше не нужно!)
	appId := req.GetAppId()

	// 4. Передаем appId (int32) в сервис авторизации
	// Убедись, что метод s.auth.Login на бэкенде теперь тоже принимает int32, а не uuid.UUID
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), appId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error: "+err.Error())
	}

	return &Autorization_servise.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *Autorization_servise.RegisterRequest) (*Autorization_servise.RegisterResponse, error) {
	//переделать с пакетом и в мидлвейр валидацию
	if s == nil {
		return nil, status.Error(codes.Internal, "server API is not initialized")
	}

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is requred")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is requred")
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		//todo - user exist
		return nil, status.Error(codes.Internal, "internal error: "+err.Error())
	}

	return &Autorization_servise.RegisterResponse{
		UserId: userID[:], // Возвращаем массив байтов, как определено в proto
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *Autorization_servise.IsAdminRequest) (*Autorization_servise.IsAdminResponse, error) {
	if s == nil {
		return nil, status.Error(codes.Internal, "server API is not initialized")
	}

	if len(req.GetUserId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is requred")
	}

	userIdStr := string(req.GetUserId())
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("failed to parse user id: %v, received length: %d, value: %s", err, len(userIdStr), userIdStr))
	}

	isAdmin, err := s.auth.IsAdmin(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error: "+err.Error())
	}

	return &Autorization_servise.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
