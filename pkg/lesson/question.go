package lesson

import (
	"github.com/alsenz/esl-games/pkg/account"
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
	ContentTypeVideo	QuestionContentType = "video"
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
	TagsJSON      datatypes.JSON      `json:"tags"`
	RulesJSON datatypes.JSON `json:"rules"`
}

func (q *Question) MakeTags() []string {
	//TODO TODO need to conver tthis to tas string
	return q.TagsJSON
}

func (q *Question) SetTags(tags []string) {
	//TODO how does this become JSON?
	q.TagsJSON =
}

func (q *Question) MakeRules() QuestionRules {
	//TODO need to convert this plus validate it...
}

func (q *Question) IsTemplated() bool {
	//TODO check that either of Content or Header is templated.
	return false
}

//TODO TODO change this actually we wanna copy this thing...
// Note- this resolves in place!
func (q *Question) Resolve(round *Round, _ *gorm.DB) {
	if !q.IsTemplated() {
		q.ID = uuid.NewV4()                             //We need to give ourselves a new UUID since technically this is a new question
		q.Content = "TODO need to turn into a template" //TODO turn into a tempalte
	}
}

