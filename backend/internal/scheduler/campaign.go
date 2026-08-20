package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CampaignScheduler periodically checks and updates campaign status
type CampaignScheduler struct {
	db     *pgxpool.Pool
	ticker *time.Ticker
	done   chan bool
}

// NewCampaignScheduler creates a new scheduler
func NewCampaignScheduler(db *pgxpool.Pool) *CampaignScheduler {
	return &CampaignScheduler{
		db:   db,
		done: make(chan bool),
	}
}

// Start begins the scheduling loop
func (s *CampaignScheduler) Start() {
	s.ticker = time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-s.done:
				s.ticker.Stop()
				return
			case <-s.ticker.C:
				s.checkCampaigns()
			}
		}
	}()
	log.Println("✅ Campaign scheduler started")
}

// Stop ends the scheduling loop
func (s *CampaignScheduler) Stop() {
	close(s.done)
}

func (s *CampaignScheduler) checkCampaigns() {
	ctx := context.Background()

	// Activate campaigns that have started
	_, err := s.db.Exec(ctx, `
		UPDATE campaigns 
		SET status = 'active', updated_at = NOW()
		WHERE status = 'scheduled' 
		AND start_date IS NOT NULL 
		AND start_date <= NOW()
	`)
	if err != nil {
		log.Printf("Scheduler: activate error: %v", err)
	}

	// Expire campaigns that have ended
	_, err = s.db.Exec(ctx, `
		UPDATE campaigns 
		SET status = 'expired', updated_at = NOW()
		WHERE status = 'active' 
		AND end_date IS NOT NULL 
		AND end_date < NOW()
	`)
	if err != nil {
		log.Printf("Scheduler: expire error: %v", err)
	}
}
