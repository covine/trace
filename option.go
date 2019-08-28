package trace

import (
	"github.com/opentracing/opentracing-go/ext"
)

type Options func(s *Span)

func SpanKindRPCClient(s *Span) {
	if !s.isNoop() {
		ext.SpanKindRPCClient.Set(s.span)
	}
}

func SpanKindRPCServer(s *Span) {
	if !s.isNoop() {
		ext.SpanKindRPCServer.Set(s.span)
	}
}

func SpanKindProducer(s *Span) {
	if !s.isNoop() {
		ext.SpanKindProducer.Set(s.span)
	}
}

func SpanKindConsumer(s *Span) {
	if !s.isNoop() {
		ext.SpanKindConsumer.Set(s.span)
	}
}
