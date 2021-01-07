package lesson

import (
	"github.com/alsenz/esl-games/pkg/account"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type Plan struct {
	account.UserObject
	Name string
	Description string
	TeamStructures TeamStructures //This is embedded
	Acts []Act // This is also embedded.
}

//TODO: this needs a better name. It's what the planner gives the controller to run the next round with
type NextQuestionDelivery struct {
	NextRoundDef RoundIdx //Planner decides what the next round is
	//TODO
}

type Planner struct {
	PlanID uuid.UUID
	Plan Plan
	// Channels from the event loop
	NextRoundChannelIn <-chan ReadyForNextRound
	// Channels with the controller
	//TODO
	QuestionFilterResolverOut chan<- QuestionFilterQuery
	QuestionFilterResolverIn <-chan QuestionFilterQuery
	//TODO bugger - HTF does Planner know what the next round will be?
	//TOOD TODO revisit

}

//TODO not sure if this exists anymore...? TODO TODO move this to planner
func (planner *Planner) LoadPlan(planId uuid.UUID, planChannel chan Plan) {
	zap.L().Info("Loading Plan")
	//TODO gorm out the plan and send it into a channel of plans
}

//TODO any utilities here that would help
