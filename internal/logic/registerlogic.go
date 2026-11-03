package logic

import (
	"context"
	"errors"

	"myuser_rpc/internal/svc"
	"myuser_rpc/user"

	"github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	if !validUsername(in.GetUsername()) {
		return nil, status.Error(codes.InvalidArgument, "username must be 3-32 letters, digits, or underscores")
	}
	if !validPassword(in.GetPassword()) {
		return nil, status.Error(codes.InvalidArgument, "password must be 8-72 bytes")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		l.Errorf("hash password: %v", err)
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	id, err := l.svcCtx.Users.Insert(l.ctx, in.GetUsername(), string(passwordHash))
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, status.Error(codes.AlreadyExists, "username already exists")
		}
		l.Errorf("insert user: %v", err)
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &user.RegisterResponse{UserId: id}, nil
}
