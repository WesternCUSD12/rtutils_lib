#!/usr/bin/env bash
set -euo pipefail

# Load environment if present
if [ -f .env ]; then
  # shellcheck disable=SC1091
  . .env
fi

if [ -z "${RT_BASE_URL:-}" ] || [ -z "${RT_TOKEN:-}" ]; then
  echo "RT_BASE_URL or RT_TOKEN unset" >&2
  exit 1
fi

BASE="$RT_BASE_URL"
TOKEN="$RT_TOKEN"
PER_PAGE=100
PAGE=1
OUT=/tmp/rt_model_values_raw.txt
: > "$OUT"

while :; do
  URL="$BASE/REST/2.0/assets?query=CustomField.%7BModel%7D+LIKE+%27%25%27&per_page=$PER_PAGE&page=$PAGE"
  echo "Fetching assets page $PAGE" >&2
  RESP=$(curl -sS -H "Authorization: token $TOKEN" "$URL")
  # extract asset URLs
  echo "$RESP" | jq -r '.items[]?._url' | while read -r U; do
    A=$(curl -sS -H "Authorization: token $TOKEN" "$U")
    echo "$A" | jq -r '.CustomFields[]? | select(.name=="Model") | .values[]?' >> "$OUT" || true
  done
  PAGES=$(echo "$RESP" | jq -r '.pages // 0')
  if [ "$PAGES" -eq 0 ] || [ "$PAGE" -ge "$PAGES" ]; then
    break
  fi
  PAGE=$((PAGE+1))
done

SORTED=/tmp/rt_model_values.txt
sort -f "$OUT" | awk '!seen[tolower($0)]++' > "$SORTED"
echo "Saved unique values to $SORTED"
wc -l "$SORTED" || true

# show first 200 lines
echo "--- sample values ---"
head -n 200 "$SORTED" || true
