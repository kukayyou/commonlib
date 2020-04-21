package uchystrixhttp

import (
	"context"
	"gnetis.com/golang/core/golib/uctracer"
	"regexp"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

//SetContext 设置上下文
func (Req *UCRequest) SetContext(ctx context.Context) *UCRequest {
	if uctracer.TracerEnable == true {
		Req.tracerContext = ctx
	}
	return Req
}

func (Req *UCRequest) startSpan() {
	if uctracer.TracerEnable == false {
		return
	}
	if Req.tracerContext == nil {
		Req.tracerContext = uctracer.Context()
	}
	if Req.tracerContext != nil {
		Req.tracerSpan, _ = opentracing.StartSpanFromContext(Req.tracerContext, "http request")
		ext.SpanKind.Set(Req.tracerSpan, "client")
		ext.HTTPMethod.Set(Req.tracerSpan, Req.method)
		ext.HTTPUrl.Set(Req.tracerSpan, Req.url)
		if requestID := uctracer.RequestID(); requestID != "" {
			Req.reqID = requestID
		}
		tracerInfo := map[string]string{}
		Req.tracerSpan.Tracer().Inject(Req.tracerSpan.Context(), opentracing.TextMap, opentracing.TextMapCarrier(tracerInfo))
		Req.SetHeader(tracerInfo)
	}
}

func (Req *UCRequest) finish() {
	if Req.tracerSpan != nil && Req.tracerContext != nil {
		reg := regexp.MustCompile(`(\"code\"|\"errCode\"|\"statusCode\")\:(\d{0,10})`)
		for _, r := range reg.FindAllSubmatch(Req.resBody, -1) {
			if len(r) == 3 {
				Req.tracerSpan.SetTag("http.response.code", string(r[2]))
				break
			}
		}
		if Req.errMsg != nil {
			Req.tracerSpan.SetTag("error", Req.errMsg)
		}
		ext.HTTPStatusCode.Set(Req.tracerSpan, uint16(Req.statusCode))
		Req.tracerSpan.Finish()
	}
}
