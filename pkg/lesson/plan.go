package lesson

import (
	"errors"
	"github.com/alsenz/esl-games/pkg/account"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Plan struct {
	account.UserObject
	Name string
	Description string
	TeamStructures TeamStructures //This is embedded
	Acts []Act // This is also embedded.
}

type Planner struct {
	PlanID uuid.UUID
	Plan Plan
	conn *gorm.DB
	//A channel for making some of the fetch from the db async
	QuestionRetrievalChannel chan map[QuestionLink]Question
}

func NewPlanner(planID uuid.UUID, conn *gorm.DB) *Planner {
	return &Planner{planID, nil, conn, make(chan map[QuestionLink]Question, 8)}
}

func (planner *Planner) Start() error {
	return planner.LoadPlan()
}

func (planner *Planner) LoadPlan() error {
	zap.L().Info("Loading Plan")
	planner.Plan = Plan{}
	result:= planner.conn.Where("id = ?", planner.PlanID).First(&planner.Plan)
	return result.Error
}

//TODO need to write a simple routine for running a quesiton filter
func (planner *Planner) FetchQuestions(query *QuestionFilterQuery) error {
	if(query.IsTemplated()) {
		return errors.New("Can't fetch templated question, needs to be resolved first")
	}
	//TODO we wanna call resolve on a model first...
	//TODO TODO review and check the form api again.
	asset := &Question{} //TODO chekc multi... TODO TODO
	result:= as.conn.Where("id = ? AND md5sum = ?", assetId, sig).First(&asset)
	//TODO shove them back in the channel TODO TODO
}
