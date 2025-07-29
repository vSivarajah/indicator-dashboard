package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"crypto-indicator-dashboard/pkg/logger"

	"github.com/robfig/cron/v3"
)

// CronScheduler implements JobScheduler using the robfig/cron library
type CronScheduler struct {
	cron        *cron.Cron
	jobs        map[string]Job
	cronEntries map[string]cron.EntryID
	executions  map[string][]*JobExecution
	stats       map[string]*JobStats
	logger      logger.Logger
	mu          sync.RWMutex
	isRunning   bool
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewCronScheduler creates a new cron-based job scheduler
func NewCronScheduler(log logger.Logger) *CronScheduler {
	return &CronScheduler{
		cron:        cron.New(cron.WithSeconds()),
		jobs:        make(map[string]Job),
		cronEntries: make(map[string]cron.EntryID),
		executions:  make(map[string][]*JobExecution),
		stats:       make(map[string]*JobStats),
		logger:      log,
	}
}

// Start begins the job scheduler
func (cs *CronScheduler) Start(ctx context.Context) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if cs.isRunning {
		return fmt.Errorf("scheduler is already running")
	}

	cs.ctx, cs.cancel = context.WithCancel(ctx)
	cs.cron.Start()
	cs.isRunning = true

	cs.logger.Info("Job scheduler started")
	return nil
}

// Stop gracefully shuts down the job scheduler
func (cs *CronScheduler) Stop() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if !cs.isRunning {
		return fmt.Errorf("scheduler is not running")
	}

	if cs.cancel != nil {
		cs.cancel()
	}

	stopCtx := cs.cron.Stop()
	<-stopCtx.Done()

	cs.isRunning = false
	cs.logger.Info("Job scheduler stopped")
	return nil
}

// AddJob registers a new job with the scheduler
func (cs *CronScheduler) AddJob(job Job) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	jobID := job.ID()

	// Check if job already exists
	if _, exists := cs.jobs[jobID]; exists {
		return fmt.Errorf("job with ID '%s' already exists", jobID)
	}

	// Validate cron schedule
	_, err := cron.ParseStandard(job.Schedule())
	if err != nil {
		return fmt.Errorf("invalid cron schedule '%s': %w", job.Schedule(), err)
	}

	// Wrap the job execution with monitoring and error handling
	wrappedJob := cs.wrapJob(job)

	// Add to cron
	entryID, err := cs.cron.AddFunc(job.Schedule(), wrappedJob)
	if err != nil {
		return fmt.Errorf("failed to add job to cron: %w", err)
	}

	// Store job and entry ID
	cs.jobs[jobID] = job
	cs.cronEntries[jobID] = entryID
	cs.executions[jobID] = make([]*JobExecution, 0)
	cs.stats[jobID] = &JobStats{
		JobID:   jobID,
		JobName: job.Name(),
	}

	cs.logger.Info("Job added to scheduler",
		"job_id", jobID,
		"job_name", job.Name(),
		"schedule", job.Schedule())

	return nil
}

// RemoveJob unregisters a job from the scheduler
func (cs *CronScheduler) RemoveJob(jobID string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	entryID, exists := cs.cronEntries[jobID]
	if !exists {
		return fmt.Errorf("job with ID '%s' not found", jobID)
	}

	// Remove from cron
	cs.cron.Remove(entryID)

	// Clean up
	delete(cs.jobs, jobID)
	delete(cs.cronEntries, jobID)
	delete(cs.executions, jobID)
	delete(cs.stats, jobID)

	cs.logger.Info("Job removed from scheduler", "job_id", jobID)
	return nil
}

// GetJob retrieves a job by ID
func (cs *CronScheduler) GetJob(jobID string) (Job, bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	job, exists := cs.jobs[jobID]
	return job, exists
}

// ListJobs returns all registered jobs
func (cs *CronScheduler) ListJobs() []Job {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	jobs := make([]Job, 0, len(cs.jobs))
	for _, job := range cs.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// IsRunning returns true if the scheduler is currently running
func (cs *CronScheduler) IsRunning() bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.isRunning
}

// GetJobStats returns statistics for a specific job
func (cs *CronScheduler) GetJobStats(jobID string) (*JobStats, bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	stats, exists := cs.stats[jobID]
	if !exists {
		return nil, false
	}

	// Create a copy to avoid race conditions
	statsCopy := *stats
	return &statsCopy, true
}

// GetAllJobStats returns statistics for all jobs
func (cs *CronScheduler) GetAllJobStats() map[string]*JobStats {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	result := make(map[string]*JobStats)
	for jobID, stats := range cs.stats {
		statsCopy := *stats
		result[jobID] = &statsCopy
	}
	return result
}

// GetJobExecutions returns execution history for a specific job
func (cs *CronScheduler) GetJobExecutions(jobID string, limit int) ([]*JobExecution, bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	executions, exists := cs.executions[jobID]
	if !exists {
		return nil, false
	}

	// Return the most recent executions
	start := 0
	if len(executions) > limit {
		start = len(executions) - limit
	}

	result := make([]*JobExecution, len(executions)-start)
	copy(result, executions[start:])
	return result, true
}

// wrapJob wraps a job with monitoring and error handling
func (cs *CronScheduler) wrapJob(job Job) func() {
	return func() {
		// Check if scheduler is still running
		select {
		case <-cs.ctx.Done():
			return
		default:
		}

		jobID := job.ID()
		startTime := time.Now()

		execution := &JobExecution{
			JobID:     jobID,
			JobName:   job.Name(),
			StartTime: startTime,
			Status:    "running",
		}

		cs.logger.Info("Starting job execution",
			"job_id", jobID,
			"job_name", job.Name())

		// Execute the job
		err := job.Execute(cs.ctx)

		endTime := time.Now()
		duration := endTime.Sub(startTime)

		// Update execution record
		execution.EndTime = endTime
		execution.Duration = duration

		if err != nil {
			execution.Status = "error"
			execution.Error = err.Error()
			job.OnError(err, duration)

			cs.logger.Error("Job execution failed",
				"job_id", jobID,
				"job_name", job.Name(),
				"duration", duration,
				"error", err)
		} else {
			execution.Status = "success"
			job.OnSuccess(duration)

			cs.logger.Info("Job execution completed successfully",
				"job_id", jobID,
				"job_name", job.Name(),
				"duration", duration)
		}

		// Update statistics and execution history
		cs.updateJobStats(jobID, execution)
	}
}

// updateJobStats updates job statistics and execution history
func (cs *CronScheduler) updateJobStats(jobID string, execution *JobExecution) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Add to execution history (keep last 100 executions)
	executions := cs.executions[jobID]
	executions = append(executions, execution)
	if len(executions) > 100 {
		executions = executions[1:]
	}
	cs.executions[jobID] = executions

	// Update statistics
	stats := cs.stats[jobID]
	stats.TotalExecutions++
	stats.LastExecution = execution.EndTime

	if execution.Status == "success" {
		stats.SuccessfulRuns++
		stats.LastSuccess = execution.EndTime
	} else {
		stats.FailedRuns++
		stats.LastError = execution.Error
	}

	// Calculate average duration
	totalDuration := time.Duration(0)
	for _, exec := range executions {
		if exec.Status != "running" {
			totalDuration += exec.Duration
		}
	}
	if stats.TotalExecutions > 0 {
		stats.AverageDuration = totalDuration / time.Duration(stats.TotalExecutions)
	}

	// Calculate next scheduled time
	if entryID, exists := cs.cronEntries[jobID]; exists {
		entry := cs.cron.Entry(entryID)
		stats.NextScheduled = entry.Next
	}
}
