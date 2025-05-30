package code

import (
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"testing"
)

func TestLoadCodes(t *testing.T) {
	tests := []struct {
		name        string
		data        map[string][]map[string]any
		expectErr   bool
		expectCodes map[string]*Code
	}{
		{
			name: "load from temporary yaml",
			data: map[string][]map[string]any{
				"codes": {
					{
						"code":      "test_error",
						"message":   "TEST",
						"error":     "test error",
						"http_code": 200,
						"grpc_code": 0,
					},
				},
			},
			expectErr: false,
			expectCodes: map[string]*Code{
				UnknownErrorCode.CustomCode: &UnknownErrorCode,
				InvalidErrorCode.CustomCode: &InvalidErrorCode,
				"test_error": {
					CustomCode:      "test_error",
					ResponseMessage: "TEST",
					ErrorMessage:    "test error",
					HttpCode:        200,
					GrpcCode:        0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare temporary yaml file
			tmpfile, err := os.CreateTemp("", "*.yaml")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			yamlData, err := yaml.Marshal(tt.data)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := tmpfile.Write(yamlData); err != nil {
				t.Fatalf("failed to write YAML to file: %v", err)
			}
			tmpfile.Close()

			if err != nil {
				t.Errorf("unexpected error when preparing yaml: %v", err)
			}

			// load codes from temporary yaml
			err = LoadCodes(tmpfile.Name())

			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
				return
			} else if err != nil {
				t.Errorf("unexpected error when load yaml codes: %v", err)
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// at least all codes from yaml loaded successfully
			gotCodes := List()
			for key, code := range tt.expectCodes {
				if !reflect.DeepEqual(gotCodes[key], code) {
					t.Errorf("expected %v but got %v", code, gotCodes[key])
				}
			}
		})
	}
}

func TestAddCustomCode(t *testing.T) {
	type args struct {
		codes []Code
	}

	tests := []struct {
		name string
		args args
		want map[string]*Code
	}{
		{
			name: "test init without add new custom code",
			args: args{},
			want: map[string]*Code{
				UnknownErrorCode.CustomCode: &UnknownErrorCode,
				InvalidErrorCode.CustomCode: &InvalidErrorCode,
			},
		},
		{
			name: "add new custom code",
			args: args{
				codes: []Code{
					{
						CustomCode:      "test_code",
						ResponseMessage: "TEST",
						ErrorMessage:    "TEST",
						HttpCode:        0,
						GrpcCode:        200,
					},
				},
			},
			want: map[string]*Code{
				"test_code": {
					CustomCode:      "test_code",
					ResponseMessage: "TEST",
					ErrorMessage:    "TEST",
					HttpCode:        0,
					GrpcCode:        200,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantList := List()

			// add new code
			if len(tt.args.codes) != 0 {
				for _, code := range tt.args.codes {
					AddCustomCode(code)
				}
			}

			// set the want list + default list
			for k, v := range wantList {
				wantList[k] = v
			}

			// check length of map codes
			if len(wantList) != len(List()) {
				t.Errorf("invalid length of codes %d, want %d", len(List()), len(tt.want))
			}

			// check list of codes
			if !reflect.DeepEqual(wantList, List()) {
				t.Errorf("AddCustomCode() = %v, want %v", List(), wantList)
			}
		})
	}
}

func TestCustomCodesInit(t *testing.T) {
	// check map is not empty
	if len(customCodes) == 0 {
		t.Fatal("custom codes map is empty after init")
	}

	// get sample from one of the default error code
	code, ok := customCodes[InternalError]
	if !ok {
		t.Fatalf("custom codes missing key %q", InternalError)
	}

	if code.HttpCode != 500 {
		t.Errorf("expected http code 500 for %q, got %d", InternalError, code.HttpCode)
	}
}
