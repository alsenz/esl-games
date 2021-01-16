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

//TODO really not sure how question logic solves the "vote on a canvas response" setup
//TODO really not sure how we draw this nicely as well...
//TODO how do we setup the question multi choices response in the console...?
//TODO how do we tie these together? Not easy.
//TODO how do we tie that each response is assoc with a Player for voting?
//TODO flibbit this is broken here as well isn't it? Urgh.

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

type QuestionDraw map[RoundIdx]map[ClientID]Question

type Question struct {
	account.UserObject
	Content string	`json:"content"`
	Header string	`json:"header,omitempty"`
	Image uuid.UUID `json:"image,omitempty"`
	ByLine string	`json:"byline,omitempty"`
	Tags pq.StringArray `json:"tags" gorm:"type:varchar(64)[]"`
	Rules QuestionRules `json:"rules"`
}

type ResolvedQuestion struct {
	Question
}

func (q *Question) Resolve(mdl *Model, _ *gorm.DB) (* ResolvedQuestion, error) {
	base := q.UserObject
	base.ID = uuid.NewV4() // Regenerate ID to denote that it's a different question.
	var content, header, byline string
	var err error
	if content, err = mdl.Eval(q.Content); err != nil {
		return nil, err
	}
	if header, err = mdl.Eval(q.Header); err != nil {
		return nil, err
	}
	if byline, err = mdl.Eval(q.ByLine); err != nil {
		return nil, err
	}
	return &ResolvedQuestion{
		Question{
			base,
			content,
			header,
			q.Image,
			byline,
			q.Tags,
			q.Rules,
		},
	}, nil
}

