import json, sys, urllib.request

req = urllib.request.Request(
    "https://api.render.com/v1/services/srv-d9ss6ae417fc73bb6oh0",
    headers={"Authorization": "Bearer rnd_iJkjH8TwGLymIhvgx6JynMpMzNvu"},
)
d = json.load(urllib.request.urlopen(req))
print("keys:", sorted(d.keys()))
ws = d.get("webService")
if ws:
    print("plan:", ws.get("plan"))
    print("runtime:", ws.get("runtime"))
    print("dockerPath:", ws.get("dockerPath"))
    print("healthCheckPath:", ws.get("healthCheckPath"))
    print("buildCommand:", ws.get("buildCommand"))
    print("startCommand:", ws.get("startCommand"))
    print("envCount:", len(ws.get("envVars") or []))
else:
    print(json.dumps(d, indent=1)[:1200])
