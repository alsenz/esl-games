package lesson2

import (
	"github.com/alsenz/esl-games/pkg/account"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Controller runs on the EventLoop go routine mostly,
// except a number of Async methods that take a channel to sync back with (with optional error!)
type Controller struct {
	//Note: Controller doesn't own the current round any more - event loop does that
	PlanID uuid.UUID
	Plan *Plan // This is the fixed recipe for the lesson
	Model Model // This is the dynamically evaluated model as questions are answered etc.
	//TODO this may need to be lazily loaded like the plan...
	LeaderID uuid.UUID
	Leader *account.User //Runs the show
	//TODO Players
	//TODO Teams - plus membership

	conn *gorm.DB
	usedQuestions map[uuid.UUID]bool //A little cache to ensure we don't pull questions multiple times
}

func NewController(planID uuid.UUID, leaderID uuid.UUID) *Controller {
	//TODO need to make the gorm connection some 'ow
	//TODO and it makes sense to make this here...?
	return &Controller{
		planID,
		nil,
		Model{},
		leaderID, //TODO make this here...
		nil,
		nil,
		make(map[uuid.UUID]bool),
	}
}

func (ctrl *Controller) LoadPlan() error {
	zap.L().Info("Loading Plan (Loading Leader User)")
	result:= ctrl.conn.First(ctrl.Leader, ctrl.LeaderID)
	if result.Error != nil {
		ctrl.Leader = nil // Reset required to make PlanIsLoaded work
		return result.Error
	}
	zap.L().Info("Loading Plan")
	result = ctrl.conn.First(ctrl.Plan, ctrl.PlanID)
	if result.Error != nil {
		ctrl.Plan = nil // Reset required to make PlanIsLoaded work
		return result.Error
	}
	return nil
}

func (ctrl *Controller) assertPermissions(questions []Question) []Question {
	result := make([]Question, 0 , len(questions))
	for _, qs := range questions {
		if ctrl.Leader.CanReadObject(&qs.UserObject) {
			result = append(result, qs)
		}
	}
	return result
}

//TODO move this into the controller TODO TODO
//TODO should this in fact perhaps sit on the controller? TODO TODO
//TODO pick this up in a bit- would be good to make this ASYNC!
//TODO TODO

//Note - let's just get rid of the idea that question filters have templates
//Turn a question set into a list of questions and resolve those questions using
//a client model //Fixme: make aspects of this async!
func (ctrl *Controller)  Resolve(qs *QuestionSet) ([]Question, error) {
	if err := qs.WhereCondition.CleanseAndValidate(); err != nil {
		return nil, err
	}
	qs.WhereCondition.AddAuth(*ctrl.Leader)
	var numberToFetch int = 0
	switch qs.Logic {
	case QuestionSetLogicOneQuestionForEverybody:
		numberToFetch = 1 + len(ctrl.usedQuestions)
	case QuestionSetLogicOneQuestionPerPlayer:
		//TODO add this method to controller
		numberToFetch = ctrl.NumPlayers() + len(ctrl.usedQuestions)
	case QuestionSetLogicOneQuestionPerTeam:
		//TODO add this method to controller
		numberToFetch = ctrl.NumTeams() + len(ctrl.usedQuestions)
	}
	questions := make([]Question, 0, numberToFetch)
	result := ctrl.conn.Where(qs.WhereCondition).Limit(numberToFetch).Find(&questions)
	if(result.Error != nil) {
		return nil, result.Error
	}
	questions = ctrl.assertPermissions(questions)
	questionsFiltered := ctrl.filterOutVisitedQuestions(questions) //TODO
	//TODO next step -
	//TODO filter out the cached questions
	//TODO assign the questions to teams.
	//TODO finally - resolve the questions using client models
	return questions, nil
}