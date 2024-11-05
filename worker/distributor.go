package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		c context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(opt *asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(opt)
	return &RedisTaskDistributor{
		client: client,
	}
}
