package tasks

import (
	"github.com/hibiken/asynq"
)

// A wrapper struct for the async client
type async struct {
	client *asynq.Client
}

// A wrapper interface for the async client
// allowing for easier testing and mocking
type Client interface {
	Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

// Initializes a new async client
// with the given redis address
func NewAsyncClient(rAddr string, rDB int) Client {
	return &async{
		client: asynq.NewClient(asynq.RedisClientOpt{Addr: rAddr, DB: rDB}),
	}
}

// Enqueues a task to the async worker
func (a *async) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	defer a.client.Close() // Close the client when done

	// Enqueue the task with the given options
	ti, err := a.client.Enqueue(task, opts...)
	if err != nil {
		return nil, err
	}

	return ti, nil
}
