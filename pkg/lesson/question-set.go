package lesson

import (
	"errors"
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
	CNF [][]QuestionFilterLiteral
}

type QuestionSet struct {
	QuestionIds []uuid.UUID `json:"questionIds,omitempty"` //Either-or
	Query QuestionFilterQuery `json:"query,omitempty"` //Either-or
}

// usedQuesIds is to ensure we don't have any repeats
func (qf QuestionSet) Resolve(model *Model, conn *gorm.DB, usedQuestionSet map[uuid.UUID]bool) () {

	//TODO we need to query out
}

type ResolvedQuestionFilter struct {
	QuestionSet
}

func (rqf ResolvedQuestionFilter) ToWhereClause(tx *gorm.DB) (*gorm.DB, error) {
	// Either we use a fixed list of question ids or we use a filter, not both
	if(len(rqf.QuestionIds) > 0) {
		return tx.Where(rqf.QuestionIds), nil

	} else { // Use the filter
		//TODO
	}
}



