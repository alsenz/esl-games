package lesson

import (
	"go.uber.org/zap"
	"time"
)

type ReadyForNextRound struct {
	CurrentRound RoundIdx
}

type ShowScreenEvent struct {
	Round RoundIdx
	Scene Scene //TODO I feel like this should probably be "screen"- a rendered version of a scene! TODO TODO - perhaps pick this up next?
	InputRequest map[PlayerToken]Message
}

type EventLoop struct {
	CurrentRound RoundIdx //Current round index as we understand it
	// Channels for sending and receiving messages to/from theatre websockets
	TheatreChannelIn <-chan Message
	TheatreChannelOut chan<- Message
	// Channels for sending and receiving messages to/from console websockets
	ConsoleChannelIn <-chan Message
	ConsoleChannelOut chan<- Message
	// Channels for telling the planner that we can move to another round
	PlannerChannelOut chan<- ReadyForNextRound
	// Channels for telling the controller about state changes
	CtrlRegisterChannelOut chan<- RegistrationEvent
	ControllerChannelOut chan<- PlayerResponseEvent
	ControllerChannelIn <-chan ShowScreenEvent
}

func (loop *EventLoop) BeforeLessonStart(planId uuid.UUID) (chan Plan, error) {
	planChannel := make(chan Plan)
	go loop.LoadPlan(planId, planChannel)
	zap.L().Info("BeforeLessonStart")
	return planChannel, nil
}

func (loop *EventLoop) BeforeRegistration() error {
	zap.L().Info("BeforeRegistration")
	return nil
}

func (loop *EventLoop) AfterRegistration() error {
	zap.L().Info("AfterRegistration")
	return nil
}

func (loop *EventLoop) DoRegistration() error {
	zap.L().Info("StartRegistration")
	if err := loop.BeforeRegistration(); err != nil {
		return err
	}
	//TODO nick the timeout code from go by example
	//TODO we need to listen to the registration channel some'ow
	time.Sleep(loop.Registration.Timeout) //TODO this needs to be based on listening for events.
	//TODO TODO
	//TODO some more logic here...
	//TODO start by fleshing out a reasonable main function using the websocket example...
	//TODO TODO we'll need to like do a go run from the run registration and then cancel it after a timeout...
	loop.Registration.Done = true
	if err := loop.AfterRegistration(); err != nil {
		return err
	}
	return nil
}

func (loop *EventLoop) AfterLessonStart(planChan chan Plan) error {
	zap.L().Info("AfterLessonStart")
	//TODO this needs to go into planner...
	select {
	case res := <-planChan:
		loop.Plan = res
	case <-time.After(10 * time.Second): //TODO parameterise
		return loop.New("took to long to receive plan from database")
	}
	return nil
}

func (loop *EventLoop) LessonStart(planId uuid.UUID) error {
	zap.L().Info("LessonStart")
	if planChan, err := loop.BeforeLessonStart(planId);err != nil {
		return err
	} else {
		if err := loop.DoRegistration(); err != nil {
			return err
		}
		if err := loop.AfterLessonStart(planChan); err != nil {
			return err
		}
	}
	return nil
}

