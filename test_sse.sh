#!/bin/bash

# 检查是否安装了curl
if ! command -v curl &>/dev/null; then
  echo "错误: 请先安装 curl 工具"
  echo "可以使用以下命令安装:"
  echo "  Mac: brew install curl"
  echo "  Ubuntu/Debian: sudo apt-get install curl"
  exit 1
fi

# 检查参数
if [ "$#" -lt 1 ]; then
  echo "用法: $0 <请求类型> [参数]"
  echo "请求类型:"
  echo "  health           - 健康检查"
  echo "  init              - 初始化请求"
  echo "  list              - 获取工具列表"
  echo "  list_channels     - 列出频道 [limit] [cursor]"
  echo "  get_thread_replies - 获取消息线程回复 <thread_url>"
  echo ""
  echo "示例:"
  echo "  $0 health"
  echo "  $0 init"
  echo "  $0 list"
  echo "  $0 list_channels 50"
  echo "  $0 get_thread_replies \"https://your-workspace.slack.com/archives/C0734812MFG/p1742788004223029\""
  exit 1
fi

# 创建日志目录
mkdir -p logs

# 使用日期作为日志文件名
log_date=$(date +"%Y%m%d")
log_file="logs/sse_requests_${log_date}.log"

# 添加时间戳到日志条目
current_time=$(date +"%H:%M:%S")
echo -e "\n[${current_time}] ====== 新SSE请求 =====" | tee -a "$log_file"

request_type=$1
shift

# 基础URL
base_url="http://localhost:3333"

# SSE请求头
sse_headers=(-H "Accept: text/event-stream")

# 检查服务器是否在运行
check_server() {
  if ! curl -s "${base_url}/health" >/dev/null; then
    echo "错误: 无法连接到服务器，请确保服务器正在运行 (make run-local)" | tee -a "$log_file"
    exit 1
  fi
}

case "$request_type" in
health)
  echo "检查服务器健康状态..." | tee -a "$log_file"
  curl -s "${base_url}/health" | tee -a "$log_file"
  echo "" # 添加换行
  ;;

init)
  check_server
  echo "发送初始化请求..." | tee -a "$log_file"
  curl -N "${sse_headers[@]}" "${base_url}/mcp/initialize" | tee -a "$log_file"
  ;;

list)
  check_server
  echo "发送工具列表请求..." | tee -a "$log_file"
  curl -N "${sse_headers[@]}" "${base_url}/mcp/tools/list" | tee -a "$log_file"
  ;;

list_channels)
  check_server
  limit=${1:-100}
  cursor=$2
  echo "发送列出频道请求..." | tee -a "$log_file"
  echo "限制: $limit" | tee -a "$log_file"
  if [ -n "$cursor" ]; then
    echo "游标: $cursor" | tee -a "$log_file"
    curl -N "${sse_headers[@]}" "${base_url}/mcp/slack/list-channels?limit=${limit}&cursor=${cursor}" | tee -a "$log_file"
  else
    curl -N "${sse_headers[@]}" "${base_url}/mcp/slack/list-channels?limit=${limit}" | tee -a "$log_file"
  fi
  ;;

get_thread_replies)
  check_server
  thread_url="$1"
  if [ -z "$thread_url" ]; then
    echo "错误: 请提供thread_url参数" | tee -a "$log_file"
    exit 1
  fi

  echo "发送获取线程回复请求..." | tee -a "$log_file"
  echo "线程URL: $thread_url" | tee -a "$log_file"
  curl -N "${sse_headers[@]}" "${base_url}/mcp/slack/get-thread-replies?thread_url=${thread_url}" | tee -a "$log_file"
  ;;

*)
  echo "错误: 未知的请求类型 '$request_type'" | tee -a "$log_file"
  exit 1
  ;;
esac

echo -e "\n完整日志已保存到: $log_file"
