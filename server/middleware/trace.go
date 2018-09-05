package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hohice/gin-web/pkg/setting"
	stdot "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	zipkinot "github.com/openzipkin/zipkin-go-opentracing"
)

const spanContextKey = "span"

var (
	ErrSpanNotFound = errors.New("span was not found in context")
)

var tracer stdot.Tracer
var collector zipkinot.Collector

var (
	serviceName, url string
	port             int
)

func init() {
	registerSelf(func(conf *setting.Configs) (error, Closeble) {
		serviceName = conf.Service
		url = conf.Trace.ZipkinUrl
		port = conf.Http.HTTPPort
		return InitTracer()
	})
}

func InitTracer() (err error, close Closeble) {
	close = func() {
		if collector != nil {
			collector.Close()
		}
	}

	if url != "" {
		collector, err = zipkinot.NewHTTPCollector(url)
		if err != nil {
			return
		}

		var (
			debug       = false
			hostPort    = fmt.Sprintf("localhost:%d", port)
			serviceName = serviceName
		)
		recorder := zipkinot.NewRecorder(collector, debug, hostPort, serviceName)
		tracer, err = zipkinot.NewTracer(recorder)
		if err != nil {
			return
		}
	} else {
		err = errors.New("zipin url is none")
	}
	return
}

// NewSpan returns gin.HandlerFunc (middleware) that starts a new span and injects it to request context.
//
// It calls ctx.Next() to measure execution time of all following handlers.
func NewSpan(operationName string, opts ...stdot.StartSpanOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span := tracer.StartSpan(operationName, opts...)
		ext.HTTPMethod.Set(span, ctx.Request.Method)
		ext.HTTPUrl.Set(span, ctx.Request.URL.String())
		ctx.Set(spanContextKey, span)

		defer span.Finish()
		ctx.Next()

		//ext.HTTPStatusCode.Set(span,ctx.Writer.Status())
		//ctx.Set(spanContextKey, span)
	}
}

// ParentSpanReferenceFunc determines how to reference parent span
// See opentracing.SpanReferenceType
type ParentSpanReferenceFunc func(stdot.SpanContext) stdot.StartSpanOption

var CPsrFunc = func(spancontext stdot.SpanContext) stdot.StartSpanOption {
	return stdot.ChildOf(spancontext)
}

var FPsrFunc = func(spancontext stdot.SpanContext) stdot.StartSpanOption {
	return stdot.FollowsFrom(spancontext)
}

// SpanFromHeaders returns gin.HandlerFunc (middleware) that extracts parent span data from HTTP headers and
// starts a new span referenced to parent with ParentSpanReferenceFunc.
//
// It calls ctx.Next() to measure execution time of all following handlers.
//
// Behaviour on errors determined by abortOnErrors option. If it set to true request handling will be aborted with error.
func SpanFromHeaders(operationName string, psr ParentSpanReferenceFunc,
	abortOnErrors bool, advancedOpts ...stdot.StartSpanOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		spanContext, err := tracer.Extract(stdot.TextMap, stdot.HTTPHeadersCarrier(ctx.Request.Header))
		if err != nil {
			if abortOnErrors {
				ctx.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		opts := append([]stdot.StartSpanOption{psr(spanContext)}, advancedOpts...)

		span := tracer.StartSpan(operationName, opts...)
		ext.HTTPMethod.Set(span, ctx.Request.Method)
		ext.HTTPUrl.Set(span, ctx.Request.URL.String())
		ctx.Set(spanContextKey, span)
		defer span.Finish()

		ctx.Next()

		//ext.HTTPStatusCode.Set(span,ctx.Writer.Status())
		//ctx.Set(spanContextKey, span)
	}
}

// SpanFromContext returns gin.HandlerFunc (middleware) that extracts parent span from request context
// and starts a new span as child of parent span.
//
// It calls ctx.Next() to measure execution time of all following handlers.
//
// Behaviour on errors determined by abortOnErrors option. If it set to true request handling will be aborted with error.
func SpanFromContext(operationName string, abortOnErrors bool,
	advancedOpts ...stdot.StartSpanOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var opts []stdot.StartSpanOption

		if parentSpan, typeOk := GetSpan(ctx); typeOk {
			opts = append(opts, stdot.ChildOf(parentSpan.Context()))
		} else {
			if abortOnErrors {
				ctx.AbortWithError(http.StatusInternalServerError, ErrSpanNotFound)
			}
			return
		}
		opts = append(opts, advancedOpts...)

		span := tracer.StartSpan(operationName, opts...)
		ext.HTTPMethod.Set(span, ctx.Request.Method)
		ext.HTTPUrl.Set(span, ctx.Request.URL.String())
		ctx.Set(spanContextKey, span)
		defer span.Finish()

		ctx.Next()

		//ext.HTTPStatusCode.Set(span,ctx.Writer.Status())
		//ctx.Set(spanContextKey, span)
	}
}

// InjectToHeaders injects span meta-information to request headers.
//
// It may be useful when you want to trace chained request (client->service 1->service 2).
// In this case you have to save request headers (ctx.Request.Header) and pass it to next level request.
//
// Behaviour on errors determined by abortOnErrors option. If it set to true request handling will be aborted with error.
func InjectToHeaders(abortOnErrors bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var spanContext stdot.SpanContext
		if span, typeOk := GetSpan(ctx); typeOk {
			spanContext = span.Context()
		} else {
			if abortOnErrors {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrSpanNotFound)
			}
			return
		}

		tracer.Inject(spanContext, stdot.HTTPHeaders, stdot.HTTPHeadersCarrier(ctx.Request.Header))
	}
}

// GetSpan extracts span from context.
func GetSpan(ctx *gin.Context) (span stdot.Span, exists bool) {
	spanI, _ := ctx.Get(spanContextKey)
	span, ok := spanI.(stdot.Span)
	exists = span != nil && ok
	return
}

// MustGetSpan extracts span from context. It panics if span was not set.
func MustGetSpan(ctx *gin.Context) stdot.Span {
	return ctx.MustGet(spanContextKey).(stdot.Span)
}
