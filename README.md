# Klikku — Photo Booth Platform

A modern photo booth platform for events with campaign management, photo processing, and print capabilities.

## Quick Start

```bash
# Clone
git clone https://github.com/lovelymondayz/klikku.git
cd klikku

# Start all services
docker compose up -d --build

# Frontend: http://localhost:3009
# Backend API: http://localhost:8083
# DB: localhost:5434 (user: klikku)
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        NGINX (80/443)                        │
│                     klikku.arjism.com → :3009                │
├─────────────────────────────────────────────────────────────┤
│  React + Vite + TS + Tailwind  │  Go + GIN + pgx + Postgres │
│        (Frontend :3009)        │       (Backend :8083)      │
├─────────────────────────────────────────────────────────────┤
│              PostgreSQL :5434  │  Local Storage (/data)     │
└─────────────────────────────────────────────────────────────┘
```

## Features

- **Photo Booth**: Capture and process photos at events
- **Campaign Management**: Create and manage photo campaigns
- **Print Integration**: Direct printing to photo printers
- **Download Tokens**: Secure download links with expiry
- **Email Delivery**: Send photos via email
- **Responsive UI**: Mobile-first design with Tailwind CSS

## API Endpoints

### Public
- `GET /api/health` — Health check
- `GET /api/campaigns` — List campaigns
- `GET /api/campaigns/:id` — Get campaign

### Authenticated (JWT required)
- `POST /api/campaigns` — Create campaign
- `POST /api/photos` — Upload photo
- `GET /api/photos/:id` — Get photo
- `POST /api/print` — Print photo
- `POST /api/email` — Email photo

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| DB_PASSWORD | - | Database password |
| DB_USER | klikku | Database user |
| DB_NAME | klikku | Database name |
| JWT_SECRET | - | JWT signing key |

## Development

```bash
# Backend only
cd backend
go run .

# Frontend only
cd frontend
npm install
npm run dev
```

## Deployment

1. Push to `main` → GitHub Action auto-deploys
2. Or manually: `ssh vps && cd /root/klikku && ./update.sh`

## License

MIT
