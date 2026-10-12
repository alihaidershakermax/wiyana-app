# Unified API Contract

## Purpose

Wiyana uses a versioned JSON API for application and administration traffic. The frontend must consume the API through `frontend/src/api/client.ts` and the stable `frontend/src/api/index.ts` facade. Views and components should not import Axios directly or compose endpoint paths for domain operations.

## Base paths

| Scope | Canonical path | Compatibility paths |
|---|---|---|
| Application data | `/api/v1` | Existing root paths remain available only where media or legacy clients require them |
| Administration | `/api/v1` with admin bearer token | Existing panel route aliases remain temporarily |
| HLS media | `/api/stream/hls/:video_id/*` | Required for segment delivery and cache headers |
| Direct media | `/api/stream/direct/:video_id` | Kept behind access checks and short-lived playback authorization |
| Protection | `/api/v1/playback/*`, `/api/v1/security/*` | `/api/check-access` and `/api/report-violation` remain compatibility aliases |

## Response envelopes

Successful domain responses should use the following shape for new endpoints:

```json
{
  "success": true,
  "data": {},
  "meta": {},
  "request_id": "request-id"
}
```

Error responses should use this shape while retaining `detail` during migration:

```json
{
  "success": false,
  "error": {
    "code": "ACCESS_DENIED",
    "message": "Access denied",
    "details": {}
  },
  "detail": "Access denied",
  "request_id": "request-id"
}
```

The frontend normalizes both old and new responses through `ApiError`, exposing `status`, `code`, `details`, and `requestId`. This lets the backend evolve without requiring every screen to understand transport-specific Axios errors.

## Authentication

User requests use `wiyana_access_token` and Telegram Mini App init data when available. Administration requests use separate `wiyana_admin_access_token` and `wiyana_admin_refresh_token` storage. Both flows share the same transport factory and refresh queue, but their token stores and session-expired events remain isolated.

## Cache and request behavior

GET requests under `/api/v1` are cacheable only when they are not authentication, playback, or security endpoints. User-scoped cache keys include a digest of bearer and Telegram headers. Mutations invalidate HTTP and catalog cache prefixes. Clients send `X-Request-ID` so logs can correlate frontend failures with backend requests.

## Migration rule

Existing modules under `frontend/src/api/*.ts` remain compatibility domain adapters. New screens should import from `@/api` and call `api.courses`, `api.lessons`, `api.videos`, `api.profile`, or another named domain. A domain adapter may be moved internally later without changing screen imports.

## Go service boundaries

The Go server bootstrap is in `cmd/server/main.go`. HTTP concerns are split into `http_routes.go`, `http_middleware.go`, `http_response.go`, `http_auth.go`, `http_playback.go`, and `http_compat.go`. Business logic remains under `internal/auth`, `internal/cache`, `internal/security`, `internal/store`, and `internal/video`.
