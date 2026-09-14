package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"klikku/internal/config"
	"klikku/internal/handlers"
	"klikku/internal/middleware"
	"klikku/internal/storage"
)

func Setup(app *gin.Engine, db *pgxpool.Pool, store storage.Storage, cfg *config.Config) {
	api := app.Group("/api")

	// Public routes
	api.POST("/auth/register", handlers.Register(db, cfg))
	api.POST("/auth/login", handlers.Login(db, cfg))
	api.POST("/auth/refresh", handlers.RefreshToken(db, cfg))
	api.POST("/auth/logout", handlers.Logout())

	// Public photobooth endpoints (device-facing)
	api.GET("/device/:token/attract", handlers.GetAttractScreen(db))
	api.POST("/device/:token/session", handlers.CreateSession(db))
	api.GET("/sessions/:id", handlers.GetSession(db))
	api.POST("/sessions/:id/capture", handlers.CapturePhoto(db, store))
	api.GET("/sessions/:id/download", handlers.DownloadSession(db))
	api.POST("/sessions/:id/finalize", handlers.FinalizeSession(db, store))
	api.GET("/download/:id", handlers.DownloadFinal(db, store))
	api.GET("/download/:id/secure", handlers.ValidateDownloadToken(db, store))
	api.POST("/sessions/:id/generate-link", handlers.GenerateSecureDownloadURL(db))
	api.POST("/sessions/:id/auto-print", handlers.AutoPrintJob(db, store))
	api.POST("/sessions/:id/send-email", handlers.SendEmailDelivery(db, store))
	api.POST("/sessions/:id/resend-email", handlers.ResendEmail(db, store))

	// Protected routes
	auth := api.Group("/")
	auth.Use(func(c *gin.Context) {
		middleware.AuthMiddleware(c)
	})
	auth.GET("/auth/me", handlers.GetCurrentUser(db))

	// Merchant-scoped routes
	merchant := auth.Group("/")
	merchant.Use(func(c *gin.Context) {
		middleware.TenantMiddleware(c)
	})

	// Campaigns
	merchant.GET("/campaigns", handlers.ListCampaigns(db))
	merchant.POST("/campaigns", handlers.CreateCampaign(db))
	merchant.GET("/campaigns/:id", handlers.GetCampaign(db))
	merchant.PUT("/campaigns/:id", handlers.UpdateCampaign(db))
	merchant.DELETE("/campaigns/:id", handlers.DeleteCampaign(db))

	// Templates
	merchant.GET("/templates", handlers.ListTemplates(db))
	merchant.POST("/templates", handlers.CreateTemplate(db))
	merchant.GET("/templates/:id", handlers.GetTemplate(db))
	merchant.PUT("/templates/:id", handlers.UpdateTemplate(db))
	merchant.DELETE("/templates/:id", handlers.DeleteTemplate(db))

	// Sessions & Gallery
	merchant.GET("/sessions", handlers.ListSessions(db))
	merchant.GET("/sessions/:id/detail", handlers.GetSessionDetail(db))
	merchant.DELETE("/sessions/:id", handlers.DeleteSession(db))

	// Print Jobs
	merchant.GET("/print-jobs", handlers.ListPrintJobs(db))
	merchant.GET("/print-jobs/:id", handlers.GetPrintJob(db))
	merchant.POST("/sessions/:id/print", handlers.CreatePrintJob(db, store))
	merchant.PUT("/print-jobs/:id/status", handlers.UpdatePrintJobStatus(db))
	merchant.POST("/sessions/:id/reprint", handlers.Reprint(db, store))
	merchant.GET("/print-jobs/pending", handlers.GetPendingPrintJobs(db))

	// Branding
	merchant.GET("/branding", handlers.GetBranding(db))
	merchant.PUT("/branding", handlers.UpdateBranding(db, store))

	// Devices
	merchant.GET("/devices", handlers.ListDevices(db))
	merchant.POST("/devices", handlers.CreateDevice(db))
	merchant.GET("/devices/:id", handlers.GetDevice(db))
	merchant.PUT("/devices/:id", handlers.UpdateDevice(db))
	merchant.DELETE("/devices/:id", handlers.DeleteDevice(db))

	// Analytics
	merchant.GET("/analytics/overview", handlers.GetAnalyticsOverview(db))

	// Super Admin routes
	superAdmin := auth.Group("/admin")
	superAdmin.Use(func(c *gin.Context) {
		middleware.SuperAdminMiddleware(c)
	})
	superAdmin.GET("/merchants", handlers.AdminListMerchants(db))
	superAdmin.POST("/merchants", handlers.AdminCreateMerchant(db))
	superAdmin.PUT("/merchants/:id", handlers.AdminUpdateMerchant(db))
	superAdmin.DELETE("/merchants/:id", handlers.AdminDeleteMerchant(db))
	superAdmin.GET("/sessions", handlers.AdminListAllSessions(db))
	superAdmin.GET("/analytics", handlers.AdminGetPlatformAnalytics(db))

	// Upload (merchant-scoped)
	merchant.POST("/upload", handlers.UploadFile(db, store))
}
