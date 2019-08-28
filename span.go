package trace

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func NewSpanFromContext(ctx context.Context) *Span {
	upstreamSpan := opentracing.SpanFromContext(ctx)
	return &Span{
		span: upstreamSpan,
	}
}

type Span struct {
	span opentracing.Span
}

func (s *Span) isNoop() bool {
	return s == nil || s.span == nil
}

func (s *Span) Finish() {
	if !s.isNoop() {
		s.span.Finish()
	}
}

func (s *Span) FinishWithOptions(opts opentracing.FinishOptions) {
	if !s.isNoop() {
		s.span.FinishWithOptions(opts)
	}
}

func (s *Span) Context() opentracing.SpanContext {
	if !s.isNoop() {
		return s.span.Context()
	}
	return nil
}

func (s *Span) SetOperationName(operationName string) *Span {
	if !s.isNoop() {
		s.span.SetOperationName(operationName)
	}
	return s
}

func (s *Span) SetTag(key, value string) *Span {
	if !s.isNoop() {
		s.span.SetTag(key, value)
	}
	return s
}

func (s *Span) SetError(err error) *Span {
	if !s.isNoop() {
		s.span.SetTag("error", err.Error())
	}
	return s
}

func (s *Span) LogFields(fields ...log.Field) {
	if !s.isNoop() {
		s.span.LogFields(fields...)
	}
}

func (s *Span) LogKV(alternatingKeyValues ...interface{}) {
	if !s.isNoop() {
		s.span.LogKV(alternatingKeyValues...)
	}
}

func (s *Span) SetBaggageItem(restrictedKey, value string) *Span {
	if !s.isNoop() {
		s.span.SetBaggageItem(restrictedKey, value)
	}
	return s
}

func (s *Span) BaggageItem(restrictedKey string) string {
	if !s.isNoop() {
		return s.span.BaggageItem(restrictedKey)
	}
	return ""
}

func (s *Span) Tracer() opentracing.Tracer {
	if !s.isNoop() {
		return s.span.Tracer()
	}
	return nil
}

func (s *Span) SetRPCClient() {
	if !s.isNoop() {
		ext.SpanKindRPCClient.Set(s.span)
	}
}

func (s *Span) SetRPCServer() {
	if !s.isNoop() {
		ext.SpanKindRPCServer.Set(s.span)
	}
}

func (s *Span) SetProducer() {
	if !s.isNoop() {
		ext.SpanKindProducer.Set(s.span)
	}
}

func (s *Span) SetConsumer() {
	if !s.isNoop() {
		ext.SpanKindConsumer.Set(s.span)
	}
}

func (s *Span) SetDBType(t string) {
	if !s.isNoop() {
		ext.DBType.Set(s.span, t)
	}
}

func (s *Span) SetPeerService(p string) {
	if !s.isNoop() {
		ext.PeerService.Set(s.span, p)
	}
}

func (s *Span) SetDBInstance(i string) {
	if !s.isNoop() {
		ext.DBInstance.Set(s.span, i)
	}
}

func (s *Span) SetDBStatement(m string) {
	if !s.isNoop() {
		ext.DBStatement.Set(s.span, m)
	}
}

func (s *Span) SetPeerHostname(h string) {
	if !s.isNoop() {
		ext.PeerHostname.Set(s.span, h)
	}
}

func (s *Span) SetPeerPort(p uint16) {
	if !s.isNoop() {
		ext.PeerPort.Set(s.span, p)
	}
}
