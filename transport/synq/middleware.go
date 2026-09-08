package synq

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/thinkgos/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func TraceId(spanName string) asynq.MiddlewareFunc {
	return func(h asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
			ctx, span := otel.GetTracerProvider().
				Tracer("asynq server").
				Start(
					ctx,
					spanName,
					trace.WithSpanKind(trace.SpanKindServer),
				)
			defer span.End()
			return h.ProcessTask(ctx, t)
		})
	}
}

func Logger(log *logger.Log) asynq.MiddlewareFunc {
	return func(h asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
			var err error

			start := time.Now()
			defer func() {
				log.OnInfoContext(ctx).
					String("pattern", t.Type()).
					String("taskId", t.ResultWriter().TaskID()).
					Duration("latency", time.Since(start)).
					HookFuncIf(err != nil, func(e *logger.Event) {
						e.Error(err)
					}).
					Msg("process task")
			}()
			err = h.ProcessTask(ctx, t)
			return err
		})
	}
}
