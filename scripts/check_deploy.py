#!/usr/bin/env python3
"""Check the latest Render deploy status for wena-app."""
import json
import os
import urllib.request

KEY = os.environ["RENDER_API_KEY"]
URL = "https://api.render.com/v1/services/srv-d9ss6ae417fc73bb6oh0/deploys?limit=3"

req = urllib.request.Request(URL, headers={"Authorization": f"Bearer {KEY}"})
data = json.load(urllib.request.urlopen(req, timeout=30))
for x in data:
    commit = (x.get("commit") or {}).get("sha", "?")[:7]
    print(x["id"], x["status"], commit, x.get("startedAt", ""), x.get("finishedAt", ""))
