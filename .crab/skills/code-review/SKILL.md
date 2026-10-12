---
name: code-review
description: Comprehensive code review and analysis tool for quality, security, and performance
---

# Code Review Skill

When to use this skill: Use when reviewing codebase changes, pull requests, or components for bugs, security vulnerabilities, edge cases, performance issues, and cleanliness.

## Steps:
1. Inspect uncommitted or recent git diffs using git_time_travel or file reading.
2. Check for correctness, security risks (XSS, injection, auth gaps), edge cases, and performance bottlenecks.
3. Report findings categorized by severity: Critical, Warning, Nit, with exact line numbers and proposed fixes.
4. If `--fix` is passed, automatically apply the fixes with `edit_file` and verify with TypeScript/build checks.
