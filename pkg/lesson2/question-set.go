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

//TODO move onto controller

//TODO "has permission..." nice simplification

//TODO move this into the controller
func (qs *QuestionSet) filterByPermissions(teacher account.User, questions []Question) []Question {
	//TODO pick this up next
	result := make([]Question, 0 , len(questions))
	for _, qs := range questions {
		if qs.OwnerID == teacher.ID {

		}
		for tGrpID := range teacher.GroupIDs() {
			if tGrpID == qs.GroupID {
				//check the
				break
			}
		}
	}
	//TODO: note - we'll need a check permissions on the questions
	//TODO: as we'll need to make sure that the group read is set!
	//TODO filter based on result
}

//TODO move this into the controller TODO TODO
//TODO should this in fact perhaps sit on the controller? TODO TODO
//TODO pick this up in a bit- would be good to make this ASYNC!
//TODO TODO

//TODO controller should know who teacher is for e.g.
//Note - let's just get rid of the idea that question filters have templates
func Resolve(qs *QuestionSet, teacher account.User, DB *gorm.DB) ([]Question, error) {
	if err := qs.WhereCondition.CleanseAndValidate(); err != nil {
		return nil, err
	}
	qs.WhereCondition.AddAuth(teacher)
	//TODO TODO we need to know what the limit is - and that will depend on
	//TODO check this here - make sure it isn't part of question logic...
	//TODO what the question mapping logic is- one-question, one-per-player, one-per-team
	WHAT_IS_THE_LIMIT := 5 //TODO
	questions := make([]Question, 0, WHAT_IS_THE_LIMIT)
	result := DB.Where(qs.WhereCondition).Limit(WHAT_IS_THE_LIMIT).Find(&questions)
	if(result.Error != nil) {
		return nil, result.Error
	}
	//TODO - make this happen
	questions = qs.filterByPermissions(teacher, questions)
	return questions, nil
}