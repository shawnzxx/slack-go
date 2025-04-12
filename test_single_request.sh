#!/bin/bash

#####################################################################
# 测试脚本：用于测试 Slack MCP 服务器的各种请求
#
# 功能：
# - 发送初始化请求 (init)
# - 获取工具列表 (list)
# - 列出 Slack 频道 (list_channels)
# - 获取消息线程回复 (thread_replies)
# - 发布消息到 Slack 频道 (post_message)
# - 获取用户资料信息 (get_users_profile)
#
# 使用方法：
# 1. 确保已安装 jq 工具
# 2. 确保 local.env 文件中包含必要的环境变量
# 3. 运行命令：./test_single_request.sh <请求类型>
#
# 环境要求：
# - jq：用于 JSON 处理
# - local.env：包含 SLACK_TOKEN 和 SLACK_TEAM_ID
#
# 输出：
# - 在终端显示请求和响应
# - 同时保存到 logs 目录下的日志文件
#####################################################################

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
  echo "  post_message      - 发布消息到 Slack 频道"
  echo "  get_users_profile - 获取多个用户资料信息"
  echo ""
  echo "示例:"
  echo "  $0 init"
  echo "  $0 list"
  echo "  $0 list_channels"
  echo "  $0 thread_replies"
  echo "  $0 post_message"
  echo "  $0 get_users_profile"
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
      "protocolVersion": "0.1.0",
      "clientInfo": {
        "name": "single-request-client",
        "version": "1.0.0",
        "publisher": "shawnzhang"
      },
      "capabilities": {
        "tools": true,
        "resources": true,
        "prompts": false
      }
    }
  }'
  ;;

list_tools)
  echo "发送工具列表请求..." | tee -a "$log_file"
  request='{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list",
    "params": {}
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

  request=$(jq -n \
    --arg url "$thread_url" \
    '{
      "jsonrpc": "2.0",
      "id": 5,
      "method": "tools/call",
      "params": {
        "name": "slack_get_thread_replies",
        "arguments": {
          "thread_url": $url
        }
      }
    }')
  ;;

post_message)
  echo -n "请输入Slack频道ID: " | tee -a "$log_file"
  read -r channel_id
  if [ -z "$channel_id" ]; then
    echo "错误: 未提供频道ID" | tee -a "$log_file"
    exit 1
  fi

  echo -n "请输入消息文本: " | tee -a "$log_file"
  read -r text
  if [ -z "$text" ]; then
    echo "错误: 未提供消息文本" | tee -a "$log_file"
    exit 1
  fi

  echo "发送发布消息请求..." | tee -a "$log_file"
  echo "频道ID: $channel_id" | tee -a "$log_file"
  echo "消息文本: $text" | tee -a "$log_file"

  request=$(jq -n \
    --arg channel_id "$channel_id" \
    --arg text "$text" \
    '{
      "jsonrpc": "2.0",
      "id": 6,
      "method": "tools/call",
      "params": {
        "name": "post_message",
        "arguments": {
          "channel_id": $channel_id,
          "text": $text
        }
      }
    }')
  ;;

get_users_profile)
  echo -n "请输入用户ID (多个ID用空格分隔): " | tee -a "$log_file"
  read -r user_ids_input
  if [ -z "$user_ids_input" ]; then
    echo "错误: 未提供用户ID" | tee -a "$log_file"
    exit 1
  fi

  # 将输入转换为数组
  IFS=' ' read -r -a user_ids_array <<<"$user_ids_input"

  # 验证是否至少有一个ID
  if [ ${#user_ids_array[@]} -eq 0 ]; then
    echo "错误: 至少需要提供一个用户ID" | tee -a "$log_file"
    exit 1
  fi

  echo "发送获取用户资料请求..." | tee -a "$log_file"
  echo "用户IDs: $user_ids_input" | tee -a "$log_file"

  # 构建JSON数组
  json_array=$(printf '%s\n' "${user_ids_array[@]}" | jq -R . | jq -s .)

  request=$(jq -n \
    --argjson user_ids "$json_array" \
    '{
      "jsonrpc": "2.0",
      "id": 7,
      "method": "tools/call",
      "params": {
        "name": "slack_get_users_profile",
        "arguments": {
          "user_ids": $user_ids
        }
      }
    }')
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

# 创建日志文件
server_log_file="logs/server_${log_date}.log"
error_log_file="logs/error_${log_date}.log"
echo -e "\n[${current_time}] ====== 服务器日志 =====" >>"$server_log_file"
echo -e "\n[${current_time}] ====== 错误日志 =====" >>"$error_log_file"

# 检查初始 JSON 格式
if ! echo "$request" | jq -c '.' >/dev/null 2>>"$error_log_file"; then
  echo "错误: JSON 格式无效" | tee -a "$error_log_file"
  exit 1
fi

# 运行命令，分别处理标准输出、标准错误和其他错误
(
  echo "$request" |
    jq -c '.' 2>>"$error_log_file" |
    env $(cat local.env | egrep -v '^#' | xargs) go run ./main/main.go 2>>"$server_log_file" |
    tee -a "$log_file" |
    jq '.' 2>>"$error_log_file"
) || {
  echo "错误: 命令执行失败，详细信息请查看错误日志" | tee -a "$error_log_file"
  exit 1
}

echo "" | tee -a "$log_file"
echo "完整日志已保存到: $log_file"
echo "服务器日志已保存到: $server_log_file"
echo "错误日志已保存到: $error_log_file"

# 如果错误日志不为空，提示查看
if [ -s "$error_log_file" ]; then
  echo "警告: 检测到错误，请查看错误日志文件"
fi
