package synq

import (
	"context"
	"encoding/json"
	"io"

	"buf.build/go/protovalidate"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

type Client interface {
	io.Closer
	Ping() error
	Enqueue(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
	EnqueueContext(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
	EnqueuePbTask(ctx context.Context, pattern string, m proto.Message, opts ...asynq.Option) (*asynq.TaskInfo, error)
	EnqueueJSONTask(ctx context.Context, pattern string, m any, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

type client struct {
	client *asynq.Client
}

func NewClient(c asynq.RedisConnOpt) Client {
	return &client{client: asynq.NewClient(c)}
}

func NewClientFromRedisClient(c redis.UniversalClient) Client {
	return &client{client: asynq.NewClientFromRedisClient(c)}
}

func NewTask(pattern string, payload []byte, opts ...asynq.Option) *asynq.Task {
	return asynq.NewTask(pattern, payload, opts...)
}

func NewPbTask(pattern string, m proto.Message, opts ...asynq.Option) (*asynq.Task, error) {
	err := protovalidate.Validate(m)
	if err != nil {
		return nil, err
	}
	py, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(pattern, py, opts...), nil
}

func NewJSONTask(pattern string, m any, opts ...asynq.Option) (*asynq.Task, error) {
	py, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(pattern, py, opts...), nil
}

func (c *client) Enqueue(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return c.client.Enqueue(task, opts...)
}

func (c *client) EnqueueContext(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return c.client.EnqueueContext(ctx, task, opts...)
}

func (c *client) EnqueuePbTask(ctx context.Context, pattern string, m proto.Message, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	task, err := NewPbTask(pattern, m)
	if err != nil {
		return nil, err
	}
	return c.client.EnqueueContext(ctx, task, opts...)
}
func (c *client) EnqueueJSONTask(ctx context.Context, pattern string, m any, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	task, err := NewJSONTask(pattern, m)
	if err != nil {
		return nil, err
	}
	return c.client.EnqueueContext(ctx, task, opts...)
}

func (c *client) Ping() error {
	return c.client.Ping()
}

func (c *client) Close() error {
	return c.client.Close()
}
