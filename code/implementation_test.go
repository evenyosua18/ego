package code

import (
	"fmt"
	"reflect"
	"testing"
)

func TestWrap(t *testing.T) {
	// add test custom code
	AddCustomCode(Code{
		CustomCode:      "test_code",
		ResponseMessage: "TEST",
		ErrorMessage:    "TEST",
		HttpCode:        0,
		GrpcCode:        200,
	})

	// args
	type args struct {
		err  error
		code string
	}

	tests := []struct {
		name string
		args args
		want Code
	}{
		{
			name: "code found",
			args: args{
				err:  fmt.Errorf("test error"),
				code: "test_code",
			},
			want: Code{
				CustomCode:      "test_code",
				ResponseMessage: "TEST",
				ErrorMessage:    "test error",
				HttpCode:        0,
				GrpcCode:        200,
			},
		},
		{
			name: "code not found",
			args: args{
				err:  fmt.Errorf("test error"),
				code: "test_unknown_code",
			},
			want: UnknownErrorCode,
		},
		{
			name: "error is nil",
			args: args{
				err:  nil,
				code: "test_code",
			},
			want: InvalidErrorCode,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Wrap(tt.args.err, tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	// add test custom code
	AddCustomCode(Code{
		CustomCode:      "test_code",
		ResponseMessage: "TEST",
		ErrorMessage:    "TEST",
		HttpCode:        0,
		GrpcCode:        200,
	})

	// args
	type args struct {
		code string
	}

	tests := []struct {
		name string
		args args
		want Code
	}{
		{
			name: "code exists",
			args: args{
				code: "test_code",
			},
			want: Code{
				CustomCode:      "test_code",
				ResponseMessage: "TEST",
				ErrorMessage:    "TEST",
				HttpCode:        0,
				GrpcCode:        200,
			},
		},
		{
			name: "code not exists",
			args: args{
				code: "test_unknown_code",
			},
			want: UnknownErrorCode,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want Code
	}{
		{
			name: "invalid error",
			args: args{
				err: fmt.Errorf("not formated error"),
			},
			want: InvalidErrorCode,
		},
		{
			name: "valid error",
			args: args{
				err: Code{
					CustomCode:      "test_code",
					ResponseMessage: "TEST",
					ErrorMessage:    "TEST",
					HttpCode:        0,
					GrpcCode:        200,
				},
			},
			want: Code{
				CustomCode:      "test_code",
				ResponseMessage: "TEST",
				ErrorMessage:    "TEST",
				HttpCode:        0,
				GrpcCode:        200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Extract(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extract() = %v, want %v", got, tt.want)
			}
		})
	}
}
