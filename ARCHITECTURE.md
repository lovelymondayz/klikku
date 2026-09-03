# Klikku — Architecture

## System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Cloudflare Edge                          │
│                     klikku.arjism.com (HTTPS)                    │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Cloudflare Tunnel (cf-tunnel)                │
│              http://192.168.88.101:8083 (plain HTTP)            │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Nginx Reverse Proxy                      │
│                    :8083 → :8083 (backend)                      │
│                    :3009 → :3009 (frontend)                     │
└────────────────────────────┬────────────────────────────────────┘
                             │
              ┌──────────────┴──────────────┐
              ▼                              ▼
┌──────────────────────┐        ┌──────────────────────┐
│   Go + GIN Backend   │        │  React + Vite + TS   │
│   :8083 (internal)   │        │  :3009 (internal)    │
│                      │        │                      │
│  - JWT Auth          │        │  - Tailwind CSS      │
│  - pgx + Postgres    │        │  - react-router-dom  │
│  - Campaign CRUD     │        │  - Photo Capture     │
│  - Photo Processing  │        │  - Campaign Manager  │
│  - Print Queue       │        │  - Print Interface   │
│  - Email Delivery    │        │  - Download Manager  │
│  - Download Tokens   │        │                      │
└──────────┬───────────┘        └──────────────────────┘
           │
           ▼
┌──────────────────────┐
│   PostgreSQL :5434   │
│                      │
│  - Merchants         │
│  - Campaigns         │
│  - Sessions          │
│  - Photos            │
│  - Print Jobs        │
│  - Payments          │
│  - Email Deliveries  │
│  - Download Tokens   │
└──────────────────────┘
```

## Tech Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Language | Go | 1.22+ |
| Web Framework | Gin | v1.10 |
| Database Driver | pgx | v4 |
| Auth | JWT (golang-jwt/jwt) | v5 |
| Migration | Goose | - |
| Frontend | React + Vite + TypeScript | Vite 5, React 18 |
| Styling | Tailwind CSS | v3 |
| Routing | react-router-dom | v6 |
| Deployment | Docker Compose | v3.8 |
| Reverse Proxy | Nginx | - |
| Tunnel | Cloudflare Tunnel | - |

## Key Design Decisions

### 1. Campaign-Based Architecture
- Each event/merchant has a campaign
- Photos belong to campaigns
- Campaigns have configurable settings

### 2. Download Token System
- SHA-256 tokens with expiry
- Tokens linked to photobooth sessions
- Automatic cleanup of expired tokens

### 3. Photo Processing Pipeline
- Upload → Validate → Resize → Compress → Store
- Multiple size variants (original, thumbnail, print)
- EXIF data preservation

### 4. Print Queue
- Async print job processing
- Printer selection per campaign
- Error handling and retry logic

### 5. Email Delivery
- Async email sending
- Template-based emails
- Delivery tracking

## API Endpoints

### Public
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/health` | Health check |
| GET | `/api/campaigns/:id` | Get campaign |
| GET | `/api/photos/:id` | Get photo |
| GET | `/api/download/:token` | Download photo |

### Authenticated (JWT required)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/campaigns` | Create campaign |
| GET | `/api/campaigns` | List campaigns |
| POST | `/api/sessions` | Create session |
| POST | `/api/photos` | Upload photo |
| POST | `/api/print` | Print photo |
| POST | `/api/email` | Email photo |

## Ports

| Service | External | Internal |
|---------|----------|----------|
| Backend | `:8083` | `:8083` |
| Frontend | `:3009` | `:80` |
| DB | `:5434` | `:5432` |
