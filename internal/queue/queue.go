package queue

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

// Queue defines the interface for job queue operations
type Queue interface {
	// Client operations
	EnqueueSendContent(contentID, jobID string) (*asynq.TaskInfo, error)
	Close() error

	// Server operations
	RegisterHandler(taskType string, handler asynq.HandlerFunc)
	Start() error
	Stop()
	Shutdown()
}

// TaskHandler defines the interface for task handlers
type TaskHandler interface {
	HandleSendContent(ctx context.Context, task *asynq.Task) error
}

// AsynqQueue implements the Queue interface using Asynq
type AsynqQueue struct {
	client *asynq.Client
	server *asynq.Server
	mux    *asynq.ServeMux
	logger interface{} // Using interface{} to avoid zap dependency in interface
}

// NewAsynqQueue creates a new Asynq-based queue
func NewAsynqQueue(redisAddr, redisPassword string, redisDB int, logger interface{}) *AsynqQueue {
	redisOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	}

	client := asynq.NewClient(redisOpt)
	
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()

	return &AsynqQueue{
		client: client,
		server: server,
		mux:    mux,
		logger: logger,
	}
}

// EnqueueSendContent enqueues a send content task
func (q *AsynqQueue) EnqueueSendContent(contentID, jobID string) (*asynq.TaskInfo, error) {
	payload := map[string]interface{}{
		"content_id": contentID,
		"job_id":     jobID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	task := asynq.NewTask("send_newsletter", payloadBytes)
	return q.client.Enqueue(task)
}

// Close closes the client connection
func (q *AsynqQueue) Close() error {
	return q.client.Close()
}

// RegisterHandler registers a task handler
func (q *AsynqQueue) RegisterHandler(taskType string, handler asynq.HandlerFunc) {
	q.mux.HandleFunc(taskType, handler)
}

// Start starts the Asynq server
func (q *AsynqQueue) Start() error {
	return q.server.Start(q.mux)
}

// Stop stops the Asynq server
func (q *AsynqQueue) Stop() {
	q.server.Stop()
}

// Shutdown gracefully shuts down the server
func (q *AsynqQueue) Shutdown() {
	q.server.Shutdown()
}
