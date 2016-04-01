package data

import (
	"database/sql/driver"
	"errors"
	"time"
)

const (
	stdTimestampFmt = time.RFC3339
)

var (
	errInvalidType = errors.New("invalid type for current timestamp")
)

// Timestamp is a fmt.Stringer, sql.Scanner and driver.Valuer implementation which is able to encode and decode
// time.Time values into and out of a database. This implementation was inspired heavily from
// https://groups.google.com/forum/#!topic/golang-nuts/P6Wrm_uVvJ0
type Timestamp struct {
	Time *time.Time
}

// Scan is the Scanner interface implementation
func (ts *Timestamp) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		ts.Time = &v
		return nil
	case string:
		t, err := time.Parse(stdTimestampFmt, v)
		if err != nil {
			return err
		}
		ts.Time = &t
		return nil
	case []byte:
		t, err := time.Parse(stdTimestampFmt, string(v))
		if err != nil {
			return err
		}
		ts.Time = &t
		return nil
	default:
		return errInvalidType
	}
}

// Value is the Valuer interface implementation
func (ts *Timestamp) Value() (driver.Value, error) {
	str := ts.Time.Format(stdTimestampFmt)
	return str, nil
}

// String is the fmt.Stringer interface implementation
func (ts *Timestamp) String() string {
	return ts.Time.Format(stdTimestampFmt)
}

func now() string {
	return time.Now().Format(stdTimestampFmt)
}
