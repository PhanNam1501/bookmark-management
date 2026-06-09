package queue

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrorNoMessage = errors.New("no message")

type redisQueue struct {
	client    *redis.Client
	queueName string
}

func NewRedisQueue(client *redis.Client, queueName string) Repository {
	return &redisQueue{
		client:    client,
		queueName: queueName,
	}
}

func (r *redisQueue) PushMessage(ctx context.Context, message []byte) error {
	return r.client.LPush(ctx, r.queueName, message).Err()
}

func (r *redisQueue) PopMessage(ctx context.Context) ([]byte, error) {
	mes, err := r.client.RPop(ctx, r.queueName).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrorNoMessage
	}
	return mes, nil
}
