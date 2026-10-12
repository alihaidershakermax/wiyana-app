# WEYANA PLATFORM — DEPLOYMENT GUIDE (Koyeb & Docker)

Guide for deploying Wiyana Platform services to production on Koyeb or container platforms.

---

## 1. Architecture Overview

- **WEYANA API Service**: Python FastAPI backend container.
- **WEYANA CONTROL PANEL**: Vue 3 single page web application deployed as static files or CDN.
- **WEYANA STUDENT APP**: Vue 3 Telegram Mini App hosted inside Telegram.
- **DATABASE**: MongoDB Atlas or self-hosted MongoDB instance.

---

## 2. Deploying Backend to Koyeb

1. **Create Service on Koyeb**:
   - Source: GitHub Repository (`/backend` or root `Dockerfile`).
   - Builder: Dockerfile.
2. **Configure Environment Variables**:
   - Set `MONGODB_URI`, `DATABASE_NAME`, `SECRET_KEY`, `TELEGRAM_BOT_TOKEN`, `YOUTUBE_CLIENT_ID`, `YOUTUBE_CLIENT_SECRET`.
3. **Port & Health Check**:
   - Port: `8000`.
   - Healthcheck Path: `/api/health`.
4. **Deploy**:
   - Push to main branch or trigger deployment in Koyeb dashboard.

---

## 3. Deploying Control Panel & Student App

1. Build static production assets:
   ```bash
   cd frontend
   npm run build
   ```
2. Host `dist` output folder on Cloudflare Pages, Vercel, Netlify, or Koyeb static service.
3. Configure domain origins in `CORS_ORIGINS` on the backend environment.
