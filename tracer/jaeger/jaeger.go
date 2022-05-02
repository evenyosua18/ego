package jaeger

import (
	"errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.6.1"
)

type SetupJaeger struct {
	attributes  []attribute.KeyValue
	ServiceName string
	Endpoints   string
}

func New(model SetupJaeger) (tp *sdkTrace.TracerProvider, err error) {
	//setup default value
	if model.Endpoints == "" {
		err = errors.New("jaeger url can't be nil")
		return
	}

	if model.ServiceName == "" {
		model.ServiceName = "example-service"
	}

	//get exporter
	var exporter *jaeger.Exporter
	if exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(model.Endpoints))); err != nil {
		return
	}

	//setup attributes
	model.attributes = append(model.attributes, semconv.ServiceNameKey.String(model.ServiceName))

	//get tracer provider from jaeger
	tp = sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(exporter),
		sdkTrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			model.attributes...,
		)),
	)

	return
}

func (j *SetupJaeger) SetAttribute(key, value string) {
	j.attributes = append(j.attributes, attribute.Key(key).String(value))
}
