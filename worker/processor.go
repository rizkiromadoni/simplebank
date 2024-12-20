package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/rizkiromadoni/simplebank/db/sqlc"
	"github.com/rizkiromadoni/simplebank/mail"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskSendVerifyEmail(c context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(opts asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProcessor {
	server := asynq.NewServer(opts, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("error while processing task")
		}),
		Logger: NewLogger(),
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mailer: mailer,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, p.ProcessTaskSendVerifyEmail)
	return p.server.Start(mux)
}

func (p *RedisTaskProcessor) Shutdown() {
	p.server.Shutdown()
}
