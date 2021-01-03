package lesson

type Op string

const (
	OpOpEquals          Op = "="
	OpOpNotEquals       Op = "!="
	OpLessThan          Op = "<"
	OpLessThanEquals    Op = "<="
	OpGreaterThan       Op = ">"
	OpGreaterThanEquals Op = ">="
	OpContains          Op = "contains"
	OpStartsWith        Op = "starts-with"
)

var OpOps = [...]Op{OpOpEquals, OpOpNotEquals, OpLessThan,
	OpLessThanEquals, OpGreaterThan, OpGreaterThanEquals, OpContains,
	OpStartsWith}