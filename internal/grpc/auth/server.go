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
	Login(ctx context.Context, email string, password string, appID int32) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID uuid.UUID, err error)
	IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
	ValidateToken(tokenString string, secret string) (userID uuid.UUID, err error)
}

type serverAPI struct {
	Autorization_servise.UnimplementedAuthServer
	auth      Auth
	jwtSecret string
}

func Register(gRPC *grpc.Server, auth Auth, jwtSecret string) { //регистрирует обработчик
	Autorization_servise.RegisterAuthServer(gRPC, &serverAPI{auth: auth, jwtSecret: jwtSecret})
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
	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
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

func (s *serverAPI) ValidateToken(
	ctx context.Context,
	req *Autorization_servise.ValidateTokenRequest,
) (*Autorization_servise.ValidateTokenResponse, error) {
	if req.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	userUUID, err := s.auth.ValidateToken(req.GetToken(), s.jwtSecret)
	if err != nil {
		// ВАЖНО: Выводим настоящую ошибку в консоль Go бэкенда!
		fmt.Printf("[gRPC Debug] Ошибка валидации токена: %v\n", err)

		// Для фронтенда пока оставляем ошибку, но можем временно пробросить детали
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("invalid token: %v", err))
	}

	return &Autorization_servise.ValidateTokenResponse{
		UserId: userUUID.String(),
	}, nil
}
