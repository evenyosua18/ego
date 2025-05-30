package code

type ICode interface {
	Error() string //builtin, return error message

	Message() string // return response message
	Code() string
	CodeGRPC() int
	CodeHTTP() int
}
