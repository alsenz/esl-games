package lesson2

import (
	"github.com/alsenz/esl-games/pkg/account"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// The key here is to leverage gorm Map WHERE as much as possible

type QuestionBankID uuid.UUID

type WhereCondition map[string]interface{}
var AllowedWhereColumns = map[string]bool {
	//TODO TODO need to look for the key question columns we deffo need.
	//TODO
}

type QuestionSet struct {
	QuestionBanks []QuestionBankID //TODO make jsonifiable
	WhereCondition WhereCondition //TODO make jsonifiable
}

// Should remove any conditions that aren't queryable
func (wc *WhereCondition) Cleanse() {
	//TODO just tidy up where clause for anything not allowed remove
	//TODO TODO
}

func (wc *WhereCondition) AddAuth(user account.User) {
	wc.Cleanse()
	//TODO TODO
	//TODO not sure how to do this right... TODO TODO
	(*wc)["SOME KIND OF ACCOUNT ID HERE"] = 34
}

func (qs *QuestionSet) AddWhereCondition(whereCondition WhereCondition) {
	whereCondition.Cleanse()
	qs.WhereCondition = whereCondition
}

func (qs QuestionSet) Resolve(teacher account.User, DB *gorm.DB, checkpoint Checkpoint) ([]Question, error) {
	//TODO

	qs.WhereCondition.AddAuth(teacher)
	//TODO extend the where condition a bit further - we need an account

	//TODO we need to do a resolve on checkpoint for where condition
	//TODO TODO here we need to add question banks where etc.

	//TODO do gorm query...

	return make([]Question, 0), nil
}