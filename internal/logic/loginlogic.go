package logic

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"strconv"

	"myuser_rpc/internal/svc"
	"myuser_rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const sessionTTL = 24 * 60 * 60

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	if !validUsername(in.GetUsername()) || !validPassword(in.GetPassword()) {
		return nil, status.Error(codes.InvalidArgument, "invalid username or password")
	}

	record, err := l.svcCtx.Users.FindByUsername(l.ctx, in.GetUsername())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.Unauthenticated, "invalid username or password")
	}
	if err != nil {
		l.Errorf("find user for login: %v", err)
		return nil, status.Error(codes.Internal, "failed to log in")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(record.PasswordHash), []byte(in.GetPassword())); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid username or password")
	}

	tokenBytes := make([]byte, 32)
	if _, err = io.ReadFull(rand.Reader, tokenBytes); err != nil {
		l.Errorf("generate access token: %v", err)
		return nil, status.Error(codes.Internal, "failed to log in")
	}
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	tokenHash := sha256.Sum256([]byte(token))
	key := "user:session:" + hex.EncodeToString(tokenHash[:])
	if err = l.svcCtx.Redis.SetexCtx(l.ctx, key, strconv.FormatInt(record.ID, 10), sessionTTL); err != nil {
		l.Errorf("store access token: %v", err)
		return nil, status.Error(codes.Internal, "failed to log in")
	}

	return &user.LoginResponse{UserId: record.ID, AccessToken: token}, nil
}
