package code

import (
	"errors"
	"testing"
)

func TestCode(t *testing.T) {
	type args struct {
		CustomCode      string
		ResponseMessage string
		ErrorMessage    string
		HttpCode        int
		GrpcCode        int
	}

	type field struct {
		err             error
		customCode      string
		responseMessage string
		errorMessage    string
		httpCode        int
		grpcCode        int
	}

	tests := []struct {
		name string
		args args
		want field
	}{
		{
			name: "test all code getter",
			args: args{
				CustomCode:      "test",
				ResponseMessage: "TEST",
				ErrorMessage:    "TEST",
				HttpCode:        200,
				GrpcCode:        200,
			},
			want: field{
				customCode:      "test",
				responseMessage: "TEST",
				errorMessage:    "TEST",
				httpCode:        200,
				grpcCode:        200,
				err: Code{
					CustomCode:      "test",
					ResponseMessage: "TEST",
					ErrorMessage:    "TEST",
					HttpCode:        200,
					GrpcCode:        200,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create code
			code := &Code{
				CustomCode:      tt.args.CustomCode,
				ResponseMessage: tt.args.ResponseMessage,
				ErrorMessage:    tt.args.ErrorMessage,
				HttpCode:        tt.args.HttpCode,
				GrpcCode:        tt.args.GrpcCode,
			}

			if code.CodeHTTP() != tt.want.httpCode {
				t.Errorf("http code error,want %d, got %d", tt.want.httpCode, code.CodeHTTP())
			}

			if code.Code() != tt.want.customCode {
				t.Errorf("custom code error, want %s, got %s", tt.want.customCode, code.Code())
			}

			if code.CodeGRPC() != tt.want.grpcCode {
				t.Errorf("grpc code error, want %d, got %d", tt.want.grpcCode, code.CodeGRPC())
			}

			if code.Error() != tt.want.errorMessage {
				t.Errorf("error message error, want %s, got %s", tt.want.errorMessage, code.Error())
			}

			if !errors.Is(tt.want.err, code.Err()) {
				t.Errorf("error result error, want %s, got %s", tt.want.err, code.Err())
			}
		})
	}
}
