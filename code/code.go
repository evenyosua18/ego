package code

type Code struct {
	CustomCode      string `yaml:"code" json:"code"`
	ResponseMessage string `yaml:"message" json:"response_message"`
	ErrorMessage    string `yaml:"error" json:"error_message"`
	HttpCode        int    `yaml:"http_code" json:"http_code"`
	GrpcCode        int    `yaml:"grpc_code" json:"grpc_code"`
}

func (e Code) Error() string {
	return e.ErrorMessage
}

func (e Code) Message() string {
	return e.ResponseMessage
}

func (e Code) Code() string {
	return e.CustomCode
}

func (e Code) CodeHTTP() int {
	return e.HttpCode
}

func (e Code) CodeGRPC() int {
	return e.GrpcCode
}

func (e Code) SetErrorMessage(errMsg string) Code {
	e.ErrorMessage = errMsg
	return e
}

func (e Code) SetMessage(msg string) Code {
	e.ResponseMessage = msg
	return e
}

func (e Code) Err() error {
	return e
}
