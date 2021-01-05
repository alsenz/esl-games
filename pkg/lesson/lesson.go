package lesson

import (
	"errors"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"time"
)

//TODO lesson's gonna need a "status" API to say what auth is required - the web app can then go <loading/><auth/><noauth/> type thing.
//TODO TODO

type Register struct {
	Timeout time.Duration
	Done bool
	RequireLogin bool
	OptDomain *string
	LessonCode string
	RegistrationSync chan bool //This is to listen for the "early ready" event of someone ending the registration period mnaually
}

type Message struct {
	Type string
	Data interface{}
}

type Client struct {
	Send chan Message
}

//TODO TODO this could get messy - we need a much clearer plan about who owns which variables and what the event loop does...
//TODO TODO I wonder if the event loop should handle all the state, in which case it may as well listen for messages itself...

//TODO no - the event loop should definitely be separate, which I fear means there is absolutely no state in the event loop...
//TODO hmm... which might essentially be fine... except that the event loop is going to ALSO wanna resolve quesitons
//TODO whereas the event handler is gonna wanna run the context...
//TODO and clients is largely gonna be handled by a hub type thing... hmm

type Lesson struct {
	Plan Plan
	Registration Register
	ResolvedQuestions map[RoundIdx]map[QuestionLink]Question	//Maps resolved questions asked at each round
	Context *Context //Used for state keeping and evaluation of strings
	InboundMessages chan Message
	Clients map[string]Client //Clients mapped by player name (//TODO or client id?)
}

func NewLesson(register Register) * Lesson {
	return &Lesson{nil, register, make(map[RoundIdx]map[QuestionLink]Question),
		NewContext(), make(chan Message, 8), make(map[string]Client)}
}

func (lesson *Lesson) Run(planId uuid.UUID) error {
	//TODO a go-routine for listening to the main channel - this is key,
	go lesson.HandleEvents() //TODO and also nick this from the websocket example...
	if err := lesson.LessonStart(planId); err != nil {
		return err
	}
	//TODO we need to start the first act now
	return nil //TODO - how do we fish out errors from HandleEvents?
}

//TODO does this even exist or does the client handle it?
func (lesson *Lesson) HandleEvents() error {
	//TODO have a look at websockets example for how we handle these events....
	//TODO probably a big range and switch based thing...
	return nil
}

func (lesson *Lesson) LoadPlan(planId uuid.UUID, planChannel chan Plan) {
	zap.L().Info("Loading Plan")
	//TODO gorm out the plan and send it into a channel of plans
}

func (lesson *Lesson) BeforeLessonStart(planId uuid.UUID) (chan Plan, error) {
	planChannel := make(chan Plan)
	go lesson.LoadPlan(planId, planChannel)
	zap.L().Info("BeforeLessonStart")
	return planChannel, nil
}

func (lesson *Lesson) BeforeRegistration() error {
	zap.L().Info("BeforeRegistration")
	return nil
}

func (lesson *Lesson) AfterRegistration() error {
	zap.L().Info("AfterRegistration")
	return nil
}

func (lesson *Lesson) DoRegistration() error {
	zap.L().Info("StartRegistration")
	if err := lesson.BeforeRegistration(); err != nil {
		return err
	}
	//TODO nick the timeout code from go by example
	//TODO we need to listen to the registration channel some'ow
	time.Sleep(lesson.Registration.Timeout) //TODO this needs to be based on listening for events.
	//TODO TODO
	//TODO some more logic here...
	//TODO start by fleshing out a reasonable main function using the websocket example...
	//TODO TODO we'll need to like do a go run from the run registration and then cancel it after a timeout...
	lesson.Registration.Done = true
	if err := lesson.AfterRegistration(); err != nil {
		return err
	}
	return nil
}

func (lesson *Lesson) AfterLessonStart(planChan chan Plan) error {
	zap.L().Info("AfterLessonStart")
	select {
	case res := <-planChan:
		lesson.Plan = res
	case <-time.After(10 * time.Second): //TODO parameterise
		return errors.New("took to long to receive plan from database")
	}
	return nil
}

func (lesson *Lesson) LessonStart(planId uuid.UUID) error {
	zap.L().Info("LessonStart")
	if planChan, err := lesson.BeforeLessonStart(planId);err != nil {
		return err
	} else {
		if err := lesson.DoRegistration(); err != nil {
			return err
		}
		if err := lesson.AfterLessonStart(planChan); err != nil {
			return err
		}
	}
	return nil
}

