# Klikku — Project Plan

**Status:** Phase 1 (MVP)
**Created:** 2026-08-20

---

## Phase 1: Core MVP (Week 1-2)

### Step 1: Project Setup
- [ ] Go backend scaffold (Fiber, pgx, golang-migrate)
- [ ] React frontend scaffold (Vite 5, TS, Tailwind, Zustand)
- [ ] Docker Compose (PostgreSQL, MinIO, Backend, Frontend)
- [ ] Nginx config for klikku.arjism.com
- [ ] update.sh deploy script

### Step 2: Auth & Tenant Foundation
- [ ] User model + migrations
- [ ] Merchant model + migrations
- [ ] JWT auth (register, login, refresh, logout)
- [ ] Tenant middleware
- [ ] Role-based access control

### Step 3: Merchant Branding & Settings
- [ ] Branding API (CRUD logo, colors, name, welcome message)
- [ ] File upload to MinIO
- [ ] Branding applies to photobooth experience

### Step 4: Campaign Management
[ ] Campaign CRUD (name, dates, status, promotion config)
[ ] Active campaign selection logic

### Step 5: Template Management
[ ] Template CRUD (name, photo_count, layout_config JSON, overlay, price)
[ ] Drag-and-drop template editor (basic)
[ ] Preview generation

### Step 6: Photobooth Customer Flow
[ ] Idle/attract screen
[ ] Template selection screen
[ ] Camera capture (getUserMedia)
[ ] Countdown + flash animation
[ ] Photo composition (libvips/ImageMagick)
[ ] Photo review screen
[ ] QR code download

### Step 7: Session & Gallery
[ ] Session state machine
[ ] Session recovery
[ ] Merchant gallery view (sessions list, photos, actions)
[ ] Download/resend/delete

### Step 8: Merchant Dashboard
[ ] Overview (today's sessions, revenue, photos, prints, emails)
[ ] Gallery with lightbox
[ ] Campaign management UI
[ ] Template management UI
[ ] Branding settings UI
[ ] Devices management UI

### Step 9: Super Admin Dashboard
[ ] List all merchants
[ ] Create/edit/suspend/delete merchants
[ ] Platform-wide analytics
[ ] All sessions search

### Step 10: Testing & Polish
[ ] Responsive/tablet testing
[ ] Idle timeout + auto-return
[ ] Loading states + error handling
[ ] Animations (countdown, flash, transitions)
[ ] Print-ready image resolution (1200x1800 min)

---

## Phase 2: Monetization (Week 3-4)

- [ ] Payment provider abstraction (PaymentProvider interface)
- [ ] Xendit integration
- [ ] Midtrans integration
- [ ] Doku integration
- [ ] Webhook handling + signature verification
- [ ] Brevo email integration (transactional)
- [ ] Email template (branded, with download link)
- [ ] Local Print Service (Go binary)
[ ] Print queue management
[ ] Reprint functionality

---

## Phase 3: Advanced (Week 5+)

- [ ] Visual drag-and-drop template editor (advanced)
- [ ] Multiple devices per merchant
- [ ] Campaign scheduling
- [ ] Advanced analytics
- [ ] Social media promotions
- [ ] Discount vouchers
- [ ] AI backgrounds / effects
- [ ] Subscription management
- [ ] Merchant self-registration

---

## Priority Order

1. Auth + tenants (foundation)
2. Customer photobooth flow (core value)
3. Template + campaign management (merchant value)
4. Session gallery (merchant value)
5. Dashboards (merchant value)
6. Payment + email + print (monetization)
