package data

import (
	"testing"
	"time"
)

func TestNewPersistentClustersFilter(t *testing.T) {
	type testCase struct {
		epoch     time.Time
		timestamp time.Time
		err       bool
	}

	testCases := []testCase{
		// test cases
		testCase{epoch: timeNow(), timestamp: timeNow(), err: true},
		testCase{epoch: timeFuture(), timestamp: timeNow(), err: true},
		testCase{epoch: timeNow(), timestamp: timePast(), err: true},
		testCase{epoch: timePast(), timestamp: timeNow(), err: false},
	}

	for i, testCase := range testCases {
		filter, err := NewPersistentClustersFilter(testCase.epoch, testCase.timestamp)
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
