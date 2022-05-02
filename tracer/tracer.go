package tracer

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"runtime"
)

const (
	requestName  = "request"
	responseName = "response"
	errorName    = "error"
)

type Tracer struct {
	ctx context.Context
}

func RootTracer(name string, reqCtx context.Context) (ctx context.Context, span trace.Span) {
	//get caller details
	caller, function := getCaller()

	//get tracer
	tracer := otel.Tracer(name)

	//create root span
	ctx, span = tracer.Start(reqCtx, function)

	//set attribute for caller description
	span.SetAttributes(attribute.Key("caller").String(caller))

	return
}

func ChildTracer(reqCtx context.Context) (ctx context.Context, span trace.Span) {
	//get caller details
	caller, function := getCaller()

	//get tracer
	tracer := otel.Tracer(function)

	//create new span
	ctx, span = tracer.Start(reqCtx, function)

	//set default attribute
	span.SetAttributes(attribute.Key("caller").String(caller)) //caller detail

	return
}

func LogRequest(span trace.Span, name string, request interface{}) {
	reqByte, _ := json.Marshal(request)
	span.AddEvent(name, trace.WithAttributes(attribute.Key(requestName).String(string(reqByte))))
}

func LogResponse(span trace.Span, name string, response interface{}) {
	resByte, _ := json.Marshal(response)
	span.AddEvent(name, trace.WithAttributes(attribute.Key(responseName).String(string(resByte))))
}

func LogError(span trace.Span, name string, err error) {
	span.AddEvent(name, trace.WithAttributes(attribute.Key(errorName).String(err.Error())))
	span.SetStatus(1, err.Error())
}

func getCaller() (description, function string) {
	// Get caller function name, file and line
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)

	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	description = fmt.Sprintf("%s - %s#%d", frame.Function, frame.File, frame.Line)
	function = frame.Function

	return
}
