package lesson

import (
	"errors"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"time"
)

//TODO lesson's gonna need a "status" API to say what auth is required - the web app can then go <loading/><auth/><noauth/> type thing.
//TODO TODO

type Client struct {
	Send chan Message
}

//TODO this mostly needs to go away
type Lesson struct {
	//TODO this to sit on a planner.
	Planner             Planner
	EventLoop			EventLoop

	//TODO not sure where this is to go
	Registration      Register
	ResolvedQuestions map[RoundIdx]map[QuestionLink]Question //Maps resolved questions asked at each round

	//TODO this to sit on the controller
	Model       *Model //Used for state keeping and evaluation of strings


	//TODO may not need this
	Clients           map[string]Client //Clients mapped by player name (//TODO or client id?)
}

func NewLesson(register Register) * Lesson {
	return &Lesson{nil, register, make(map[RoundIdx]map[QuestionLink]Question),
		NewLessonModel(), make(chan Message, 8), make(map[string]Client)}
}

func (lesson *Lesson) Run(planId uuid.UUID) error {
	go lesson.EventLoop.Run()
	//TODO run the other go routings
	return nil //TODO - how do we fish out errors from various go routines?
}

