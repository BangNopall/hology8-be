package time

import "time"

type TimeInterface interface {
	Now() time.Time
	Add(duration time.Duration) time.Time
}

type TimeStruct struct{}

var Time = getTime()

func getTime() TimeInterface {
	return &TimeStruct{}
}

func (t *TimeStruct) Now() time.Time {
	return time.Now()
}

func (t *TimeStruct) Add(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}
