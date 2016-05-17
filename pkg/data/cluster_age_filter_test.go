package data

import (
	"testing"
	"time"
)

func TestNewClusterAgeFilter(t *testing.T) {
	type testCase struct {
		chB time.Time
		chA time.Time
		crB time.Time
		crA time.Time
		err bool
	}

	testCases := []testCase{
		// checked in time test cases
		testCase{chB: timeNow(), chA: timeNow(), crB: timeFuture(), crA: timeNow(), err: true},
		testCase{chB: timeNow(), chA: timeFuture(), crB: timeFuture(), crA: timeNow(), err: true},
		testCase{chB: timePast(), chA: timeNow(), crB: timeFuture(), crA: timeNow(), err: true},
		testCase{chB: timeFuture(), chA: timeNow(), crB: timeFuture(), crA: timeNow(), err: false},
		// create time test cases
		testCase{chB: timeFuture(), chA: timeNow(), crB: timeNow(), crA: timeNow(), err: true},
		testCase{chB: timeFuture(), chA: timeNow(), crB: timeNow(), crA: timeFuture(), err: true},
		testCase{chB: timeFuture(), chA: timeNow(), crB: timePast(), crA: timeNow(), err: true},
		testCase{chB: timeFuture(), chA: timeNow(), crB: timeNow(), crA: timePast(), err: false},
	}

	for i, testCase := range testCases {
		filter, err := NewClusterAgeFilter(testCase.chB, testCase.chA, testCase.crB, testCase.crA)
		if testCase.err && err == nil {
			t.Errorf("expected error on iteration %d but got none", i)
			continue
		} else if !testCase.err && err != nil {
			t.Errorf("expected no error on iteration %d but got %s", i, err)
			continue
		}
		if filter == nil && err == nil {
			t.Errorf("got no error but resulting filter was nil on iteration %d", i)
			continue
		}
		if filter != nil && err != nil {
			t.Errorf("got an error but resulting filter was not nil on iteration %d", i)
			continue
		}
	}
}
