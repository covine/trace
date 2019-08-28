package trace

import (
	"errors"

	"github.com/opentracing/opentracing-go"
)

type Closer interface {
	Close() error
}

type Tracer struct {
	tracer opentracing.Tracer // nil tracer 可以关闭整个链路追踪的功能，不用担心程序panic，Tracer封装会做好保护
	closer Closer             // opentracing.Tracer 没有Close接口，有些情况下，比如日志文件是需要关闭的
}

func NewTracer(t opentracing.Tracer, c Closer) (*Tracer, error) {
	if t == nil {
		return nil, errors.New("nil opentracing tracer")
	}
	return &Tracer{
		tracer: t,
		closer: c,
	}, nil
}

func NoopTracer() *Tracer {
	return &Tracer{
		tracer: nil,
		closer: nil,
	}
}

func (t *Tracer) isNoop() bool {
	return t == nil || t.tracer == nil
}

// 假如A函数内部支持链路追踪，如果使用原生的startSpan，那么在上游span nil的情况下，A函数会启动一个root span
// 我们希望控制A函数，使其按照预期打开或关闭链路追踪，那么我们可以通过控制span 是否 nil，来控制A函数是否启动链路追踪
// 因此，需要区分 StartRootSpan 和 StartSpan
// StartRootSpan 即便是 parent span 为 nil，也会创建一个 root span，StartSpan 则忽略这次调用，不会创建可用 span
// 所以，程序员需要自己知道，A函数是不是一个链路的开始，如果是，就要用StartRootSpan，即便上游没传有效span，也会创建一个root span
// 如果A函数只是一个中间过程，上游传入nil span的时候，不要做任何处理，就要用 StartSpan
func (t *Tracer) StartRootSpan(name string, ref opentracing.SpanReferenceType, s *Span, options ...Options) *Span {
	if t.isNoop() {
		return nil
	}
	return t.startSpan(name, ref, s.Context(), options...)
}

func (t *Tracer) StartSpan(name string, ref opentracing.SpanReferenceType, s *Span, options ...Options) *Span {
	if t.isNoop() {
		return nil
	}
	if s.isNoop() {
		return nil
	}
	return t.startSpan(name, ref, s.Context(), options...)
}

func (t *Tracer) startSpan(name string, ref opentracing.SpanReferenceType, ctx opentracing.SpanContext, options ...Options) *Span {
	var sp opentracing.Span
	if ref == opentracing.FollowsFromRef {
		sp = t.tracer.StartSpan(name, opentracing.FollowsFrom(ctx))
	} else {
		sp = t.tracer.StartSpan(name, opentracing.ChildOf(ctx))
	}

	if sp == nil {
		return nil
	} else {
		rs := &Span{
			span: sp,
		}

		for _, op := range options {
			op(rs)
		}
		return rs
	}
}

func (t *Tracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	if t != nil && t.tracer != nil {
		return t.tracer.Inject(sm, format, carrier)
	}
	return nil
}

func (t *Tracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	if t != nil && t.tracer != nil {
		return t.tracer.Extract(format, carrier)
	}
	return nil, nil
}

func (t *Tracer) OpentracingTracer() opentracing.Tracer {
	if t.isNoop() {
		return nil
	}
	return t.tracer
}

func (t *Tracer) SetGlobalTracer() {
	if t.isNoop() {
		return
	}
	opentracing.SetGlobalTracer(t.tracer)
}

func (t *Tracer) Close() error {
	if t != nil && t.closer != nil {
		return t.closer.Close()
	}
	return nil
}
