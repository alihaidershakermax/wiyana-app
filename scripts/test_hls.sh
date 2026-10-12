#!/bin/bash
# Local integration test: encrypted HLS playback flow (device binding + AES key).
set -u
BASE=http://localhost:8380

echo "== 1. dev login =="
TOKEN=$(curl -s -X POST "$BASE/api/v1/auth/dev-login" -H 'Content-Type: application/json' -d '{"role":"owner"}' | python3 -c 'import sys,json;print(json.load(sys.stdin)["access_token"])')
echo "login OK (token head: ${TOKEN:0:16}...)"

echo "== 2. video id (pre-created on Convex) =="
VIDEO_ID="j579m4f79rmf1yspt558ct797s8cpmzy"

echo "== 3. create playback session =="
SESS=$(curl -s -X POST "$BASE/api/v1/playback/session" -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d "{\"videoId\":\"$VIDEO_ID\",\"courseId\":\"test-course\"}")
echo "$SESS" | python3 -c 'import sys,json;d=json.load(sys.stdin);print("session ok:",d.get("success"))'
PLAYBACK_TOKEN=$(echo "$SESS" | python3 -c 'import sys,json;print(json.load(sys.stdin)["token"])')

echo "== 4. create grant (with device id) =="
GRANT=$(curl -s -X POST "$BASE/api/stream/hls/$VIDEO_ID" -H "Authorization: Bearer $TOKEN" -H "X-Playback-Token: $PLAYBACK_TOKEN" -H "X-Device-ID: device-abc-123")
echo "$GRANT" | head -c 300; echo
PLAYLIST_URL=$(echo "$GRANT" | python3 -c 'import sys,json;print(json.load(sys.stdin)["playlist_url"])')

echo "== 5. fetch playlist through the route =="
PL=$(curl -s "$BASE$PLAYLIST_URL" -H "Authorization: Bearer $TOKEN" -H "X-Device-ID: device-abc-123")
echo "$PL" | head -16

echo "== 6. fetch first segment WITH device header (expect 200) =="
SEG=$(echo "$PL" | grep -v '^#' | grep '\.ts' | head -1)
echo "segment url: $SEG"
curl -s -o /tmp/seg_ok.ts -w 'with-device: %{http_code} bytes=%{size_download}\n' "$BASE$SEG" -H "Authorization: Bearer $TOKEN" -H "X-Device-ID: device-abc-123"

echo "== 7. fetch same segment WITHOUT device header (expect 403) =="
curl -s -o /dev/null -w 'without-device: %{http_code}\n' "$BASE$SEG" -H "Authorization: Bearer $TOKEN"

echo "== 8. fetch key file WITH device header (expect 200, 16 bytes) =="
KEY_LINE=$(echo "$PL" | grep 'stream.key' | head -1)
KEY_URL=$(echo "$KEY_LINE" | python3 -c "import sys,re;m=re.search(r'URI=\"([^\"]+)\"',sys.stdin.read().strip());print(m.group(1))")
echo "key url: $KEY_URL"
curl -s -o /tmp/key.bin -w 'key: %{http_code} bytes=%{size_download}\n' "$BASE$KEY_URL" -H "Authorization: Bearer $TOKEN" -H "X-Device-ID: device-abc-123"

echo "== 9. fetch key WITHOUT device header (expect 403) =="
curl -s -o /dev/null -w 'key without device: %{http_code}\n' "$BASE$KEY_URL" -H "Authorization: Bearer $TOKEN"

echo "== 10. fetch segment with WRONG device id (expect 403) =="
curl -s -o /dev/null -w 'wrong-device: %{http_code}\n' "$BASE$SEG" -H "Authorization: Bearer $TOKEN" -H "X-Device-ID: device-evil-999"

echo "== 11. fetch segment WITHOUT auth (expect 401) =="
curl -s -o /dev/null -w 'no-auth: %{http_code}\n' "$BASE$SEG"

echo "== 12. verify: segment is encrypted and key decrypts it =="
KEY_HEX=$(od -An -tx1 /tmp/key.bin | tr -d ' \n')
IV_LINE=$(echo "$PL" | grep 'EXT-X-KEY' | sed 's/.*IV=0x//')
openssl enc -d -aes-128-cbc -in /tmp/seg_ok.ts -K "$KEY_HEX" -iv "$IV_LINE" -out /tmp/seg_plain.ts 2>/dev/null && echo "segment decrypted with the route-served key — AES-128 encryption WORKS" || echo "openssl decrypt failed (key/IV mismatch)"
