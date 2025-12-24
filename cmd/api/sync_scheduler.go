package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

/**
 * Sync Job Scheduler
 * Manages background jobs for metrics synchronization
 * Uses cron-like scheduling for hourly/daily sync tasks
 */

type SyncJobScheduler struct {
	jobs map[string]*SyncJob
	done chan bool
}

type SyncJob struct {
	Name       string
	Interval   time.Duration
	LastRun    time.Time
	NextRun    time.Time
	Handler    func(ctx context.Context) error
	Enabled    bool
	RunCount   int
	ErrorCount int
	LastError  error
}

type SyncJobStats struct {
	Name       string
	LastRun    time.Time
	NextRun    time.Time
	RunCount   int
	ErrorCount int
	LastError  string
	Status     string
}

// Global scheduler instance
var scheduler *SyncJobScheduler

// Initialize the sync job scheduler
func initSyncJobScheduler() *SyncJobScheduler {
	return &SyncJobScheduler{
		jobs: make(map[string]*SyncJob),
		done: make(chan bool),
	}
}

// Register a new sync job
func (s *SyncJobScheduler) RegisterJob(job *SyncJob) error {
	if _, exists := s.jobs[job.Name]; exists {
		return fmt.Errorf("job %s already registered", job.Name)
	}

	job.NextRun = time.Now().Add(job.Interval)
	s.jobs[job.Name] = job

	log.Printf("[SCHEDULER] Registered job: %s (interval: %v)", job.Name, job.Interval)
	return nil
}

// Start the scheduler
func (s *SyncJobScheduler) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("[SCHEDULER] Scheduler stopped")
				s.done <- true
				return
			case <-ticker.C:
				s.checkAndRunJobs(ctx)
			}
		}
	}()

	log.Println("[SCHEDULER] Sync job scheduler started")
}

// Check and run jobs that are due
func (s *SyncJobScheduler) checkAndRunJobs(ctx context.Context) {
	now := time.Now()

	for jobName, job := range s.jobs {
		if !job.Enabled {
			continue
		}

		if now.After(job.NextRun) {
			go s.runJob(ctx, jobName, job)
		}
	}
}

// Run a single job
func (s *SyncJobScheduler) runJob(ctx context.Context, jobName string, job *SyncJob) {
	defer func() {
		job.LastRun = time.Now()
		job.NextRun = job.LastRun.Add(job.Interval)
		job.RunCount++

		if r := recover(); r != nil {
			job.ErrorCount++
			job.LastError = fmt.Errorf("panic: %v", r)
			log.Printf("[SCHEDULER] Job %s panicked: %v", jobName, r)
		}
	}()

	log.Printf("[SCHEDULER] Running job: %s", jobName)
	start := time.Now()

	err := job.Handler(ctx)
	duration := time.Since(start)

	if err != nil {
		job.ErrorCount++
		job.LastError = err
		log.Printf("[SCHEDULER] Job %s failed after %v: %v", jobName, duration, err)
	} else {
		log.Printf("[SCHEDULER] Job %s completed successfully in %v", jobName, duration)
	}
}

// Get job status
func (s *SyncJobScheduler) GetJobStatus(jobName string) (*SyncJobStats, error) {
	job, exists := s.jobs[jobName]
	if !exists {
		return nil, fmt.Errorf("job %s not found", jobName)
	}

	lastErrorStr := ""
	if job.LastError != nil {
		lastErrorStr = job.LastError.Error()
	}

	status := "idle"
	if job.Enabled {
		status = "enabled"
		if time.Now().Before(job.NextRun) {
			status = "waiting"
		}
	} else {
		status = "disabled"
	}

	return &SyncJobStats{
		Name:       job.Name,
		LastRun:    job.LastRun,
		NextRun:    job.NextRun,
		RunCount:   job.RunCount,
		ErrorCount: job.ErrorCount,
		LastError:  lastErrorStr,
		Status:     status,
	}, nil
}

// Get all job statuses
func (s *SyncJobScheduler) GetAllJobStatus() map[string]*SyncJobStats {
	stats := make(map[string]*SyncJobStats)
	for jobName := range s.jobs {
		if stat, err := s.GetJobStatus(jobName); err == nil {
			stats[jobName] = stat
		}
	}
	return stats
}

// Enable/disable a job
func (s *SyncJobScheduler) SetJobEnabled(jobName string, enabled bool) error {
	job, exists := s.jobs[jobName]
	if !exists {
		return fmt.Errorf("job %s not found", jobName)
	}
	job.Enabled = enabled
	log.Printf("[SCHEDULER] Job %s enabled: %v", jobName, enabled)
	return nil
}

// Stop the scheduler
func (s *SyncJobScheduler) Stop() {
	<-s.done
}
