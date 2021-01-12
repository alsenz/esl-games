package lesson

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"time"
)


type Lesson struct {
	EventLoop			EventLoop
	Model 				Model
	Register      		Register
	Controller			Controller

	//Channels (to be closed)

}

//TODO rework this
func NewLesson(register Register, inactivityTimeout time.Duration) * Lesson {
	/*	CurrentRound        RoundIdx //Current round index as we understand it. A 0, 0, 0 round = registration
		InactivityTimeout   time.Duration
		LeaderPlayerToken   ClientID
		RegistrationTimeout time.Duration
		// Channels for sending and receiving messages to/from theatre websockets.
		// Go routines handling websockets must convert their respective events to EventLoopEvents
		TheatreChannelIn <-chan EventLoopEvent
		TheatreChannelOuts []chan<- EventLoopEvent
		// Channels for sending and receiving messages to/from console websockets
		// Go routines handling websockets must convert their respective events to EventLoopEvents
		ConsoleChannelIn <-chan EventLoopEvent
		ConsoleChannelOuts map[ClientID]chan<- EventLoopEvent
		// Channels for telling the controller about state changes
		ControllerChannelOut chan<- ControllerEvent
		EventLoopChannelIn <-chan EventLoopEvent*/
	theatreChannelIn := make(chan EventLoopEvent, 16)
	theatreChannelOuts := make([]chan EventLoopEvent, 0)
	consoleChannelIn := make(chan EventLoopEvent, 16)
	//TODO various channels to create
	return &Lesson{
		EventLoop{nil,
			inactivityTimeout,
			nil,
			register.Timeout,
			theatreChannelIn,
			theatreChannelOuts,
			consoleChannelIn,
			
		}, //TODO we need to make some channels here, mostly... TODO create a New for this
		NewLessonModel(),
		register,
		Controller{}, //TODO create a New for this!
	}
}

func (lesson *Lesson) Close() error {
	//TODO need to close the various channels..
}

func (lesson *Lesson) Run(planId uuid.UUID) error {
	//Do this up-front to ensure that there are no funky races
	if err := lesson.Controller.Planner.LoadPlan(); err != nil {
		zap.L().Error("Unable to load plan due to " + err.Error())
		return err
	}
	go lesson.Controller.Run()
	go lesson.EventLoop.Run()
	return nil //TODO - how do we fish out errors from various go routines? TODO?
}


//TODO move the play handler in here so it can be a receiver on lesson