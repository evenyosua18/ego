package code

import "errors"

// Wrap for combine error with custom code
func Wrap(err error, code string) Code {
	customCode, ok := customCodes[code]
	if !ok {
		return *customCodes[defaultUnknownCode]
	}

	if customCode == nil {
		return *customCodes[defaultUnknownCode]
	} else if err == nil {
		return *customCodes[defaultInvalidError]
	}

	return Code{
		CustomCode:      customCode.CustomCode,
		ResponseMessage: customCode.ResponseMessage,
		ErrorMessage:    err.Error(),
		HttpCode:        customCode.HttpCode,
		GrpcCode:        customCode.GrpcCode,
	}
}

// Get custom code
func Get(code string) Code {
	customCode, ok := customCodes[code]
	if !ok {
		return *customCodes[defaultUnknownCode]
	}

	if customCode == nil {
		return UnknownErrorCode
	}

	return *customCode
}

// Extract error to struct
func Extract(err error) Code {
	var customCode Code
	if errors.As(err, &customCode) {
		return customCode
	}

	return InvalidErrorCode
}
