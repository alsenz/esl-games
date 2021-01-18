package lesson2

import (
	"errors"
	"github.com/alsenz/esl-games/pkg/account"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
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
	Players map[ClientID]Player
	Teams map[ClientID]TeamName

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
		make(map[ClientID]Player),
		make(map[ClientID]TeamName),
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

func (ctrl *Controller) NumPlayers() int {
	return len(ctrl.Players)
}

func (ctrl *Controller) NumTeams() int {
	return len(ctrl.Teams)
}

func (ctrl *Controller) PutPlayersInTeams() error {
	switch ctrl.Plan.TeamRules.TeamLogic {
		case TeamLogicNoTeams:
		for _, player := range ctrl.Players {
			ctrl.Teams[player.ClientID] = "default"
		}
		case TeamLogicByNumberOf:
		idx := 0
		for _, player := range ctrl.Players {
			ctrl.Teams[player.ClientID] = TeamName("Team " + strconv.FormatInt(int64((idx % ctrl.Plan.TeamRules.Count) + 1), 10))
			idx++
		}
		case TeamLogicBySize: //Fixme - fill slightly more intelligent
		teamNo := 1
		size := 0
		for _, player := range ctrl.Players {
			ctrl.Teams[player.ClientID] = TeamName("Team " + strconv.FormatInt(int64(teamNo), 10))
			size++
			if size >= ctrl.Plan.TeamRules.Count {
				teamNo++
				size = 0
			}
		}
		default:
			return errors.New("Unknown team logic: " + string(ctrl.Plan.TeamRules.TeamLogic))
	}
	return nil
}

//TODO pick up from here

//TODO we need a simple addPlayer!

//TODO we need a simple add round!

//TODO we need a simple add response!

//TODO we need a simple score round!

//TODO we need a simple (pre)load scene!

func (ctrl *Controller) assertPermissions(questions []Question) []Question {
	result := make([]Question, 0 , len(questions))
	for _, qs := range questions {
		if ctrl.Leader.CanReadObject(&qs.UserObject) {
			result = append(result, qs)
		}
	}
	return result
}

//TODO pick up with removeVisited

//Note - let's just get rid of the idea that question filters have templates
//Turn a question set into a list of questions and resolve those questions using
//a client model //Fixme: make aspects of this async!
func (ctrl *Controller)  Resolve(qs *QuestionSet) ([]Question, error) {
	if err := qs.WhereCondition.CleanseAndValidate(); err != nil {
		return nil, err
	}
	qs.WhereCondition.AddAuth(*ctrl.Leader)
	var numberToResolve int = 0
	switch qs.Logic {
	case QuestionSetLogicOneQuestionForEverybody:
		numberToResolve = 1
	case QuestionSetLogicOneQuestionPerPlayer:
		numberToResolve = ctrl.NumPlayers()
	case QuestionSetLogicOneQuestionPerTeam:
		numberToResolve = ctrl.NumTeams()
	}
	numberToFetch := numberToResolve + len(ctrl.usedQuestions)
	questions := make([]Question, 0, numberToFetch)
	result := ctrl.conn.Where(qs.WhereCondition).Limit(numberToFetch).Find(&questions)
	if(result.Error != nil) {
		return nil, result.Error
	}
	questions = ctrl.assertPermissions(questions)
	//TODO build this
	questionsFiltered := ctrl.removedVisitedQuestions(questions, numberToResolve)
	//TODO next step -
	//TODO filter out the cached questions
	//TODO assign the questions to teams.
	//TODO finally - resolve the questions using client models
	return questions, nil
}