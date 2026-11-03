package svc

import (
	"myuser_rpc/internal/config"
	"myuser_rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	DB     sqlx.SqlConn
	Redis  *redis.Redis
	Users  *model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  c.Redis.NewRedis(),
		Users:  model.NewUserModel(db),
	}
}
