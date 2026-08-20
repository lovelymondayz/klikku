package models

import "time"

// Roles
const (
	RoleSuperAdmin  = "SUPER_ADMIN"
	RoleMerchantAdmin = "MERCHANT_ADMIN"
	RoleMerchantStaff = "MERCHANT_STAFF"
)

// User represents a system user (super admin, merchant admin/staff)
type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	MerchantID   string    `json:"merchant_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// Merchant represents a business using the platform
type Merchant struct {
	ID               string    `json:"id"`
	BusinessName     string    `json:"business_name"`
	Slug             string    `json:"slug"`
	LogoURL          string    `json:"logo_url,omitempty"`
	PrimaryColor     string    `json:"primary_color,omitempty"`
	SecondaryColor   string    `json:"secondary_color,omitempty"`
	Font             string    `json:"font,omitempty"`
	WelcomeMessage   string    `json:"welcome_message,omitempty"`
	IdleBackgroundURL string   `json:"idle_background_url,omitempty"`
	SubscriptionStatus string  `json:"subscription_status"`
	EmailDesign      JSONMap   `json:"email_design,omitempty"`
	SocialLinks      JSONMap   `json:"social_links,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

// Campaign represents a promotional campaign
type Campaign struct {
	ID              string    `json:"id"`
	MerchantID      string    `json:"merchant_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Status          string    `json:"status"`
	PromotionConfig JSONMap   `json:"promotion_config,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// Template represents a photobooth template
type Template struct {
	ID            string    `json:"id"`
	MerchantID    string    `json:"merchant_id"`
	CampaignID    string    `json:"campaign_id,omitempty"`
	Name          string    `json:"name"`
	PreviewURL    string    `json:"preview_url,omitempty"`
	LayoutConfig  JSONMap   `json:"layout_config"`
	OverlayURL    string    `json:"overlay_url,omitempty"`
	OutputWidth   int       `json:"output_width"`
	OutputHeight  int       `json:"output_height"`
	PhotoCount    int       `json:"photo_count"`
	Price         float64   `json:"price"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"created_at"`
}

// Device represents a physical photobooth device
type Device struct {
	ID             string    `json:"id"`
	MerchantID     string    `json:"merchant_id"`
	Name           string    `json:"name"`
	DeviceToken    string    `json:"device_token"`
	Status         string    `json:"status"`
	CurrentCampaignID string `json:"current_campaign_id,omitempty"`
	PrinterConfig  JSONMap   `json:"printer_config,omitempty"`
	LastSeenAt     time.Time `json:"last_seen_at"`
	CreatedAt      time.Time `json:"created_at"`
}

// PhotoboothSession represents a customer photo session
type PhotoboothSession struct {
	ID            string    `json:"id"`
	MerchantID    string    `json:"merchant_id"`
	DeviceID      string    `json:"device_id,omitempty"`
	CampaignID    string    `json:"campaign_id,omitempty"`
	TemplateID    string    `json:"template_id,omitempty"`
	Status        string    `json:"status"`
	PaymentStatus string    `json:"payment_status"`
	Email         string    `json:"email,omitempty"`
	FinalImageURL string    `json:"final_image_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	CompletedAt   time.Time `json:"completed_at,omitempty"`
}

// Photo represents a captured or processed photo
type Photo struct {
	ID           string    `json:"id"`
	SessionID    string    `json:"session_id"`
	OriginalURL  string    `json:"original_url,omitempty"`
	ProcessedURL string    `json:"processed_url,omitempty"`
	FinalURL     string    `json:"final_url,omitempty"`
	Position     int       `json:"position"`
	CreatedAt    time.Time `json:"created_at"`
}

// Payment represents a payment transaction
type Payment struct {
	ID                   string    `json:"id"`
	SessionID            string    `json:"session_id"`
	Provider             string    `json:"provider"`
	Amount               float64   `json:"amount"`
	Status               string    `json:"status"`
	TransactionReference string    `json:"transaction_reference,omitempty"`
	PaidAt               time.Time `json:"paid_at,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

// PrintJob represents a print task
type PrintJob struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"session_id"`
	DeviceID    string    `json:"device_id"`
	Status      string    `json:"status"`
	Copies      int       `json:"copies"`
	PrinterName string    `json:"printer_name,omitempty"`
	ErrorMessage string   `json:"error_message,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// EmailDelivery represents an email delivery record
type EmailDelivery struct {
	ID           string    `json:"id"`
	SessionID    string    `json:"session_id"`
	Email        string    `json:"email"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message,omitempty"`
	SentAt       time.Time `json:"sent_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// JSONMap is a flexible JSON object
type JSONMap map[string]interface{}
