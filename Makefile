.PHONY: all build run clean

# 默认目标
all: build

# 创建必要的目录
bin:
	mkdir -p bin

# 构建目标
build: bin
	@echo "Building slack-mcp binary..."
	go build -o bin/slack-mcp ./main
	@echo "Build successful! Binary located at bin/slack-mcp"

# 运行目标
run-local:
	env $$(cat local.env | egrep -v '^#' | xargs) \
		go run ./main/main.go

# 清理目标binary
clean:
	rm -rf bin/

# 清理测试logs
clean-logs:
	rm -rf logs/