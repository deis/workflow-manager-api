package data

import (
	"database/sql/driver"
	"errors"
	"time"
)

const (
	// StdTimestampFmt is the standard format of timestamps used to store times in the database
	// and accept and send timestamps over the wire
	StdTimestampFmt = time.RFC3339
)

var (
	errInvalidType = errors.New("invalid type for current timestamp")
)

// Timestamp is a fmt.Stringer, sql.Scanner and driver.Valuer implementation which is able to encode and decode
// time.Time values into and out of a database. This implementation was inspired heavily from
// https://groups.google.com/forum/#!topic/golang-nuts/P6Wrm_uVvJ0
type Timestamp struct {
	Time time.Time
}

// Scan is the Scanner interface implementation
func (ts *Timestamp) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		ts.Time = v
		return nil
	case string:
		t, err := time.Parse(StdTimestampFmt, v)
		if err != nil {
			return err
		}
		ts.Time = t
		return nil
	case []byte:
		t, err := time.Parse(StdTimestampFmt, string(v))
		if err != nil {
			return err
		}
		ts.Time = t
		return nil
	default:
		return errInvalidType
	}
}

func newTimestampFromStr(tstr string) (Timestamp, error) {
	t, err := time.Parse(StdTimestampFmt, tstr)
	if err != nil {
		return Timestamp{}, err
	}
	return Timestamp{Time: t}, nil
}

// Value is the Valuer interface implementation
func (ts Timestamp) Value() (driver.Value, error) {
	str := ts.Time.Format(StdTimestampFmt)
	return str, nil
}

// String is the fmt.Stringer interface implementation
func (ts Timestamp) String() string {
	return ts.Time.Format(StdTimestampFmt)
}
