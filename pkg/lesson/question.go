package lesson

import (
	"github.com/alsenz/esl-games/pkg/account"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strings"
)

//TODO let's simplify further and say we have a question

//TODO content types...

// We have a parameter N clearly
//SingleInput, MultipleInput, SingleMultipleChoice, MultiMultipleChoice,
//ApproximateInput (has error margin)
//RankedMultipleChoice
//FPTPVote, SingleTransferableVote, ProportionalVote, AlternativeVote (often has N)

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
)

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

var QuestionLogics = [...]QuestionLogic{}

//TODO let's use go-templates liberally.
//TODO everything is either a question or a question filter
//TODO question filter -> resolve to question filter -> question -> resolve to question
//TODO slides now ALWAYS have a question filter

type QuestionBehaviour struct {
	Logic QuestionLogic //TODO json these up!
	ContentType QuestionContentType
	Tags []string //May include x=y pairs
}


//TODO
//TODO 1) gorm question - doesn't work :( TODO unless we use the JSON type in postgres... TODO this is gonna be pretty important.

type Resolvable interface {
	IsTemplated() bool
	Resolve(*Round, *gorm.DB) *Question
}

type Question struct {
	account.UserObject
	Content string
	Header string
	Image string
	//TODO will have to split this manually TODO TODO JSON this
	Tags []string //TODO gorm question - how to embed this - ah well that's not happening apparently... TODO TODO no json it!
	Behaviour QuestionBehaviour //TODO json this up!
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

type QuestionFilterOp string

const (
	QuestionFilterOpEquals          QuestionFilterOp = "="
	QuestionFilterOpNotEquals       QuestionFilterOp = "!="
	QuestionFilterLessThan          QuestionFilterOp = "<"
	QuestionFilterLessThanEquals    QuestionFilterOp = "<="
	QuestionFilterGreaterThan       QuestionFilterOp = ">"
	QuestionFilterGreaterThanEquals QuestionFilterOp = ">="
	QuestionFilterContains          QuestionFilterOp = "contains"
	QuestionFilterStartsWith        QuestionFilterOp = "starts-with"
)

var QuestionFilterOps = [...]QuestionFilterOp{QuestionFilterOpEquals, QuestionFilterOpNotEquals, QuestionFilterLessThan,
	QuestionFilterLessThanEquals, QuestionFilterGreaterThan, QuestionFilterGreaterThanEquals, QuestionFilterContains,
	QuestionFilterStartsWith}

type QuestionFilterLiteral struct {
	Field string //TODO need to json this.
	Op    QuestionFilterOp
	Value string
}

func (qfl QuestionFilterLiteral) ToString() string {
	return strings.Join([]string{qfl.Field, string(qfl.Op), qfl.Value}, "")
}

//TODO this basically needs to be a list of conditions
//TODO check how gorm structures things
type QuestionFilter struct {
	QuestionBanks []string
	//TODO TODO this just needs to get translated into a nice big where query
	CNF [][]QuestionFilterLiteral
}

func (qf *QuestionFilter) ToWhere() string {
	if qf.IsTemplated() {
		//TODO raise an exception
	}
	//TODO
}

func (qf *QuestionFilter) IsTemplated() bool {
	return false //TODO TODO
}

func (qf *QuestionFilter) Resolve(round *Round, ctx *gorm.DB) {
	if !qf.IsTemplated() {
		//TODO TODO this needs to call a database now as well...
	}
}

//TODO this guy needs to be resolvable
