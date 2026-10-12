FROM node:20-slim AS frontend-builder
WORKDIR /build
COPY frontend/package*.json ./
RUN npm ci --no-audit --no-fund --loglevel=error
COPY frontend/ ./
RUN npm run build

FROM golang:1.22-bookworm AS go-builder
WORKDIR /src
COPY backend-go/go.mod backend-go/go.sum ./
RUN go mod download
COPY backend-go/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/wiyana ./cmd/server

FROM debian:bookworm-slim AS runtime
WORKDIR /app
ENV PORT=8000 \
    ENVIRONMENT=production \
    STORE_BACKEND=redis \
    CACHE_ENABLED=true \
    HLS_OUTPUT_DIR=/app/hls_cache \
    UPLOAD_DIR=/app/uploads \
    FFMPEG_BINARY=ffmpeg \
    HOME=/app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl ffmpeg \
    && rm -rf /var/lib/apt/lists/* \
    && addgroup --system --gid 10001 app \
    && adduser --system --uid 10001 --gid 10001 --home /app app \
    && mkdir -p /app/static/spa /app/uploads /app/hls_cache \
    && chown -R app:app /app

COPY --from=go-builder --chown=app:app /out/wiyana /app/wiyana
COPY --from=frontend-builder --chown=app:app /build/dist /app/static/spa

USER app
EXPOSE 8000
HEALTHCHECK --interval=15s --timeout=5s --start-period=30s --retries=3 \
  CMD curl -fsS http://127.0.0.1:${PORT:-8000}/api/health || exit 1

ENTRYPOINT ["/app/wiyana"]
