# Wiyana Platform API and UI/UX Audit

## Design read

This is a redesign of an existing Arabic Telegram Mini App for learners, teachers, and administrators. The product should feel premium, calm, and trustworthy, with an iOS-inspired Glass Liquid language applied selectively rather than as a generic glass effect. The implementation should preserve the orange Wiyana accent, RTL layout, mobile-first behavior, and clear content hierarchy.

The selected reference skill is `design-taste-frontend` from https://github.com/Leonxlnx/taste-skill. Its main constraints adopted here are: audit before redesign, one coherent visual system, reduced-motion support, dark/light parity, WCAG-conscious contrast, motion limited to transform/opacity, and no decorative UI that does not carry meaning.

## Current architecture findings

The frontend already has Axios, but transport concerns are duplicated between `frontend/src/api/client.ts` and `frontend/src/admin/api/adminApi.ts`. Both clients attach tokens and implement response handling, while the admin client separately implements refresh queuing. This makes auth behavior inconsistent and makes future API changes expensive.

The frontend API surface is split between `src/api/*`, `src/services/*`, direct `apiClient` calls inside views, and a second admin client. `src/api/protectionAdmin.ts` probes multiple endpoint shapes, which is a compatibility workaround rather than a stable contract. The new architecture should have one transport factory, one normalized error type, one token-store abstraction, and domain modules under `src/api/domains/`.

The Go server is functionally modular under `internal/`, but `backend-go/cmd/server/main.go` is approximately 1,032 lines. It combines bootstrap, middleware, caching, CORS, authentication, playback, HLS delivery, security administration, compatibility routes, response helpers, and empty fallback payloads. The first safe split should preserve `package main` and routes while moving shared HTTP helpers, middleware, playback/HLS, and compatibility response builders into focused files. A later step can move route handlers behind an explicit `/api/v1` contract without breaking the existing frontend.

## Canonical API contract direction

The canonical JSON API remains `/api/v1`. Root-level media and protection endpoints remain explicit exceptions because HLS assets and existing playback clients depend on `/api/stream/*`, `/api/check-access`, and `/api/report-violation`. The frontend should consume these exceptions through named domain methods rather than constructing paths in view components.

Responses should retain compatibility fields during migration, but new client code should prefer typed envelopes: `{ success, data, meta, request_id }` for ordinary API responses and `{ success: false, error: { code, message, details }, request_id }` for failures. The Go server can add `X-Request-ID` and the normalized error shape without removing `detail` and `error` until all callers migrate.

## Proposed module boundaries

| Area | Target location | Responsibility |
|---|---|---|
| HTTP transport | `frontend/src/api/client.ts` | Base URL, auth headers, Telegram init data, refresh, normalized errors, request IDs |
| Token storage | `frontend/src/api/authStore.ts` | User/admin access and refresh tokens with SSR-safe guards |
| Domain contracts | `frontend/src/api/domains/*.ts` | Typed courses, home, lessons, videos, security, notifications, admin methods |
| Compatibility facade | `frontend/src/api/index.ts` | Stable exports for existing callers during migration |
| Go response helpers | `backend-go/cmd/server/http_response.go` | JSON, errors, request IDs, consistent headers |
| Go middleware | `backend-go/cmd/server/http_middleware.go` | CORS, rate limiting, cache, security headers |
| Go media routes | `backend-go/cmd/server/http_playback.go` | Playback sessions, grants, HLS, direct streams |
| Go compatibility | `backend-go/cmd/server/http_compat.go` | Temporary legacy route responses and migration comments |
| UI tokens | `frontend/src/styles.css` | Semantic color, surface, spacing, focus, reduced motion, glass fallback |
| UI primitives | `frontend/src/components/ui/` | Shared page header, state panels, buttons, surfaces, skeletons |

## UI/UX priorities

The redesign should first improve the application shell, page headers, loading/empty/error states, touch targets, focus visibility, and mobile bottom navigation. Glass Liquid should be reserved for navigation, elevated surfaces, and focused controls; content cards should remain readable and not become translucent by default. All non-essential motion must collapse under `prefers-reduced-motion`, and `prefers-reduced-transparency` must receive a solid fallback.
