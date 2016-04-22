package data

import (
	"fmt"
	"time"
)

var (
	zeroTime   time.Time
	nowTime    = time.Now()
	futureTime = nowTime.Add(1 * time.Hour)
	pastTime   = nowTime.Add(-1 * time.Hour)
)

func timeFuture() time.Time {
	return futureTime
}

func timePast() time.Time {
	return pastTime
}

func timeNow() time.Time {
	return nowTime
}

type errNoPossibleBetweenTime struct {
	lt time.Time
	gt time.Time
}

func (e errNoPossibleBetweenTime) Error() string {
	return fmt.Sprintf("no possible time between %s and %s", e.lt, e.gt)
}

// returns a time such that t1 < retval < t2. if that inequality is
// impossible, returns the zero time and errNoPossibleBetweenTime
func betweenTimes(t1, t2 time.Time) (time.Time, error) {
	if t1.Before(t2) {
		diff := t2.Sub(t1)
		return t1.Add(diff / 2), nil
	}
	return zeroTime, errNoPossibleBetweenTime{lt: t1, gt: t2}
}

// returns a time that is less than both t1 and t2
func lessThanTimes(t1, t2 time.Time) time.Time {
	if t1.Before(t2) {
		return t1.Add(-1 * time.Hour)
	}
	return t2.Add(-1 * time.Hour)
}

// returns a time that is greater than both t1 and t2
func greaterThanTimes(t1, t2 time.Time) time.Time {
	if t1.Before(t2) {
		return t2.Add(1 * time.Hour)
	}
	return t1.Add(1 * time.Hour)
}
