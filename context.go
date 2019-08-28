package trace

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

func ContextWithSpan(ctx context.Context, span *Span) context.Context {
	if span.isNoop() {
		return ctx
	}

	return opentracing.ContextWithSpan(ctx, span.span)
}
