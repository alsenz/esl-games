package lesson2

import (
	"errors"
	"fmt"
	"github.com/alsenz/esl-games/pkg/account"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"sort"
	"strconv"
)

type Draw map[*QuestionSet]map[Round][]Question

// Controller runs on the EventLoop go routine mostly,
// except a number of Async methods that take a channel to sync back with (with optional error!)
type Controller struct {
	//Note: Controller doesn't own the current round any more - event loop does that
	PlanID uuid.UUID
	Plan *Plan // This is the fixed recipe for the lesson
	Model Model // This is the dynamically evaluated model as questions are answered etc.
	LeaderID uuid.UUID
	Leader *account.User //Runs the show
	Players map[ClientID]Player
	Teams map[ClientID]TeamName
	Rounds []Round
	Draw Draw

	conn *gorm.DB
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
		make([]Round, 0),
		make(Draw),
		nil,
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
			ctrl.Teams[player.ClientID] = DefaultTeamName
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


func (ctrl *Controller) AddPlayer(player Player) {
	//Fixme: anything else?
	ctrl.Players[player.ClientID] = player
}

// Returns the index of the smallest round which is greater or equal to this one
func (ctrl *Controller) getRoundIdxLeastUpperBound(round Round) int {
	sort.Search(len(ctrl.Rounds), func(i int) bool { return !ctrl.Rounds[i].LessThan(round) })
}

// Returns the index of the greatest round less than or equal to
func (ctrl *Controller) getRoundIdxLowerBound(round Round) int {
	ub := ctrl.getRoundIdxLeastUpperBound(round)
	if round == ctrl.Rounds[ub] || ub == 0 {
		return ub
	} else {
		return ub - 1
	}
}

func (ctrl *Controller) insertRound(round Round) {
	if len(ctrl.Rounds) == 0 {
		ctrl.Rounds = append(ctrl.Rounds, round)
		return
	}
	ub := ctrl.getRoundIdxLeastUpperBound(round)
	if ctrl.Rounds[ub] == round {
		return
	}
	ctrl.Rounds = append(ctrl.Rounds[:ub+1], ctrl.Rounds[ub:]...)
	ctrl.Rounds[ub] = round
}


func(ctrl *Controller) GetPlayerTeam(player Player) TeamName {
	if team, found := ctrl.Teams[player.ClientID]; found {
		return team
	}
	return DefaultTeamName
}

//TODO need to pass through a draw channel
// Returns whether or not a round was created - i.e. false if the round already exists
func (ctrl *Controller) CreateRound(round Round, drawChan *chan<-Draw) (bool, error) { //TODO draw channel
	if _, found := ctrl.Model[round]; found {
		return false, nil
	}
	ctrl.insertRound(round)
	for _, player := range ctrl.Players {
		ctrl.Model.AddRecord(Record{
			player,
			round,
			make(ScoreCards),
			nil, //We'll assign these later
			ctrl.GetPlayerTeam(player),
			nil,
		})
	}
	if drawChan != nil {
		if err := ctrl.FetchQuestions(round, *drawChan); err != nil {
			return false, err
		}
	}
	return true, nil
}

//TODO pick up from here
//TODO we need a simple add response!

//TODO we need a simple score round! //TODO -- this will need to score forward and update those scores in future rounds...

//TODO we need a simple (pre)load scene!
//TODO we need a ResolveRound that will use the Draw map or questions pre-assigned in the model to resovle the screens
//TODO TODO

func (ctrl *Controller) trimRepeatQuestions(qs *QuestionSet, questions []Question) []Question {
	// Proceed to identify new questions by counter example
	newQuestions := make(map[uuid.UUID]Question)
	for _, q := range questions {
		newQuestions[q.ID] = q
	}
	oldQuestions := make([]Question, 0)
	for _, inner := range ctrl.Draw {
		for _, questions := range inner {
			for _, existingQ := range questions {
				if newQuestionDisproof, found := newQuestions[existingQ.ID]; found {
					oldQuestions = append(oldQuestions, newQuestionDisproof)
					delete(newQuestions, newQuestionDisproof.ID)
				}
			}
		}
	}
	// Now we'll need to add some old questions back in
	numRequired := ctrl.NumberRequired(qs)
	result := make([]Question, 0, numRequired)
	for _, newQ := range newQuestions {
		result = append(result, newQ)
	}
	for i := len(result); i <= numRequired; i++ {
		oldIdx := i - len(result)
		if oldIdx < len(oldQuestions) {
			result = append(result, oldQuestions[oldIdx])
		}
	}
	return result
}

func (ctrl *Controller) UsedQuestionCount() int {
	ct := 0
	for _, inner := range ctrl.Draw {
		for _, inner := range inner {
			ct += len(inner)
		}
	}
	return ct
}


func (ctrl *Controller) NumberRequired(qs *QuestionSet) int {
	switch qs.Logic {
	case QuestionSetLogicOneQuestionForEverybody:
		return 1
	case QuestionSetLogicOneQuestionPerPlayer:
		return ctrl.NumPlayers()
	case QuestionSetLogicOneQuestionPerTeam:
		return ctrl.NumTeams()
	}
	return 0
}

//Note - let's just get rid of the idea that question filters have templates
//Turn a question set into a list of questions and resolve those questions using
//a client model
func (ctrl *Controller) FetchQuestions(round Round, drawChan chan<-Draw) error {
	//First check if we can find a question set in the plan
	registerRound := 1
	if uint(len(ctrl.Plan.Acts) + registerRound) < round.Act {
		return errors.New(fmt.Sprintf("Act number %s not found in plan", round.Act))
	}
	act := ctrl.Plan.Acts[round.Act]
	if uint(len(act.Scenes) + registerRound) < round.Scene {
		return errors.New(fmt.Sprintf("Act number %s, Scene number %s not found in plan", round.Act, round.Scene))
	}
	scene := act.Scenes[round.Scene]
	if scene.QuestionSet == nil {
		return nil
	}
	if err := scene.QuestionSet.WhereCondition.CleanseAndValidate(); err != nil {
		return err
	}
	numberToFetch := ctrl.NumberRequired(scene.QuestionSet) + ctrl.UsedQuestionCount()

	go func() {
		//TODO check this
		//TODO this needs to be in a go routine TODO TODO
		query := ctrl.conn.Where(scene.QuestionSet.WhereCondition)
		query = ctrl.Leader.InjectReadAuth(query)
		questions := make([]Question, 0, numberToFetch)
		result := query.Limit(numberToFetch).Find(&questions)
		if(result.Error != nil) {
			zap.L().Fatal("Error retrieving questions: " + result.Error.Error())
		}
		//TODO put into the channel
	}()
	return  nil
}

func (ctrl *Controller) addToDraw(qs *QuestionSet, round Round, questions []Question) {
	if _, found := ctrl.Draw[qs]; !found {
		ctrl.Draw[qs] = make(map[Round][]Question)
	}
	if _, found := ctrl.Draw[qs][round]; !found {
		ctrl.Draw[qs][round] = make([]Question, len(questions))
	}
	ctrl.Draw[qs][round] = append(ctrl.Draw[qs][round], questions...)
}

//TODO we need a DealDraw(Draw) which should extend draw map and try to pre-resolve as much as possible
func (ctrl *Controller) DealDraw(draw Draw) error {

	for qs, inner := range draw {
		for round, questions := range inner {
			if len(questions) == 0 {
				return errors.New("no questions drawn from question set")
			}
			questionsFiltered := ctrl.trimRepeatQuestions(qs, questions)
			// We know we're going to use these, so we can add them to the draw
			ctrl.addToDraw(qs, round, questionsFiltered)
			// Assign to teams / players and resolve
			switch qs.Logic {
			case QuestionSetLogicOneQuestionForEverybody:
				if teamMap, found := ctrl.Model[round]; found {
					for team, playerMap := range teamMap {
						for playerId, record := range playerMap {
							playerQ := questions[0]
							playerQ.Resolve(&ClientModel{
								ctrl.Model,
								record.Player,
								round,
							})
							record.Question = &playerQ
							ctrl.Model[round][team][playerId] = record
						}
					}
				} else {
					return errors.New("round not found in model")
				}

			case QuestionSetLogicOneQuestionPerTeam:
				//TODO - do similar to above
			case QuestionSetLogicOneQuestionPerPlayer:
				//TODO
			default:
				zap.L().Fatal("Logic error - unknown question set logic: " + string(qs.Logic))
			}
			//TODO
			//TODO next step -
			//TODO assign the questions to teams-- not clear how this is going to happen, but we'll make it work
			//TODO finally - resolve the questions using client models and add them to the model
		}
	}


}