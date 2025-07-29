package scheduler

import (
	"context"
	"time"
)

// Job represents a scheduled task
type Job interface {
	// ID returns the unique identifier for this job
	ID() string
	
	// Name returns the human-readable name for this job
	Name() string
	
	// Schedule returns the cron expression for this job
	Schedule() string
	
	// Execute runs the job with the provided context
	Execute(ctx context.Context) error
	
	// OnSuccess is called when the job completes successfully
	OnSuccess(duration time.Duration)
	
	// OnError is called when the job fails
	OnError(err error, duration time.Duration)
}

// JobScheduler manages and executes scheduled jobs
type JobScheduler interface {
	// Start begins the job scheduler
	Start(ctx context.Context) error
	
	// Stop gracefully shuts down the job scheduler
	Stop() error
	
	// AddJob registers a new job with the scheduler
	AddJob(job Job) error
	
	// RemoveJob unregisters a job from the scheduler
	RemoveJob(jobID string) error
	
	// GetJob retrieves a job by ID
	GetJob(jobID string) (Job, bool)
	
	// ListJobs returns all registered jobs
	ListJobs() []Job
	
	// IsRunning returns true if the scheduler is currently running
	IsRunning() bool
}

// JobExecution represents a single execution of a job
type JobExecution struct {
	JobID     string        `json:"job_id"`
	JobName   string        `json:"job_name"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Status    string        `json:"status"` // "success", "error", "running"
	Error     string        `json:"error,omitempty"`
}

// JobStats contains statistics about job executions
type JobStats struct {
	JobID            string        `json:"job_id"`
	JobName          string        `json:"job_name"`
	TotalExecutions  int           `json:"total_executions"`
	SuccessfulRuns   int           `json:"successful_runs"`
	FailedRuns       int           `json:"failed_runs"`
	LastExecution    time.Time     `json:"last_execution"`
	LastSuccess      time.Time     `json:"last_success"`
	LastError        string        `json:"last_error,omitempty"`
	AverageDuration  time.Duration `json:"average_duration"`
	NextScheduled    time.Time     `json:"next_scheduled"`
}

// BaseJob provides a basic implementation of the Job interface
type BaseJob struct {
	id       string
	name     string
	schedule string
}

// NewBaseJob creates a new base job
func NewBaseJob(id, name, schedule string) *BaseJob {
	return &BaseJob{
		id:       id,
		name:     name,
		schedule: schedule,
	}
}

// ID returns the job ID
func (b *BaseJob) ID() string {
	return b.id
}

// Name returns the job name
func (b *BaseJob) Name() string {
	return b.name
}

// Schedule returns the cron schedule
func (b *BaseJob) Schedule() string {
	return b.schedule
}

// OnSuccess default implementation - can be overridden
func (b *BaseJob) OnSuccess(duration time.Duration) {
	// Default implementation does nothing
}

// OnError default implementation - can be overridden
func (b *BaseJob) OnError(err error, duration time.Duration) {
	// Default implementation does nothing
}