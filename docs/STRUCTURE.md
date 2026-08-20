# Klikku — Project Structure

```
klikku/
├── docs/
│   ├── ARCHITECTURE.md      # System design & decisions
│   ├── PLAN.md              # Roadmap & phases
│   └── STRUCTURE.md         # This file
├── scripts/
│   └── update.sh            # Deploy script (webhook target)
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── handlers/
│   │   │   ├── auth.go
│   │   │   ├── merchant.go
│   │   │   ├── campaign.go
│   │   │   ├── template.go
│   │   │   ├── session.go
│   │   │   ├── device.go
│   │   │   └── upload.go
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── tenant.go
│   │   │   ├── cors.go
│   │   │   └── ratelimit.go
│   │   ├── models/
│   │   │   ├── user.go
│   │   │   ├── merchant.go
│   │   │   ├── campaign.go
│   │   │   ├── template.go
│   │   │   ├── device.go
│   │   │   ├── session.go
│   │   │   ├── photo.go
│   │   │   ├── payment.go
│   │   │   ├── print_job.go
│   │   │   └── email_delivery.go
│   │   ├── repository/
│   │   │   ├── user_repo.go
│   │   │   ├── merchant_repo.go
│   │   │   ├── campaign_repo.go
│   │   │   ├── template_repo.go
│   │   │   ├── session_repo.go
│   │   │   └── photo_repo.go
│   │   ├── services/
│   │   │   ├── auth_service.go
│   │   │   ├── image_service.go
│   │   │   ├── storage_service.go
│   │   │   ├── payment_service.go
│   │   │   ├── email_service.go
│   │   │   └── session_service.go
│   │   ├── utils/
│   │   │   ├── jwt.go
│   │   │   ├── hash.go
│   │   │   ├── s3.go
│   │   │   ├── qr.go
│   │   │   └── response.go
│   │   └── routes/
│   │       └── routes.go
│   ├── migrations/
│   │   ├── 000001_create_merchants.up.sql
│   │   ├── 000001_create_merchants.down.sql
│   │   ├── 000002_create_users.up.sql
│   │   ├── 000002_create_users.down.sql
│   │   └── ...
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── App.tsx
│   │   │   ├── router.tsx
│   │   │   └── providers.tsx
│   │   ├── components/
│   │   │   ├── ui/
│   │   │   ├── layout/
│   │   │   └── photobooth/
│   │   ├── features/
│   │   │   ├── auth/
│   │   │   ├── dashboard/
│   │   │   ├── campaigns/
│   │   │   ├── templates/
│   │   │   ├── sessions/
│   │   │   ├── branding/
│   │   │   └── devices/
│   │   ├── hooks/
│   │   │   ├── useAuth.ts
│   │   │   ├── useApi.ts
│   │   │   └── usePhotobooth.ts
│   │   ├── lib/
│   │   │   ├── api.ts
│   │   │   ├── utils.ts
│   │   │   └── constants.ts
│   │   ├── stores/
│   │   │   ├── authStore.ts
│   │   │   └── merchantStore.ts
│   │   ├── types/
│   │   │   ├── index.ts
│   │   │   └── api.ts
│   │   ├── styles/
│   │   │   └── globals.css
│   │   ├── main.tsx
│   │   └── vite-env.d.ts
│   ├── index.html
│   ├── package.json
│   ├── tailwind.config.js
│   ├── tsconfig.json
│   ├── vite.config.ts
│   └── Dockerfile
├── docker-compose.yml
├── .env.example
├── .gitignore
└── Makefile
```

---

## File Conventions

- Go: PascalCase for exported, camelCase for internal
- TS/TSx: PascalCase for components, camelCase for utilities
- Migrations: `{version}_{name}.up.sql` / `{version}_{name}.down.sql`
- CSS: Tailwind utility-first, CSS variables for theme colors
