#!/usr/bin/env python3
"""Extract human-readable Arabic/English text runs from Vue files."""
import re, glob, sys

ar = re.compile(r'[\u0600-\u06FF]{5,}')
files = sorted(glob.glob('/home/ubuntu/wena-app-platfrom/frontend/src/**/*.vue', recursive=True))
seen = set()
for f in files:
    src = open(f, encoding='utf-8').read()
    # strip <script> and <style>
    src = re.sub(r'<script[\s\S]*?</script>', '', src)
    src = re.sub(r'<style[\s\S]*?</style>', '', src)
    for m in ar.finditer(src):
        run = m.group(0)
        if run in seen:
            continue
        seen.add(run)
        print(f"{f.split('/frontend/src/')[1]:55s} {run}")
