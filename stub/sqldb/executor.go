package sqldb

type StubExecutor struct {
	ExecError         error
	ResultValue       int64
	ResultError       error
	QueryRowValues    []any
	QueryRowErr       error
	QueryValues       []any
	QueryDestination  any
	QueryErr          error
	QueryCloseErr     error
	QueryMapResultErr error
	QueryScanAllErr   error
	QueryScanOneErr   error
	QueryRowIndex     int
}
