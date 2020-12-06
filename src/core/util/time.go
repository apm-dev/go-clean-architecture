package util

import "time"

var Time TimeInterface = &timeImpl{}

type TimeInterface interface {
	FromUnix(int64) time.Time
	Now() time.Time
	NowUnix() int64
}

type timeImpl struct {
}

func (*timeImpl) FromUnix(i int64) time.Time {
	return time.Unix(i, 0)
}

func (*timeImpl) Now() time.Time {
	return time.Now().UTC()
}

func (t *timeImpl) NowUnix() int64 {
	return t.Now().Unix()
}
