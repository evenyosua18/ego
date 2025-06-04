package http

import (
	"context"
	"github.com/evenyosua18/ego/code"
	"github.com/gofiber/fiber/v3"
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
