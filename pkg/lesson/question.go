package lesson

import (
	"github.com/alsenz/esl-games/pkg/account"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

//TODO we have a few options
//TODO - we have a concrete question

//TODO let's use go-templates liberally.
//TODO everything is either a question or a question filter
//TODO questions have

//TODO
//TODO 1) gorm question
//TODO 2) go variants
//TODO 3) go interfaces - mocking out response...

type Resolvable interface {
	IsTemplated() bool
	Resolve(*Round, *gorm.DB) *Question
}

type Question struct {
	account.UserObject
	Content string //TODO this may involve images of course... I think we're gonna allow this via markdown...
	Header  string
	Tags    []string //TODO gorm question - how to embed this
	//TODO response type - probably here we just need the response type
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

//TODO this basically needs to be a list of conditions
//TODO check how gorm structures things
type QuestionFilter struct {
	QuestionBanks []string
	//TODO TODO
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
