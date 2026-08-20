# Klikku — Architecture Document

**Status:** Active
**Created:** 2026-08-20
**Stack:** React + Vite + TS + Tailwind | Go + Fiber | PostgreSQL | MinIO

---

## 1. System Overview

Klikku is a multi-merchant white-label photobooth platform. One central platform owner (Super Admin) onboards multiple merchants. Each merchant gets isolated branding, campaigns, templates, pricing, and dashboards.

The customer-facing photobooth runs on a tablet browser (kiosk-style), guiding users through: Touch → Choose Template → Pay → Pose → Capture → Print → Email/QR → Done.

---

## 2. High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Cloudflare CDN                            │
│                    (klikku.arjism.com)                           │
└──────────────────────────┬──────────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────────┐
│                         Nginx Reverse Proxy                      │
│    /api/* → backend:8083    /    → frontend:3009                 │
└──────────────────────────┬──────────────────────────────────────┘
                           │
        ┌──────────────────┴──────────────────┐
        │                                      │
┌───────▼───────┐                    ┌─────────▼─────────┐
│   Frontend     │                    │    Backend (Go)    │
│  React + Vite  │◄──────────────────►│   Fiber + pgx      │
│  (port 3009)   │                    │   (port 8083)      │
└────────────────┘                    └─────────┬─────────┘
                                                 │
                    ┌────────────────────────────┼──────────────────┐
                    │                            │                  │
             ┌──────▼──────┐            ┌────────▼───────┐  ┌───────▼──────┐
             │ PostgreSQL  │            │     MinIO       │  │  Brevo API   │
             │ (port 5440) │            │ (port 8089)     │  │  (email)     │
             └─────────────┘            └────────────────┘  └──────────────┘
```

---

## 3. Multi-Tenant Isolation

All tenant-scoped queries include `merchant_id` filtering at the repository layer. No cross-merchant data access.

**Roles:**
- `SUPER_ADMIN` — platform owner, sees all merchants
- `MERCHANT_ADMIN` — merchant owner, manages own data only
- `MERCHANT_STAFF` — staff, limited permissions within merchant

**Isolation strategy:**
- JWT token contains `user_id`, `role`, and `merchant_id`
- Every API route uses `TenantMiddleware` that extracts merchant_id from JWT
- Repository layer auto-applies `WHERE merchant_id = $1` to all queries
- Super Admin requests bypass tenant filter (explicitly)

---

## 4. Backend Module Structure (Go + Fiber)

```
klikku/backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/        # Env config loader
│   ├── handlers/      # HTTP handlers (Fiber)
│   ├── middleware/    # Auth, tenant, CORS, logger
│   ├── models/        # GORM/pgx models
│   ├── repository/    # DB queries (per entity)
│   ├── services/      # Business logic
│   ├── utils/         # Helpers (hash, JWT, S3)
│   └── routes/        # Route registration
├── migrations/        # SQL migrations (golang-migrate)
├── go.mod
└── Dockerfile
```

---

## 5. Frontend Module Structure (React + Vite + TS)

```
klikku/frontend/src/
├── app/               # App shell, providers, router
├── components/        # Reusable UI
│   ├── ui/            # Buttons, inputs, modals
│   ├── layout/        # Sidebar, header, container
│   └── photobooth/    # Countdown, flash, template preview
├── features/
│   ├── auth/          # Login, register
│   ├── dashboard/     # Merchant overview, analytics
│   ├── campaigns/     # CRUD
│   ├── templates/     # CRUD, visual editor
│   ├── sessions/      # Gallery, detail, actions
│   ├── branding/      # Merchant settings
│   └── devices/       # Device management
├── hooks/             # Custom React hooks
├── lib/               # API client, utils
├── stores/            # Zustand stores
└── types/             # TypeScript types
```

---

## 6. Session State Machine

```
IDLE → STARTED → TEMPLATE_SELECTED → PAYMENT_PENDING → PAYMENT_CONFIRMED
     → CAMERA_READY → CAPTURING → PROCESSING → PHOTO_READY → PRINTING
     → PRINT_COMPLETE → EMAIL_OPTIONAL → COMPLETED → RETURN_TO_IDLE
```

**Recovery:** Session state persists in DB. On reconnect, frontend polls `GET /api/sessions/:id/state` to resume from current state. Payment confirmation is server-side only (never trust frontend).

---

## 7. Storage Strategy

| Type | Bucket | Access |
|------|--------|--------|
| Original photos | `klikku-originals` | Private, signed URL |
| Processed photos | `klikku-processed` | Private, signed URL |
| Final compositions | `klikku-finals` | Private, signed URL |
| Template assets | `klikku-templates` | Private |
| Merchant logos | `klikku-assets` | Public read |

All uploads go through backend (never direct-to-S3 from frontend for originals). MinIO lifecycle policies auto-expire old originals per merchant retention config.

---

## 8. Payment Architecture (Phase 2)

```
Frontend selects template
→ Backend creates session + pending payment
→ Backend calls Xendit/Midtrans/Doku API → returns payment URL/QR
→ Frontend redirects/shows QR
→ Provider sends webhook → Backend verifies signature
→ Backend marks session PAID → unlocks camera
```

**Pluggable design:** `PaymentProvider` interface with implementations per provider. Provider chosen per-merchant in DB.

```go
type PaymentProvider interface {
    CreatePayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error)
    VerifyWebhook(payload []byte, signature string) (*WebhookResult, error)
    HandleCallback(payload []byte) (*CallbackResult, error)
}
```

---

## 9. Print Architecture (Phase 2)

```
Session READY
→ Backend creates PrintJob (PENDING)
→ Local Print Service polls GET /api/devices/:token/print-jobs
→ Service downloads final image from signed URL
→ Service sends to CUPS/local printer
→ Service PATCH /api/print-jobs/:id (status=PRINTING/DONE/FAILED)
```

The Local Print Service is a separate lightweight Go binary running on the kiosk device. It authenticates via device token.

---

## 10. Security

- **Auth:** bcrypt password hashing + JWT (access token 15min, refresh token 7d, httpOnly cookie)
- **Tenant isolation:** middleware-enforced on every route
- **S3:** all buckets private, presigned URLs for access (expiring)
- **Payment:** webhook signature verification (HMAC), idempotency keys
- **Photos:** unguessable UUID v7 session IDs, no sequential IDs exposed
- **Rate limiting:** Fiber rate limiter on auth + API routes
- **CORS:** strict origin allowlist

---

## 11. MVP Scope (Phase 1)

- ✅ Multi-merchant (CRUD)
- ✅ Merchant branding (logo, colors, name)
- ✅ Customer photobooth flow (touch → template → capture → compose → QR download)
- ✅ Camera capture via getUserMedia
- ✅ Template selection with preview
- ✅ Photo composition (ImageMagick/libvips)
- ✅ QR code download (unguessable URL)
- ✅ Merchant dashboard (sessions, gallery, analytics)
- ✅ Campaign management
- ✅ Template management (JSON config, drag-and-drop editor)
- ✅ Device management
- ✅ Super Admin dashboard
- ❌ Payment integration (Phase 2)
- ❌ Print service (Phase 2)
- ❌ Email delivery (Phase 2)

---

## 12. Ports & Deployment

| Service | Port |
|---------|------|
| PostgreSQL | 5440 |
| Backend | 8083 |
| Frontend | 3009 |
| MinIO API | 8089 |
| MinIO Console | 9091 |

Deploy via `docker compose build --no-cache && docker compose up -d --force-recreate` + `update.sh` webhook.
