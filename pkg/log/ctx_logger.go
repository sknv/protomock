package log

import (
	"context"
	"log/slog"
)

type ctxKey string

type ctxFields struct {
	fields []slog.Attr
}

const (
	_slogLogger ctxKey = "slog_logger"
	_slogFields ctxKey = "slog_fields"
)

// ToContext returns a context that contains the given Logger.
// Use FromContext to retrieve the Logger.
func ToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, _slogLogger, logger)
}

// FromContext returns the Logger stored in ctx by ToContext, or the default Logger if there is none.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(_slogLogger).(*slog.Logger); ok && logger != nil {
		return logger
	}

	return slog.Default()
}

type ContextHandler struct {
	slog.Handler
}

func NewContextHandler(handler slog.Handler) *ContextHandler {
	return &ContextHandler{
		Handler: handler,
	}
}

// Handle adds contextual attributes to the Record before calling the underlying handler.
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if fields, ok := ctx.Value(_slogFields).(*ctxFields); ok && fields != nil {
		r.AddAttrs(fields.fields...)
	}

	return h.Handler.Handle(ctx, r) //nolint:wrapcheck // proxy
}

// AppendCtx adds an slog attribute to the provided context so that it will be
// included in any Record created with such context.
func AppendCtx(ctx context.Context, attrs ...slog.Attr) context.Context {
	fields, ok := ctx.Value(_slogFields).(*ctxFields)
	if !ok || fields == nil {
		fields = &ctxFields{fields: nil}
	}

	fields.fields = append(fields.fields, attrs...)

	return context.WithValue(ctx, _slogFields, fields)
}
