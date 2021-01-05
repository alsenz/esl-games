package lesson

import (
	"github.com/gofrs/uuid"
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

type QuestionIDList []uuid.UUID

type QuestionFilterQuery struct {
	QuestionBanks []string
	//TODO TODO this just needs to get translated into a nice big where query
	CNF [][]QuestionFilterLiteral
}

func (qfq *QuestionFilterQuery) IsTemplated() bool {
	return false //TODO TODO
}

func (qfq *QuestionFilterQuery) ToWhereClause() string {
	if qfq.IsTemplated() {
		//TODO raise an exception
	}
	//TODO
}

//TODO this basically needs to be a list of conditions
//TODO check how gorm structures things
type QuestionFilter struct {
	QuestionIds []uuid.UUID //Either-or
	Query QuestionFilterQuery //Eitehr-or
}

func (qf *QuestionFilter) ToWhereClause() string {
	if qf.IsTemplated() {
		//TODO raise an exception
	}
	//TODO
}

func (qf *QuestionFilter) IsTemplated() bool {
	return false //TODO TODO
}

// usedQuesIds is to ensure we don't have any repeats
func (qf *QuestionFilter) Resolve(round *Round, ctx *gorm.DB, usedQuesIds map[uuid.UUID]bool) {
	if !qf.IsTemplated() {
		//TODO TODO this needs to call a database now as well...
	}
	//TODO we need an ORDER BY where we put the used questions at the bottom of the list!
}

