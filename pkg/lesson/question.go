package lesson

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/alsenz/esl-games/pkg/account"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

//Note: scoring systems for each of these... there's going to need to be a validate...

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
	ScoringRules ScoringRules	`json:"scoringRules,omitempty"`
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




//TODO let's use go-templates liberally.
//TODO everything is either a question or a question filter
//TODO question filter -> resolve to question filter -> question -> resolve to question
//TODO slides now ALWAYS have a question filter
//TODO when reading off inputs, we lookup resolved questions by rounds Act 1, Scene 3, Repetition 5- nice 'n' simple

//TODO
//TODO use gorm datatypes of string - will need to write validate routine for question to make sure is correct
//TODO 1) gorm question - doesn't work :( TODO unless we use the JSON type in postgres... TODO this is gonna be pretty important.



type Question struct {
	account.UserObject
	Content string	`json:"content"`
	Header string	`json:"header,omitempty"`
	Image uuid.UUID `json:"image,omitempty"`
	ByLine string	`json:"byline,omitempty"`
	TagsJSON pq.StringArray `json:"tags" gorm:"type:varchar(64)[]"`
	Rules QuestionRules `json:"rules"`
	//TODO answers?
}

type ResolvedQuestion struct {
	Question
}

//TODO
//TODO TODO change this actually we wanna copy this thing...
// Note- this resolves in place!
func (q *Question) Resolve(ctx *Context, _ *gorm.DB) *ResolvedQuestion {
	q.ID = uuid.NewV4()                             //We need to give ourselves a new UUID since technically this is a new question
	q.Content = "TODO need to turn into a template" //TODO turn into a template
}

