package data

import (
	"fmt"
)

type txErr struct {
	orig error
	err  error
	op   string
}

func (t txErr) Error() string {
	if t.orig != nil {
		return fmt.Sprintf("%s transaction error error (%s). original error (%s)", t.op, t.err, t.orig)
	}
	return fmt.Sprintf("%s transaction error (%s)", t.op, t.err)
}
