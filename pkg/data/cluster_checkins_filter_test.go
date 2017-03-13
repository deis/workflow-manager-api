package data

import (
	"testing"
	"time"
)

func TestNewClusterCheckinsFilter(t *testing.T) {
	type testCase struct {
		crA time.Time
		crB time.Time
		err bool
	}

	testCases := []testCase{
		// test cases
		testCase{crA: timeNow(), crB: timeNow(), err: true},
		testCase{crA: timeFuture(), crB: timeNow(), err: true},
		testCase{crA: timeNow(), crB: timePast(), err: true},
		testCase{crA: timePast(), crB: timeNow(), err: false},
	}

	for i, testCase := range testCases {
		filter, err := NewClusterCheckinsFilter(testCase.crA, testCase.crB)
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
