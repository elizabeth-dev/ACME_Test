package operators

type Comparison int

const (
	EQUALS Comparison = iota
	NOT_EQUALS
	GREATER_THAN
	GREATER_THAN_EQ
	LESS_THAN
	LESS_THAN_EQ
)
