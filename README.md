# openim demo
opemim的简易测试demo

## openim-docker
openim 依赖服务 docker一键部署启动

- 启动服务：
```bash
docker compose up -d
```

- 停止服务：
```bash
docker compose down
```

- 查看日志：
```bash
docker logs -f openim-server
docker logs -f openim-chat
```

## 创建用户

#### openim后台
http://127.0.0.1:11002
账号密码都是chatAdmin
也可以通过这个后台创建用户

#### openim web
http://127.0.0.1:11001
直接注册用户，验证吗注册，默认验证码为666666


## frontend
简易的前端聊天页面，在后台找到用户UserID类似：4319292610，输入后点击登陆

```ts
npm install
npm run dev
```

## backend
golang http服务，用来获取token
```go
go mod tidy

go run main.go
```


## test
token生成测试main.go
```
go mod tidy

go run main.go
```