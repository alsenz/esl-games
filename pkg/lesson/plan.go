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
	QuestionRetrievalChannel chan<- QuestionDraw
}

func NewPlanner(planID uuid.UUID, conn *gorm.DB) *Planner {
	//TODO review questionlink here
	return &Planner{planID, nil, conn, make(chan map[QuestionLink]Question, 8)}
}

func (planner *Planner) Start() error {
	return planner.LoadPlan()
}

func (planner *Planner) LoadPlan() error {
	zap.L().Info("Loading Plan")
	planner.Plan = Plan{}
	result:= planner.conn.First(&planner.Plan, planner.PlanID)
	return result.Error
}

//TODO change this - should be for a question filter - a resolved one at that probably
//TODO need to write a simple routine for running a question filter
func (planner *Planner) FetchQuestions(rqf ResolvedQuestionFilter) error {
	var matched []Question
	//TODO no we wanna build the clause separately.
	//TODO we wanna on question filter build a map[string] interface
	planner.conn.Where(map[string]interface{}{}).Find(&matched)
	result:= as.conn.Where("id = ? AND md5sum = ?", assetId, sig).First(&asset)
	//TODO shove them back in the channel TODO TODO
}
