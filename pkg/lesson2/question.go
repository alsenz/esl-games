package lesson2

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/alsenz/esl-games/pkg/account"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

//Note: Vote here means "vote for person"- it's up to the slide designer
//... to dress this up as "vote for their answer"
type QuestionLogic string
const (
	SingleInput              QuestionLogic = "Single Input"
	SingleInputOpen          QuestionLogic = "Single Input No Right Answer"
	MultipleInput            QuestionLogic = "Multiple Input"
	MultipleInputOpen        QuestionLogic = "Multiple Input No Right Answer"
	SingleMultipleChoice     QuestionLogic = "Single Answer Multiple Choice"
	SingleMultipleChoiceOpen QuestionLogic = "Single Answer Multiple Choice No Right Answer"
	MultiMultipleChoice      QuestionLogic = "Multiple Answer Multiple Choice"
	MultiMultipleChoiceOpen  QuestionLogic = "Multiple Answer Multiple Choice No Right Answer"
	MultiBestMultipleChoice  QuestionLogic = "Multiple Answer Multiple Choice With Best Answer"
	ApproximateSingleInput   QuestionLogic = "Approximate Single Input"
	ApproximateMultipleInput QuestionLogic = "Approximate Multiple Inputs"
	VoteForPlayerFPTP				 QuestionLogic = "Vote For Player (First Past The Post)"
	VoteForPlayerSTV	           	 QuestionLogic = "Vote For Player (Single Transferable Vote)"
	VoteForPlayerAV					 QuestionLogic = "Vote For Player (Alternative Vote)"
	//Fixme: at some point - vote for team!
)
var QuestionLogics = [...]QuestionLogic{SingleInput, SingleInputOpen, MultipleInput, MultipleInputOpen, SingleMultipleChoice,
	SingleMultipleChoiceOpen, MultiMultipleChoice, MultiMultipleChoiceOpen, MultiBestMultipleChoice, ApproximateSingleInput,
	ApproximateMultipleInput, VoteForPlayerFPTP, VoteForPlayerSTV, VoteForPlayerAV}

type QuestionRules struct {
	Logic QuestionLogic	`json:"logic,omitempty"`
	ContentType ResponseContentType `json:"contentType,omitempty"`
	Data datatypes.JSON `json:"data,omitempty"` //Extra data - e.g. best answer plus multiple choices
}

func (qr *QuestionRules) Value() (driver.Value, error) {
	if raw, err := json.Marshal(qr); err != nil {
		return nil, err
	} else {
		return datatypes.JSON(raw).Value()
	}
}

func (qr *QuestionRules) Scan(src interface{}) error {
	jsn := &datatypes.JSON{}
	if err := jsn.Scan(src); err != nil {
		return err
	}
	return json.Unmarshal(*jsn, qr)
}

type Question struct { //TODO this needs to be gorm'd up  //TODO this first
	account.UserObject
	Content string	`json:"content" gorm:"column:content"`
	Header string	`json:"header,omitempty" gorm:"column:header"`
	Media string `json:"media,omitempty" gorm:"column:media"`
	MediaContentType string `json:"mediaContentType,omitempty" gorm:"column:media_content_type"`
	ByLine string	`json:"byline,omitempty" gorm:"column:by_line"`
	Tags pq.StringArray `json:"tags" gorm:"type:varchar(64)[];column:tags"`
	Rules QuestionRules `json:"rules" gorm:"column:rules"`
}

func (question *Question) Resolve(mdl *ClientModel) error {
	var err error
	if question.Content, err = mdl.Eval(question.Content); err != nil {
		return err
	}
	if question.Header, err = mdl.Eval(question.Header); err != nil {
		return err
	}
	if question.Media, err = mdl.Eval(question.Media); err != nil {
		return err
	}
	if question.MediaContentType, err = mdl.Eval(question.MediaContentType); err != nil {
		return err
	}
	if question.ByLine, err = mdl.Eval(question.ByLine); err != nil {
		return err
	}
	return nil
}