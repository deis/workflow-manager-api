package data

import (
	"testing"
	"time"
)

// since various times have been decoded from different encodings, they have different bits of
// information. for example, a time decoded from time.Kitchen (3:04PM) has much less information
// than time.RFC3339Nano (2006-01-02T15:04:05.999999999Z07:00), so t1.Equal(t2) will be false.
// this func tries to determine whether they are "close enough"
func fuzzyTimeEqual(t1 time.Time, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

type timestampScanTestCase struct {
	val          interface{}
	expectedTime time.Time
	expectedErr  bool
}

func TestTimestampScan(t *testing.T) {
	now := time.Now()
	testCases := []timestampScanTestCase{
		timestampScanTestCase{val: now.Format(stdTimestampFmt), expectedTime: now, expectedErr: false},
		timestampScanTestCase{val: []byte(now.Format(stdTimestampFmt)), expectedTime: now, expectedErr: false},
		timestampScanTestCase{val: now.Format(time.ANSIC), expectedTime: now, expectedErr: true},
		timestampScanTestCase{val: now.Format(time.UnixDate), expectedTime: now, expectedErr: true},
		timestampScanTestCase{val: now.Format(time.RubyDate), expectedTime: now, expectedErr: true},
		timestampScanTestCase{val: now.Format(time.RFC822), expectedTime: now, expectedErr: true},
		timestampScanTestCase{val: now.Format(time.RFC822Z), expectedTime: now, expectedErr: true},
		timestampScanTestCase{val: now.Format(time.RFC850), expectedTime: now, expectedErr: true},
		timestampScanTestCase{val: now.Format(time.RFC1123), expectedTime: now, expectedErr: true},
		timestampScanTestCase{val: now.Format(time.RFC1123Z), expectedTime: now, expectedErr: true},
		timestampScanTestCase{val: now.Format(time.RFC3339), expectedTime: now, expectedErr: false},
		// RFC3339Nano is a superset of RFC3339, so the time package can parse it
		timestampScanTestCase{val: now.Format(time.RFC3339Nano), expectedTime: now, expectedErr: false},
		timestampScanTestCase{val: now.Format(time.Kitchen), expectedTime: now, expectedErr: true},
		timestampScanTestCase{val: true, expectedTime: now, expectedErr: true},
	}
	for i, testCase := range testCases {
		ts := new(Timestamp)
		err := ts.Scan(testCase.val)
		if testCase.expectedErr && err == nil {
			t.Errorf("test case %d expected err", i+1)
			continue
		}
		if !testCase.expectedErr && err != nil {
			t.Errorf("test case %d didn't expect err", i+1)
			continue
		}
		if testCase.expectedErr && err != nil {
			continue
		}
		if ts.Time == nil {
			t.Errorf("test case %d expected non-nil time", i+1)
			continue
		}
		if !fuzzyTimeEqual(testCase.expectedTime, *ts.Time) {
			t.Errorf("test case %d expected time %s doesn't match actual %s", i+1, testCase.expectedTime, *ts.Time)
			continue
		}
	}
}

func TestTimestampNow(t *testing.T) {
	n := now()
	if _, err := time.Parse(stdTimestampFmt, n); err != nil {
		t.Fatal(err)
	}
}

func TestTimestampString(t *testing.T) {
	now := time.Now()
	ts := new(Timestamp)
	ts.Time = &now
	if _, err := time.Parse(stdTimestampFmt, ts.String()); err != nil {
		t.Fatal(err)
	}
}

func TestTimestampValue(t *testing.T) {
	now := time.Now()
	ts := new(Timestamp)
	ts.Time = &now
	val, err := ts.Value()
	if err != nil {
		t.Fatal(err)
	}
	str, ok := val.(string)
	if !ok {
		t.Fatalf("returned value was not a string")
	}
	if _, err := time.Parse(stdTimestampFmt, str); err != nil {
		t.Fatalf("returned value was not in expected format %s (%s)", stdTimestampFmt, err)
	}
}
