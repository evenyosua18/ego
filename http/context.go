package http

import (
	"context"
	"encoding/json"
	"github.com/evenyosua18/ego/code"
	"github.com/gofiber/fiber/v3"
	"strings"
)

type fiberContext struct {
	ctx fiber.Ctx
}

func (f *fiberContext) Param(key string) string {
	return f.ctx.Params(key)
}

func (f *fiberContext) Query(key string) string {
	return f.ctx.Query(key)
}

func (f *fiberContext) Body() []byte {
	return f.ctx.Body()
}

func (f *fiberContext) RequestBody(res any) error {
	return json.Unmarshal(f.ctx.Body(), &res)
}

func (f *fiberContext) Send(status int, body []byte) error {
	return f.ctx.Status(status).Send(body)
}

func (f *fiberContext) JSON(status int, data any) error {
	return f.ctx.Status(status).JSON(data)
}

func (f *fiberContext) Context() context.Context {
	return f.ctx.Context()
}

func (f *fiberContext) ResponseError(err error) error {
	c := code.Extract(err)

	return f.ctx.Status(c.HttpCode).JSON(c)
}

func (f *fiberContext) ResponseSuccess(data any) error {
	return f.ctx.Status(200).JSON(data)
}

func (f *fiberContext) HttpData() (data map[string]any, operationName string) {
	return map[string]any{
		"method":     f.ctx.Method(),
		"path":       f.ctx.Path(),
		"query":      f.ctx.Queries(),
		"ip_address": f.ctx.IP(),
		"body":       truncate(f.ctx.Get("Content-Type"), f.ctx.Body()),
	}, strings.ToUpper(f.ctx.Method()) + " " + f.ctx.Path()
}

const (
	maxLogLength = 2048

	truncatedText     = "...(truncated)"
	unableMarshalText = "<<unable to marshal object>>"
	omittedFileText   = "<<file omitted>>"
)

func truncate(contentType string, obj any) string {
	if !strings.HasPrefix(contentType, "application/json") &&
		!strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		return omittedFileText
	}

	switch v := obj.(type) {
	case string:
		if len(v) > maxLogLength {
			return v[:maxLogLength] + truncatedText
		}
		return v
	case []byte:
		if len(v) > maxLogLength {
			return string(v[:maxLogLength]) + truncatedText
		}
		return string(v)
	}

	// Marshal to JSON for structured types
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return unableMarshalText
	}

	if len(jsonBytes) > maxLogLength {
		return string(jsonBytes[:maxLogLength]) + truncatedText
	}

	return string(jsonBytes)
}
