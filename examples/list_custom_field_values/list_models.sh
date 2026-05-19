#!/usr/bin/env bash
set -euo pipefail

# Usage: bash list_models.sh [CUSTOM_FIELD_ID]
# Defaults to CUSTOM_FIELD_ID=4 (Model)

CF_ID=${1:-4}

command -v curl >/dev/null 2>&1 || { echo "curl is required" >&2; exit 2; }
command -v jq >/dev/null 2>&1 || { echo "jq is required" >&2; exit 2; }

# Ensure env vars are present
if [ -z "${RT_BASE_URL:-}" ] || [ -z "${RT_TOKEN:-}" ]; then
  echo "Please ensure RT_BASE_URL and RT_TOKEN are set in the environment." >&2
  exit 2
fi

# sanitize
BASE=$(printf "%s" "$RT_BASE_URL" | awk '{print $1}')
TOKEN=$(printf "%s" "$RT_TOKEN" | awk -F'=' '{print $NF}' | awk '{print $1}')

echo "Using RT base: $BASE"

TMP=$(mktemp)
trap 'rm -f "$TMP"' EXIT

try_get() {
  local url=$1
  echo "\n-- GET $url --"
  http_code=$(curl -sS -w "%{http_code}" -H "Authorization: token $TOKEN" "$url" -o "$TMP") || true
  echo "HTTP $http_code"
  if [ -s "$TMP" ] && jq -e . "$TMP" >/dev/null 2>&1; then
    echo "JSON response summary:"
    jq -r '(if has("total") then .total elif has("count") then .count elif (.items? and (.items|type=="array")) then (.items|length) elif (type=="array") then (.|length) elif (.items? and (.items|type=="object")) then (.items|length) else 0 end) as $n | "count: "+($n|tostring)' "$TMP"
    echo "sample items (first 5):"
    jq -c 'def sample: if .items then (if (.items|type=="array") then .items[0:5] elif (.items|type=="object") then ([.items|to_entries[]|.value][0:5]) else empty end) elif (type=="array") then .[0:5] else empty end; sample' "$TMP" | sed 's/^/  /'
  else
    echo "Non-JSON or empty response (first 200 chars):"
    head -c 200 "$TMP" | sed 's/^/  /'
  fi
}

try_post_json() {
  local url=$1
  local body=$2
  echo "\n-- POST $url -- body: $body --"
  http_code=$(curl -sS -w "%{http_code}" -H "Authorization: token $TOKEN" -H "Content-Type: application/json" -d "$body" "$url" -o "$TMP") || true
  echo "HTTP $http_code"
  if [ -s "$TMP" ] && jq -e . "$TMP" >/dev/null 2>&1; then
    echo "JSON response summary:"
    jq -r '(if has("total") then .total elif has("count") then .count elif (.items? and (.items|type=="array")) then (.items|length) elif (type=="array") then (.|length) elif (.items? and (.items|type=="object")) then (.items|length) else 0 end) as $n | "count: "+($n|tostring)' "$TMP"
    echo "sample items (first 5):"
    jq -c 'def sample: if .items then (if (.items|type=="array") then .items[0:5] elif (.items|type=="object") then ([.items|to_entries[]|.value][0:5]) else empty end) elif (type=="array") then .[0:5] else empty end; sample' "$TMP" | sed 's/^/  /'
  else
    echo "Non-JSON or empty response (first 200 chars):"
    head -c 200 "$TMP" | sed 's/^/  /'
  fi
}

echo "Trying common endpoints for custom field id $CF_ID"

# 1. GET /customfield/:id/values
try_get "$BASE/REST/2.0/customfield/${CF_ID}/values"

# 2. POST /customfield/:id/values (search) - try empty and simple body
try_post_json "$BASE/REST/2.0/customfield/${CF_ID}/values" '{}'
try_post_json "$BASE/REST/2.0/customfield/${CF_ID}/values" '{"query":{}}'

# 3. GET /customfield/:id
try_get "$BASE/REST/2.0/customfield/${CF_ID}"

# 4. GET with category param (if supported)
try_get "$BASE/REST/2.0/customfield/${CF_ID}?category=Model"

# 5. POST /customfields to search by name
try_post_json "$BASE/REST/2.0/customfields" "{\"query\":\"Name = \\\"Model\\\"\"}"

# 6. Try singular 'value' endpoints
try_get "$BASE/REST/2.0/customfield/${CF_ID}/value"
try_get "$BASE/REST/2.0/customfield/${CF_ID}/value/1"

# 7. Asset-sampling fallback (first page)
echo "\nAsset-sampling fallback: searching assets and showing sample asset URLs"
curl -sS -G -H "Authorization: token $TOKEN" --data-urlencode "query=CustomField.{Model} LIKE '%'" "$BASE/REST/2.0/assets" -o "$TMP"
echo "Assets search returned count:"; jq -r '.total // .count // 0' "$TMP" || true
echo "First 10 asset URLs:"; jq -r '.items[]._url' "$TMP" | sed -n '1,10p'

echo "\nDone. Review the outputs above to see which endpoint returned usable values."

exit 0
