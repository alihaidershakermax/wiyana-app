# WEYANA PLATFORM — ENVIRONMENT SPECIFICATION

| Variable Name | Required | Default | Description |
|---|---|---|---|
| `MONGODB_URI` | Yes | - | MongoDB connection string |
| `DATABASE_NAME` | Yes | `wiyana_db` | Target database name |
| `SECRET_KEY` | Yes | - | 64+ char secret for JWT authentication tokens |
| `TELEGRAM_BOT_TOKEN` | Yes | - | Bot token from @BotFather |
| `TELEGRAM_WEBHOOK_SECRET` | Yes | - | Secret string for verifying webhook calls |
| `YOUTUBE_CLIENT_ID` | Optional | - | Google Cloud OAuth Client ID for YouTube Data API |
| `YOUTUBE_CLIENT_SECRET` | Optional | - | Google Cloud OAuth Client Secret |
| `CORS_ORIGINS` | Yes | `*` | Allowed origins for cross-origin requests |
