package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	stdot "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	zipkinot "github.com/openzipkin/zipkin-go-opentracing"
)

var Tracer stdot.Tracer
var collector zipkinot.Collector

func InitTracer(url string, port int) error {
	var err error
	if url != "" {
		collector, err = zipkinot.NewHTTPCollector(url)
		if err != nil {
			return err
		}

		var (
			debug       = false
			hostPort    = fmt.Sprintf("localhost:%d", port)
			serviceName = "walm"
		)
		recorder := zipkinot.NewRecorder(collector, debug, hostPort, serviceName)
		Tracer, err = zipkinot.NewTracer(recorder)
		if err != nil {
			return err
		}
	} else {
		return errors.New("zipin url is none")
	}
	return nil
}

/*
func EnableTrace() gin.HandlerFunc {
	return trace.SpanFromHeaders(Tracer, "Walm", stdot.ChildOf, false)
}
*/

func EndTrace() {
	defer collector.Close()
}

const spanContextKey = "span"

// Errors which may occur at operation time.
var (
	ErrSpanNotFound = errors.New("span was not found in context")
)

// NewSpan returns gin.HandlerFunc (middleware) that starts a new span and injects it to request context.
//
// It calls ctx.Next() to measure execution time of all following handlers.
func NewSpan(tracer stdot.Tracer, operationName string, opts ...stdot.StartSpanOption) gin.HandlerFunc {
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
//
// See opentracing.SpanReferenceType
type ParentSpanReferenceFunc func(stdot.SpanContext) stdot.StartSpanOption

// SpanFromHeaders returns gin.HandlerFunc (middleware) that extracts parent span data from HTTP headers and
// starts a new span referenced to parent with ParentSpanReferenceFunc.
//
// It calls ctx.Next() to measure execution time of all following handlers.
//
// Behaviour on errors determined by abortOnErrors option. If it set to true request handling will be aborted with error.
func SpanFromHeaders(tracer stdot.Tracer, operationName string, psr ParentSpanReferenceFunc,
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
	}
}

// SpanFromContext returns gin.HandlerFunc (middleware) that extracts parent span from request context
// and starts a new span as child of parent span.
//
// It calls ctx.Next() to measure execution time of all following handlers.
//
// Behaviour on errors determined by abortOnErrors option. If it set to true request handling will be aborted with error.
func SpanFromContext(tracer stdot.Tracer, operationName string, abortOnErrors bool,
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
	}
}

// InjectToHeaders injects span meta-information to request headers.
//
// It may be useful when you want to trace chained request (client->service 1->service 2).
// In this case you have to save request headers (ctx.Request.Header) and pass it to next level request.
//
// Behaviour on errors determined by abortOnErrors option. If it set to true request handling will be aborted with error.
func InjectToHeaders(tracer stdot.Tracer, abortOnErrors bool) gin.HandlerFunc {
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
