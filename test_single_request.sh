#!/bin/bash

# 检查是否安装了jq
if ! command -v jq &>/dev/null; then
  echo "错误: 请先安装 jq 工具"
  echo "可以使用以下命令安装:"
  echo "  Mac: brew install jq"
  echo "  Ubuntu/Debian: sudo apt-get install jq"
  exit 1
fi

# 检查参数
if [ "$#" -lt 1 ]; then
  echo "用法: $0 <请求类型> [参数]"
  echo "请求类型:"
  echo "  init              - 初始化请求"
  echo "  list              - 获取工具列表"
  echo "  list_channels     - 列出频道"
  echo "  thread_replies    - 获取消息线程回复"
  echo ""
  echo "示例:"
  echo "  $0 init"
  echo "  $0 list"
  echo "  $0 list_channels"
  echo "  $0 thread_replies"
  exit 1
fi

# 创建日志目录
mkdir -p logs

# 使用日期作为日志文件名
log_date=$(date +"%Y%m%d")
log_file="logs/requests_${log_date}.log"

# 添加时间戳到日志条目
current_time=$(date +"%H:%M:%S")
echo -e "\n[${current_time}] ====== 新请求 =====" | tee -a "$log_file"

request_type=$1
shift

case "$request_type" in
init)
  echo "发送初始化请求..." | tee -a "$log_file"
  request='{
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {
                "protocolVersion": "2024-11-05",
                "clientInfo": {
                    "name": "single-request-client",
                    "version": "1.0.0"
                },
                "capabilities": {}
            }
        }'
  ;;

list)
  echo "发送工具列表请求..." | tee -a "$log_file"
  request='{
            "jsonrpc": "2.0",
            "id": 2,
            "method": "tools/list"
        }'
  ;;

list_channels)
  echo "发送列出频道请求..." | tee -a "$log_file"
  request='{
            "jsonrpc": "2.0",
            "id": 3,
            "method": "tools/call",
            "params": {
                "name": "slack_list_channels",
                "arguments": {
                    "limit": 100
                }
            }
        }'
  ;;

thread_replies)
  echo -n "请输入Slack消息URL (例如: https://workspace.slack.com/archives/C0734812MFG/p1742788004223029): " | tee -a "$log_file"
  read -r thread_url
  if [ -z "$thread_url" ]; then
    echo "错误: 未提供URL" | tee -a "$log_file"
    exit 1
  fi

  echo "发送获取线程回复请求..." | tee -a "$log_file"
  echo "消息URL: $thread_url" | tee -a "$log_file"

  request=$(
    cat <<EOF
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "tools/call",
  "params": {
    "name": "slack_get_thread_replies",
    "arguments": {
      "thread_url": "$thread_url"
    }
  }
}
EOF
  )
  ;;

*)
  echo "错误: 未知的请求类型 '$request_type'" | tee -a "$log_file"
  exit 1
  ;;
esac

echo "请求内容:" | tee -a "$log_file"
echo "$request" | jq '.' | tee -a "$log_file"
echo "" | tee -a "$log_file"

echo "响应内容:" | tee -a "$log_file"

# 发送请求并获取响应
echo "$request" | env $(cat local.env | egrep -v '^#' | xargs) go run ./main/main.go | tee -a "$log_file" | jq '.'

echo "" | tee -a "$log_file"
echo "完整日志已保存到: $log_file"
