package lesson

import (
	"gorm.io/gorm"
	"strings"
)

type QuestionFilterLiteral struct {
	Field string //TODO need to json this.
	Op    Op
	Value string
}

func (qfl QuestionFilterLiteral) ToString() string {
	return strings.Join([]string{qfl.Field, string(qfl.Op), qfl.Value}, "")
}

//TODO this basically needs to be a list of conditions
//TODO check how gorm structures things
type QuestionFilter struct {
	QuestionBanks []string
	//TODO TODO this just needs to get translated into a nice big where query
	CNF [][]QuestionFilterLiteral
}

func (qf *QuestionFilter) ToWhere() string {
	if qf.IsTemplated() {
		//TODO raise an exception
	}
	//TODO
}

func (qf *QuestionFilter) IsTemplated() bool {
	return false //TODO TODO
}

func (qf *QuestionFilter) Resolve(round *Round, ctx *gorm.DB) {
	if !qf.IsTemplated() {
		//TODO TODO this needs to call a database now as well...
	}
}

//TODO this guy needs to be resolvable

