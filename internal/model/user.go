package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type User struct {
	ID           int64  `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
}

type UserModel struct {
	conn sqlx.SqlConn
}

func NewUserModel(conn sqlx.SqlConn) *UserModel {
	return &UserModel{conn: conn}
}

func (m *UserModel) Insert(ctx context.Context, username, passwordHash string) (int64, error) {
	result, err := m.conn.ExecCtx(ctx,
		"INSERT INTO users (username, password_hash) VALUES (?, ?)",
		username, passwordHash,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (m *UserModel) FindByUsername(ctx context.Context, username string) (*User, error) {
	var record User
	err := m.conn.QueryRowCtx(ctx, &record,
		"SELECT id, username, password_hash FROM users WHERE username = ? LIMIT 1",
		username,
	)
	if err != nil {
		return nil, err
	}

	return &record, nil
}
