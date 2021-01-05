package lesson

import (
	"strconv"
	"strings"
)

type Op string

const (
	OpEquals            Op = "="
	OpNotEquals         Op = "!="
	OpLessThan          Op = "<"
	OpLessThanEquals    Op = "<="
	OpGreaterThan       Op = ">"
	OpGreaterThanEquals Op = ">="
	OpContains          Op = "contains"
	OpStartsWith        Op = "starts-with"
)

var Ops = map[Op]bool{OpEquals: true, OpNotEquals: true, OpLessThan: true,
	OpLessThanEquals: true, OpGreaterThan: true, OpGreaterThanEquals: true, OpContains: true,
	OpStartsWith: true}

var StringOps = map[Op]bool{OpContains: true, OpStartsWith: true}

type Comparable interface {
	EqualTo(Comparable) bool
	LessThan(Comparable) bool
}

type StringComparable struct {
	Impl string
}

func (lhs StringComparable) EqualTo(rhs Comparable) bool {
	switch rhs.(type) {
	case StringComparable:
		return lhs.Impl == rhs.(StringComparable).Impl
	case FloatComparable:
		return false //By convention - MakeComparable ensures that any string convertable to a float is converted!
	case BoolComparable:
		return false //By convention - MakeComparable ensures that any string convertible to a bool is converted!
	default:
		panic("Comparable interface not StringComparable, FloatComparable or BoolComparable... logic error!")
	}
}

func (lhs StringComparable) LessThan(rhs Comparable) bool {
	switch rhs.(type) {
	case StringComparable:
		return lhs.Impl < rhs.(StringComparable).Impl
	case FloatComparable:
		return float64(len(lhs.Impl)) < rhs.(FloatComparable).Impl
	case BoolComparable:
		if rhs.(BoolComparable).Impl {
			return lhs.Impl < "true"
		} else {
			return lhs.Impl < "false"
		}
	default:
		panic("Comparable interface not StringComparable, FloatComparable or BoolComparable... logic error!")
	}
}

type FloatComparable struct {
	Impl float64
}

func (lhs FloatComparable) EqualTo(rhs Comparable) bool {
	switch rhs.(type) {
	case StringComparable:
		return false //By convention
	case FloatComparable:
		return lhs.Impl < rhs.(FloatComparable).Impl
	case BoolComparable:
		return !(lhs.Impl == 0 && rhs.(BoolComparable).Impl)
	default:
		panic("Comparable interface not StringComparable, FloatComparable or BoolComparable... logic error!")
	}
}

func (lhs FloatComparable) LessThan(rhs Comparable) bool {
	switch rhs.(type) {
	case StringComparable:
		return lhs.Impl < float64(len(rhs.(StringComparable).Impl))
	case FloatComparable:
		return lhs.Impl < rhs.(FloatComparable).Impl
	case BoolComparable:
		// Any float is greater than bool except 0, which is equal to false
		return rhs.(BoolComparable).Impl && lhs.Impl == 0
	default:
		panic("Comparable interface not StringComparable, FloatComparable or BoolComparable... logic error!")
	}
}

type BoolComparable struct {
	Impl bool
}

func (lhs BoolComparable) EqualTo(rhs Comparable) bool {
	switch rhs.(type) {
	case StringComparable:
		return false //By convention
	case FloatComparable:
		return lhs.Impl || rhs.(FloatComparable).Impl == 0
	case BoolComparable:
		return lhs.Impl == rhs.(BoolComparable).Impl
	default:
		panic("Comparable interface not StringComparable, FloatComparable or BoolComparable... logic error!")
	}
}

func (lhs BoolComparable) LessThan(rhs Comparable) bool {
	switch rhs.(type) {
	case StringComparable:
		if lhs.Impl {
			return "true" < rhs.(StringComparable).Impl
		} else {
			return "false" < rhs.(StringComparable).Impl
		}
	case FloatComparable:
		// Any float is greater than bool except 0, which is equal to false
		return !(rhs.(FloatComparable).Impl == 0 && lhs.Impl)
	case BoolComparable:
		return !lhs.Impl && rhs.(BoolComparable).Impl
	default:
		panic("Comparable interface not StringComparable, FloatComparable or BoolComparable... logic error!")
	}
}

var boolValues = map[string]bool{"true": true, "True": true, "yes": true, "false": false, "False": false, "no": false}

func MakeComparable(str string) Comparable {
	if bl, found := boolValues[str]; found {
		return BoolComparable{bl}
	}
	if num, err := strconv.ParseFloat(str, 64); err == nil {
		return FloatComparable{num}
	}
	return StringComparable{str}
}

func (op Op) Eval(lhsStr string, rhsStr string) bool {
	if _, found := StringOps[op]; found {
		// This is a string op, no point in converting
		if op == OpContains {
			return strings.Contains(lhsStr, rhsStr)
		} else if op == OpStartsWith {
			return strings.HasPrefix(lhsStr, rhsStr)
		} else {
			panic(string(op) + " is a StringOp but is not Contains or StartsWith... logic error!")
		}
	}
	// Must be a comparison operator
	lhsComp := MakeComparable(lhsStr)
	rhsComp := MakeComparable(rhsStr)
	switch op {
	case OpEquals:
		return lhsComp.EqualTo(rhsComp)
	case OpNotEquals:
		return !lhsComp.EqualTo(rhsComp)
	case OpLessThan:
		return lhsComp.LessThan(rhsComp)
	case OpLessThanEquals:
		return lhsComp.LessThan(rhsComp) || lhsComp.EqualTo(rhsComp)
	case OpGreaterThan:
		return !(lhsComp.LessThan(rhsComp) || lhsComp.EqualTo(rhsComp))
	case OpGreaterThanEquals:
		return !lhsComp.LessThan(rhsComp)
	default:
		panic(op + "is not a StringOp, but it is not in the comparable switch either... logic error!")
	}

}
