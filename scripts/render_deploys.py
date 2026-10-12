import json, sys, urllib.request

url = "https://api.render.com/v1/services/srv-d9ss6ae417fc73bb6oh0/deploys?limit=40"
req = urllib.request.Request(url, headers={"Authorization": "Bearer rnd_iJkjH8TwGLymIhvgx6JynMpMzNvu"})
d = json.load(urllib.request.urlopen(req))
if not isinstance(d, list):
    print(json.dumps(d)[:800]); sys.exit(1)
print("keys:", list(d[0].keys()) if d else "empty")
for item in d:
    x = item.get("deploy") if isinstance(item, dict) else item
    if not isinstance(x, dict):
        continue
    c = x.get("commit") or {}
    print(x.get("createdAt", "")[:19], (x.get("status") or "?").ljust(20), str(c.get("id", ""))[:8])
