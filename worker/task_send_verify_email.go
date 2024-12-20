package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	db "github.com/rizkiromadoni/simplebank/db/sqlc"
	"github.com/rizkiromadoni/simplebank/util"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task_send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (d *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	c context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := d.client.EnqueueContext(c, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (p *RedisTaskProcessor) ProcessTaskSendVerifyEmail(c context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	user, err := p.store.GetUser(c, payload.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user not found: %w", err)
		}

		return fmt.Errorf("failed to get user: %w", err)
	}

	verifyEmail, err := p.store.CreateVerifyEmail(c, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	// TODO: send verify email
	verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s",
		verifyEmail.ID, verifyEmail.SecretCode)
	subject := "verify your email"
	content := fmt.Sprintf(`Hello %s, please verify your email address by clicking the link below.\n\n<a href="%s">Click Here</a>\n\nIf you did not request this, please ignore this email.\n`, user.Username, verifyUrl)
	to := []string{user.Email}

	err = p.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", user.Email).Msg("processed task")
	return nil
}
