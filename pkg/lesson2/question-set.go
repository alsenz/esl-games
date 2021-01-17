package lesson2

import (
	"errors"
	"github.com/alsenz/esl-games/pkg/account"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// The key here is to leverage gorm Map WHERE as much as possible


type WhereCondition map[string]interface{}
var AllowedWhereColumns = map[string]bool {
	"content": true,
	"header": true,
	"by_line": true,
	"tags": true,
	"group_id": true,
	"owner_id": true,
	"id": true,
}

type QuestionSetLogic string
const (
	QuestionSetLogicOneQuestionForEverybody QuestionSetLogic = "one-question-for-everbody"
	QuestionSetLogicOneQuestionPerPlayer QuestionSetLogic = "one-question-per-player"
	QuestionSetLogicOneQuestionPerTeam QuestionSetLogic = "one-question-per-team"
)

type QuestionSet struct {
	Logic QuestionSetLogic `json:"logic"`
	WhereCondition WhereCondition `json:"whereCondition,omitempty"`
}

// Should remove any conditions that aren't queryable
// Should check that all the types are right
func (wc *WhereCondition) CleanseAndValidate() error {
	// Clean out any "disallowed" keys
	for key, _ := range *wc {
		if val, found := AllowedWhereColumns[key]; !found || !val {
			delete(*wc, key)
		}
	}
	for key, val := range *wc {
		switch key {
		case "content":
			fallthrough
		case "by_line":
			fallthrough
		case "header":
			if _, ok := val.(string); !ok {
				if _, ok = val.([]string); !ok {
					return errors.New("cannot validate WhereCondition, type of \""+key+"\" is" +
						" not string or []string")
				}
			}
		case "tags":
			newVal := make(pq.StringArray, 0)
			if strVal, ok := val.(string); ok {
				newVal = append(newVal, strVal)
				(*wc)[key] = newVal
			} else if sliceVal, ok := val.([]string); ok {
				(*wc)[key] = pq.StringArray(sliceVal)
			} else {
				errors.New("cannot validate WhereCondition, type of \""+key+"\" is" +
					" not string or []string")
			}
		case "owner_id":
			fallthrough
		case "id":
			fallthrough
		case "group_id":
			if strVal, ok := val.(string); ok {
				newUUID, err := uuid.FromString(strVal)
				if err != nil {
					return err
				}
				(*wc)[key] = newUUID
			} else if sliceVal, ok := val.([]string); ok {
				newVal := make([]uuid.UUID, 0)
				for _, val := range sliceVal {
					newUUID, err := uuid.FromString(val)
					if err != nil {
						return err
					}
					newVal = append(newVal, newUUID)
				}
				(*wc)[key] = newVal
			} else {
				errors.New("cannot validate WhereCondition, type of \""+key+"\" is" +
					" not string or []string")
			}

		}
	}
}

func (wc *WhereCondition) AddAuth(user account.User) {
	(*wc)["owner_id"] = user.ID
	(*wc)["group_id"] = user.GroupIDs()
}