.PHONY: all build run clean vendor

# 默认目标
all: build

# 创建必要的目录
bin:
	mkdir -p bin

# 更新 vendor 目录
vendor:
	@echo "更新 vendor 目录..."
	go mod vendor
	@echo "vendor 目录更新完成！"

# 构建目标
build: bin
	@echo "更新 vendor 目录..."
	go mod vendor
	@echo "构建 slack-mcp 二进制文件..."
	go build -o bin/slack-mcp ./main
	@echo "构建成功！二进制文件位于 bin/slack-mcp"

# 运行目标
run:
	env $$(cat local.env | egrep -v '^#' | xargs) \
		go run ./main/main.go

# 清理目标binary
clean:
	rm -rf bin/

# 清理测试logs
clean-logs:
	rm -rf logs/

# ----------------------------
# Development utility
# ----------------------------
go-lint:
	(cd ./main && golangci-lint run --config ../.golangci.toml)
