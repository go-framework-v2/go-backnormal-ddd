#!/usr/bin/env bash
# 登录接口限流测试：连续请求 35 次，前 30 次通过，第 31 次起应返回 code 429
# 用法：./scripts/test_rate_limit.sh [BASE_URL] [PROJECT_ID]
# 示例：./scripts/test_rate_limit.sh
#       ./scripts/test_rate_limit.sh http://127.0.0.1:8091 1

BASE="${1:-http://127.0.0.1:8080}"
PROJECT_ID="${2:-1}"
TOTAL=35

echo "限流测试: $TOTAL 次请求 -> $BASE/user/login/guest (projectId=$PROJECT_ID)"
echo "预期: 前 30 次 code=200，第 31 次起 code=429"
echo "---"

for i in $(seq 1 "$TOTAL"); do
  body=$(curl -s -X POST "$BASE/user/login/guest" \
    -H "Content-Type: application/json" \
    -H "projectId: $PROJECT_ID" \
    -d '{"data":{"deviceId":"rate-test-dev","oaid":"","model":"test","realChannel":"test"}}')
  # 不依赖 jq：从 JSON 里取出 "code": 数字
  code=$(echo "$body" | grep -o '"code"[[:space:]]*:[[:space:]]*[0-9]*' | head -1 | grep -o '[0-9]*$')
  if [ "$code" = "429" ]; then
    echo "第 $i 次: 被限流 (429)"
  else
    echo "第 $i 次: 通过 (code=${code:-?})"
  fi
done

echo "---"
echo "完成。若第 31 次起为 429 则限流正常。"
