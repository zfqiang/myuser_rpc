# 用户 RPC 服务

基于 go-zero 的 gRPC 服务，提供用户注册和登录。项目使用 Go 1.27.3、MySQL 保存用户数据、bcrypt 哈希密码，并将登录 token 的 SHA-256 摘要存入 Redis（默认有效期 24 小时）。

## 初始化

1. 确认 Go 版本为 `1.27.3`，并启动 MySQL 和 Redis。
2. 执行 `sql/schema.sql` 创建数据库和用户表。
3. 按本机环境修改 `etc/user.yaml` 中的 MySQL DSN 和 Redis 地址。
4. 在 `user` 目录启动服务：

```powershell
go mod tidy
go run . -f etc/user.yaml
```

RPC 默认监听 `0.0.0.0:8080`。开发/测试模式下启用 gRPC reflection。

## RPC

- `Register(username, password)`：用户名限定为 3-32 位 ASCII 字母、数字或下划线；密码限定为 8-72 字节。重复用户名返回 `ALREADY_EXISTS`。
- `Login(username, password)`：凭据错误统一返回 `UNAUTHENTICATED`；成功响应包含 `user_id` 和 `access_token`。

服务端 Redis 键为 `user:session:<token 的 SHA-256>`，值为用户 ID。不要将明文密码或 token 写入日志；生产环境应通过安全的配置注入机制提供数据库密码和 Redis 凭据。

## 重新生成 RPC 代码

修改 `user.proto` 后，在 `user` 目录执行：

```powershell
goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=. --module myuser_rpc
```
