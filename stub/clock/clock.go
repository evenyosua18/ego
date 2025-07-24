package clock

import "time"

var (
	Now = func() time.Time {
		return time.Now().Local()
	}

	StubValue = time.Date(2024, time.December, 1, 1, 0, 0, 0, time.Local)
)

func Stub() time.Time {
	return StubValue
}
