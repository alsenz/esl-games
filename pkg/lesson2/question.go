package lesson2

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/alsenz/esl-games/pkg/account"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type QuestionContentType string
const (
	ContentTypeText		QuestionContentType = "text"
	ContentTypeCanvas	QuestionContentType = "canvas"
	ContentTypeImage	QuestionContentType = "image"
	ContentTypeVideo	QuestionContentType = "gif"
	ContentTypeNumber	QuestionContentType = "number"
	ContentTypeEmoji	QuestionContentType = "emoji"
	ContentTypeBoolean	QuestionContentType = "boolean"
	ContentTypeMathjax	QuestionContentType = "mathjax"
	ContentTypeAudio	QuestionContentType = "audio"
	ContentTypeChart	QuestionContentType = "chart"
)
var QuestionContentTypes = [...]QuestionContentType{ContentTypeText, ContentTypeCanvas, ContentTypeImage,
	ContentTypeVideo, ContentTypeNumber, ContentTypeEmoji, ContentTypeBoolean, ContentTypeMathjax, ContentTypeAudio,
	ContentTypeChart}

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
	VoteFPTP				 QuestionLogic = "Vote (First Past The Post)"
	VoteSTV					 QuestionLogic = "Vote (Single Transferable Vote)"
	VotePR					 QuestionLogic = "Vote (Proportional Scoring)"
	VoteAV					 QuestionLogic = "Vote (Alternative Vote)"
)
var QuestionLogics = [...]QuestionLogic{SingleInput, SingleInputOpen, MultipleInput, MultipleInputOpen, SingleMultipleChoice,
	SingleMultipleChoiceOpen, MultiMultipleChoice, MultiMultipleChoiceOpen, MultiBestMultipleChoice, ApproximateSingleInput,
	ApproximateMultipleInput, VoteFPTP, VoteSTV, VotePR, VoteAV}

type QuestionRules struct {
	Logic QuestionLogic	`json:"logic,omitempty"`
	ContentType QuestionContentType `json:"contentType,omitempty"`
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

type Question struct {
	account.UserObject
	Content string	`json:"content"`
	Header string	`json:"header,omitempty"`
	Media string `json:"media,omitempty"`
	MediaContentType string `json:"mediaContentType,omitempty"`
	ByLine string	`json:"byline,omitempty"`
	Tags pq.StringArray `json:"tags" gorm:"type:varchar(64)[]"`
	Rules QuestionRules `json:"rules"`
}

//TODO return to QuestionSet after this and resolve...
//TODO circle back round and sort out model
//TODO needs to be resolved based on model
func (question *Question) Resolve(mdl *Model) {
	//TODO just eval everything that can be TODO TODO
	question.Content = mdl.Eval(question.Content)
}