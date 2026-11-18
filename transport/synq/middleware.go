package synq

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
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

func Logger(log zerolog.Logger) asynq.MiddlewareFunc {
	return func(h asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
			var err error

			start := time.Now()
			defer func() {
				log.Info().Ctx(ctx).
					Str("pattern", t.Type()).
					Str("taskId", t.ResultWriter().TaskID()).
					Dur("latency", time.Since(start)).
					Func(func(e *zerolog.Event) {
						if err != nil {
							e.Err(err)
						}
					}).
					Msg("process task")
			}()
			err = h.ProcessTask(ctx, t)
			return err
		})
	}
}
