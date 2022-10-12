package query_utils

import "github.com/elizabeth-dev/FACEIT_Test/internal/pkg/utils/query_utils/operators"

type Filter struct {
	Field    string
	Operator operators.Comparison
	Value    interface{}
}

type Sort struct {
	Field     string
	Direction operators.Sort
}

type Pagination struct {
	Limit  int64
	Offset int64
}
