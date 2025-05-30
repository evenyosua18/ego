package code

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var (
	customCodes         map[string]*Code
	defaultUnknownCode  = "unknown_code"
	defaultInvalidError = "invalid_error"

	// default error code

	UnknownErrorCode = Code{
		CustomCode:      defaultUnknownCode,
		ResponseMessage: "need to register your custom code",
		ErrorMessage:    "unknown code",
		HttpCode:        500,
		GrpcCode:        2,
	}

	InvalidErrorCode = Code{
		CustomCode:      defaultInvalidError,
		ResponseMessage: "error is invalid or got an empty error",
		ErrorMessage:    "invalid error",
		HttpCode:        500,
		GrpcCode:        13,
	}

	// universal error code

	InternalError     = "internal_error"
	NotFoundError     = "not_found_error"
	DatabaseError     = "database_error"
	BadRequestError   = "bad_request"
	UnauthorizedError = "unauthorized"
)

func init() {
	customCodes = map[string]*Code{
		defaultUnknownCode:  &UnknownErrorCode,
		defaultInvalidError: &InvalidErrorCode,

		// universal error code

		InternalError: {
			CustomCode:      InternalError,
			ResponseMessage: "something when wrong, will be fixed as soon as possible",
			ErrorMessage:    "internal server error",
			HttpCode:        500,
			GrpcCode:        13,
		},
		NotFoundError: {
			CustomCode:      NotFoundError,
			ResponseMessage: "data not found",
			ErrorMessage:    "not found",
			HttpCode:        404,
			GrpcCode:        5,
		},
		DatabaseError: {
			CustomCode:      DatabaseError,
			ResponseMessage: "something went wrong, will be fixed as soon as possible",
			ErrorMessage:    "not found",
			HttpCode:        500,
			GrpcCode:        13,
		},
		BadRequestError: {
			CustomCode:      BadRequestError,
			ResponseMessage: "please check the request again",
			ErrorMessage:    "bad request",
			HttpCode:        400,
			GrpcCode:        9,
		},
		UnauthorizedError: {
			CustomCode:      UnauthorizedError,
			ResponseMessage: "permission not found",
			ErrorMessage:    "unauthorized",
			HttpCode:        400,
			GrpcCode:        16,
		},
	}
}

func LoadCode(path string) {
	//read file
	f, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	e := struct {
		Codes []Code `yaml:"codes"`
	}{}

	//unmarshal yaml file
	if err = yaml.Unmarshal(f, &e); err != nil {
		panic(err)
	}

	//save to map
	for _, code := range e.Codes {
		AddCustomCode(code)
	}

	// TODO replace to internal log
	log.Printf("success register %d codes", len(e.Codes))
}

func AddCustomCode(code Code) {
	customCodes[code.CustomCode] = &code
}

func List() map[string]*Code {
	return customCodes
}
