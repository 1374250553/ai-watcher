.PHONY: build run clean test

# 编译
build:
	go build -o ai-watcher cmd/main.go

# 运行
run: build
	./ai-watcher

# 清理
clean:
	rm -f ai-watcher

# 测试
test:
	go test ./...

# 格式化
fmt:
	go fmt ./...

# 获取依赖
deps:
	go mod tidy

# 初始化数据库
init-db:
	mysql -u root -p < migrations/001_init.sql

# 启动服务
start:
	sudo systemctl start ai-watcher

# 停止服务
stop:
	sudo systemctl stop ai-watcher

# 重启服务
restart:
	sudo systemctl restart ai-watcher

# 查看状态
status:
	sudo systemctl status ai-watcher

# 查看日志
logs:
	sudo journalctl -u ai-watcher -f
