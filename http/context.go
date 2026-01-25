package http

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/evenyosua18/ego/code"
	fiber "github.com/gofiber/fiber/v3"
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

func (f *fiberContext) GetRequestHeader(key string) string {
	return f.ctx.Get(key)
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

func (f *fiberContext) BindQuery(destination any) error {
	v := reflect.ValueOf(destination)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return code.Get(code.BadRequestError).SetErrorMessage("destination must be a non-nil pointer to struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return code.Get(code.BadRequestError).SetErrorMessage("destination must point to a struct")
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		// Skip unexported fields
		if !fieldVal.CanSet() {
			continue
		}

		tag := field.Tag.Get("query")
		if tag == "" {
			tag = field.Name
		}

		val := f.ctx.Query(tag)
		if val == "" {
			continue // skip empty
		}

		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				fieldVal.SetInt(i)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if u, err := strconv.ParseUint(val, 10, 64); err == nil {
				fieldVal.SetUint(u)
			}
		case reflect.Float32, reflect.Float64:
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				fieldVal.SetFloat(f)
			}
		case reflect.Bool:
			if b, err := strconv.ParseBool(val); err == nil {
				fieldVal.SetBool(b)
			}
		}
	}

	return nil
}

func (f *fiberContext) Cookie(name string, defaultValue ...string) string {
	return f.ctx.Cookies(name, defaultValue...)
}

func (f *fiberContext) SetCookie(name, value string, expiredAt time.Time) {
	f.ctx.Cookie(&fiber.Cookie{
		Expires:  expiredAt,
		Name:     name,
		Value:    value,
		Path:     "/",
		SameSite: fiber.CookieSameSiteLaxMode,
		Secure:   true,
		HTTPOnly: true,
	})
}

func (f *fiberContext) ResponseRedirect(to string) error {
	f.ctx.Location(to)
	f.ctx.Redirect()
	f.ctx.Status(http.StatusFound)

	return nil
}

func (f *fiberContext) Render(name string, data map[string]any) error {
	return f.ctx.Render(name, data)
}

func (f *fiberContext) FormValue(key string) string {
	return f.ctx.FormValue(key)
}
