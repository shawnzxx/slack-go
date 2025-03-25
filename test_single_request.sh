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
  echo "  post_message      - 发送消息 (需要提供频道ID和消息内容)"
  echo "  get_thread_replies - 获取消息线程回复 (需要提供Slack消息URL)"
  echo ""
  echo "示例:"
  echo "  $0 init"
  echo "  $0 list"
  echo "  $0 list_channels"
  echo "  $0 post_message general \"Hello World!\""
  echo "  $0 get_thread_replies \"https://your-workspace.slack.com/archives/C0734812MFG/p1742788004223029\""
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

get_thread_replies)
  thread_url="$1"

  if [ -z "$thread_url" ]; then
    echo -n "请输入Slack消息URL (例如: https://workspace.slack.com/archives/C0734812MFG/p1742788004223029): " | tee -a "$log_file"
    read -r thread_url
    if [ -z "$thread_url" ]; then
      echo "错误: 未提供URL" | tee -a "$log_file"
      exit 1
    fi
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

post_message)
  if [ "$#" -lt 2 ]; then
    echo "错误: 发送消息需要提供频道ID和消息内容" | tee -a "$log_file"
    exit 1
  fi

  channel="$1"
  message="$2"

  echo "发送消息请求..." | tee -a "$log_file"
  echo "频道: $channel" | tee -a "$log_file"
  echo "消息: $message" | tee -a "$log_file"

  request=$(
    cat <<EOF
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "slack_post_message",
    "arguments": {
      "channel": "$channel",
      "text": "$message"
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
# 发送请求到已运行的服务
if [ ! -p "/tmp/slack-mcp-in" ]; then
  echo "错误: 请先启动服务器 (make run-local)" | tee -a "$log_file"
  exit 1
fi

# 发送请求
echo "$request" >/tmp/slack-mcp-in

# 读取一个响应后自动退出
read -r response </tmp/slack-mcp-out
if [[ "$response" =~ ^\{ ]]; then
  echo "$response" | jq '.' | tee -a "$log_file"
else
  echo "$response" | tee -a "$log_file"
fi

echo "" | tee -a "$log_file"
echo "完整日志已保存到: $log_file"
