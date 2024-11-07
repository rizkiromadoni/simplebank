package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/rizkiromadoni/simplebank/db/sqlc"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(c context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(opts asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(opts, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, p.ProcessTaskSendVerifyEmail)
	return p.server.Start(mux)
}
