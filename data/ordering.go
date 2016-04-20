package data

import (
	"fmt"
)

type orderBy struct {
	expression string
	sorting    string
}

func newOrderBy(expr, sort string) *orderBy {
	return &orderBy{expression: expr, sorting: sort}
}

func (o orderBy) String() string {
	if o.sorting == "" {
		return fmt.Sprintf("ORDER BY %s", o.expression)
	}
	return fmt.Sprintf("ORDER BY %s %s", o.expression, o.sorting)
}
