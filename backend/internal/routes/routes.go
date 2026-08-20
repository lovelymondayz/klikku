package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"klikku/internal/config"
	"klikku/internal/handlers"
	"klikku/internal/middleware"
	"klikku/internal/utils"
)

func Setup(app *fiber.App, db *pgxpool.Pool, storage *utils.Storage, cfg *config.Config) {
	api := app.Group("/api")

	// Public routes
	api.Post("/auth/register", handlers.Register(db, cfg))
	api.Post("/auth/login", handlers.Login(db, cfg))
	api.Post("/auth/refresh", handlers.RefreshToken(db, cfg))
	api.Post("/auth/logout", handlers.Logout())

	// Public photobooth endpoints (device-facing)
	api.Get("/devices/:token/attract", handlers.GetAttractScreen(db))
	api.Post("/devices/:token/session", handlers.CreateSession(db))
	api.Get("/sessions/:id", handlers.GetSession(db))
	api.Post("/sessions/:id/capture", handlers.CapturePhoto(db, storage, cfg))
	api.Get("/sessions/:id/download", handlers.DownloadSession(db))
	api.Post("/sessions/:id/finalize", handlers.FinalizeSession(db, storage, cfg))
	api.Get("/download/:id", handlers.DownloadFinal(db, storage))
	api.Get("/download/:id/secure", handlers.ValidateDownloadToken(db, storage))
	api.Post("/sessions/:id/generate-link", handlers.GenerateSecureDownloadURL(db))
	api.Post("/sessions/:id/auto-print", handlers.AutoPrintJob(db, storage))
	api.Post("/sessions/:id/send-email", handlers.SendEmailDelivery(db, storage, cfg))
	api.Post("/sessions/:id/resend-email", handlers.ResendEmail(db, storage, cfg))

	// Protected routes
	auth := api.Group("/", middleware.AuthMiddleware)
	auth.Get("/auth/me", handlers.GetCurrentUser(db))

	// Merchant-scoped routes
	merchant := auth.Group("/", middleware.TenantMiddleware)

	// Campaigns
	merchant.Get("/campaigns", handlers.ListCampaigns(db))
	merchant.Post("/campaigns", handlers.CreateCampaign(db))
	merchant.Get("/campaigns/:id", handlers.GetCampaign(db))
	merchant.Put("/campaigns/:id", handlers.UpdateCampaign(db))
	merchant.Delete("/campaigns/:id", handlers.DeleteCampaign(db))

	// Templates
	merchant.Get("/templates", handlers.ListTemplates(db))
	merchant.Post("/templates", handlers.CreateTemplate(db))
	merchant.Get("/templates/:id", handlers.GetTemplate(db))
	merchant.Put("/templates/:id", handlers.UpdateTemplate(db))
	merchant.Delete("/templates/:id", handlers.DeleteTemplate(db))

	// Sessions & Gallery
	merchant.Get("/sessions", handlers.ListSessions(db))
	merchant.Get("/sessions/:id", handlers.GetSessionDetail(db))
	merchant.Delete("/sessions/:id", handlers.DeleteSession(db))
	merchant.Post("/sessions/:id/resend-email", handlers.ResendEmail(db, storage, cfg))

	// Print Jobs
	merchant.Get("/print-jobs", handlers.ListPrintJobs(db))
	merchant.Get("/print-jobs/:id", handlers.GetPrintJob(db))
	merchant.Post("/sessions/:id/print", handlers.CreatePrintJob(db, storage))
	merchant.Put("/print-jobs/:id/status", handlers.UpdatePrintJobStatus(db))
	merchant.Post("/sessions/:id/reprint", handlers.Reprint(db, storage))
	merchant.Get("/print-jobs/pending", handlers.GetPendingPrintJobs(db))

	// Branding
	merchant.Get("/branding", handlers.GetBranding(db))
	merchant.Put("/branding", handlers.UpdateBranding(db, storage, cfg))

	// Devices
	merchant.Get("/devices", handlers.ListDevices(db))
	merchant.Post("/devices", handlers.CreateDevice(db))
	merchant.Get("/devices/:id", handlers.GetDevice(db))
	merchant.Put("/devices/:id", handlers.UpdateDevice(db))
	merchant.Delete("/devices/:id", handlers.DeleteDevice(db))

	// Analytics
	merchant.Get("/analytics/overview", handlers.GetAnalyticsOverview(db))

	// Super Admin routes
	superAdmin := auth.Group("/admin", middleware.SuperAdminMiddleware)
	superAdmin.Get("/merchants", handlers.AdminListMerchants(db))
	superAdmin.Post("/merchants", handlers.AdminCreateMerchant(db, cfg))
	superAdmin.Put("/merchants/:id", handlers.AdminUpdateMerchant(db))
	superAdmin.Delete("/merchants/:id", handlers.AdminDeleteMerchant(db))
	superAdmin.Get("/sessions", handlers.AdminListAllSessions(db))
	superAdmin.Get("/analytics", handlers.AdminGetPlatformAnalytics(db))

	// Upload (merchant-scoped)
	merchant.Post("/upload", handlers.UploadFile(db, storage, cfg))
}
