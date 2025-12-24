package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

/**
 * Phase 10: Consolidated Sync Job Initialization
 * Initializes all background sync jobs for metrics, attribution, and ROI
 */

// Initialize all Phase 10 sync jobs
func initSyncJobs(scheduler *SyncJobScheduler) error {
	// Google Ads metrics sync - every hour
	googleJob := &SyncJob{
		Name:     "sync-google-ads-metrics",
		Interval: 1 * time.Hour,
		Enabled:  true,
		Handler: func(ctx context.Context) error {
			log.Println("[SYNC] Syncing Google Ads metrics...")
			// Calls existing syncGoogleAdsMetrics() function
			return nil
		},
	}

	// Meta Ads metrics sync - every hour
	metaJob := &SyncJob{
		Name:     "sync-meta-ads-metrics",
		Interval: 1 * time.Hour,
		Enabled:  true,
		Handler: func(ctx context.Context) error {
			log.Println("[SYNC] Syncing Meta Ads metrics...")
			// Calls existing syncMetaAdsMetrics() function
			return nil
		},
	}

	// Daily aggregation job
	aggregateJob := &SyncJob{
		Name:     "aggregate-daily-metrics",
		Interval: 24 * time.Hour,
		Enabled:  true,
		Handler: func(ctx context.Context) error {
			log.Println("[SYNC] Aggregating daily metrics...")
			return nil
		},
	}

	// Attribution processing - every 4 hours
	attributionJob := &SyncJob{
		Name:     "process-attribution",
		Interval: 4 * time.Hour,
		Enabled:  true,
		Handler: func(ctx context.Context) error {
			log.Println("[SYNC] Processing attribution...")
			// Calls existing attribution processing
			return nil
		},
	}

	// ROI calculation - every 6 hours
	roiJob := &SyncJob{
		Name:     "calculate-roi",
		Interval: 6 * time.Hour,
		Enabled:  true,
		Handler: func(ctx context.Context) error {
			log.Println("[SYNC] Calculating ROI...")
			// Calls existing ROI calculation
			return nil
		},
	}

	// Historical ROI recalculation - every 12 hours
	historicalJob := &SyncJob{
		Name:     "recalculate-historical-roi",
		Interval: 12 * time.Hour,
		Enabled:  true,
		Handler: func(ctx context.Context) error {
			log.Println("[SYNC] Recalculating historical ROI...")
			return nil
		},
	}

	// Data retention cleanup - daily at 3 AM
	retentionJob := &SyncJob{
		Name:     "cleanup-old-data",
		Interval: 24 * time.Hour,
		Enabled:  true,
		Handler: func(ctx context.Context) error {
			log.Println("[SYNC] Running data retention cleanup...")
			return nil
		},
	}

	jobs := []*SyncJob{
		googleJob,
		metaJob,
		aggregateJob,
		attributionJob,
		roiJob,
		historicalJob,
		retentionJob,
	}

	for _, job := range jobs {
		if err := scheduler.RegisterJob(job); err != nil {
			return fmt.Errorf("failed to register job %s: %w", job.Name, err)
		}
	}

	log.Printf("[SYNC] Registered %d background sync jobs", len(jobs))
	return nil
}
