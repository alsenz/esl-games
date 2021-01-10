package lesson

import (
	"encoding/json"
	"go.uber.org/zap"
	"time"
)

type MoveOnEvent struct {
	CurrentRound RoundIdx
}

type ShowScreenEvent struct {
	Round RoundIdx
	View RenderedView
	InputRequest map[PlayerToken]ConsoleMessageOut
}

//TODO lesson needs to construct, and make the channels.
type EventLoop struct {
	CurrentRound        RoundIdx //Current round index as we understand it
	InactivityTimeout	time.Duration
	LeaderPlayerToken   PlayerToken
	RegistrationTimeout time.Duration
	// Channels for sending and receiving messages to/from theatre websockets
	TheatreChannelIn <-chan TheatreMessageIn
	TheatreChannelOuts []chan<- TheatreMessageOut
	// Channels for sending and receiving messages to/from console websockets
	ConsoleChannelIn <-chan ConsoleMessageIn
	ConsoleChannelOuts map[PlayerToken]chan<- ConsoleMessageOut
	// Channels for telling the planner that we can move to another round
	PlannerChannelOut chan<- MoveOnEvent
	// Channels for telling the controller about state changes
	CtrlRegisterChannelOut chan<- RegistrationEvent
	ControllerChannelOut chan<- PlayerResponseEvent
	ControllerChannelIn <-chan ShowScreenEvent
}

func (loop *EventLoop) BeforeLessonStart() error {
	zap.L().Info("BeforeLessonStart")
	return nil
}

func (loop *EventLoop) BeforeRegistration() error {
	zap.L().Info("BeforeRegistration")
	return nil
}

func (loop *EventLoop) AfterRegistration() error {
	zap.L().Info("AfterRegistration")
	return nil
}

// Console Messages sometimes come with new output channels to map players to
func (loop *EventLoop) updatePlayerOut(playerToken PlayerToken, playerOut chan<- ConsoleMessageOut) {
	loop.ConsoleChannelOuts[playerToken] = playerOut
}

func (loop *EventLoop) DoRegistration() error {
	zap.L().Info("StartRegistration")
	if err := loop.BeforeRegistration(); err != nil {
		return err
	}
	Loop:
		for {
		select {
		case msg := <-loop.ConsoleChannelIn:
			if msg.OptOutChan != nil {
				loop.updatePlayerOut(msg.PlayerToken, *msg.OptOutChan)
			}
			switch msg.Type {
			case ConsoleSkipMessage: //Only the leader can skip the registration
				if msg.PlayerToken == loop.LeaderPlayerToken {
					break Loop
				}
			case ConsoleRegisterMessage:
				//TODO the first step is to ensure that the output channel is right
				//TODO the next step is to add the player to the controller
			default:
				zap.L().Warn("ConsoleInMessageType " + string(msg.Type) + " is not accepted at this stage.")
			}
		case msg := <-loop.TheatreChannelIn:
			switch msg.Type {
			case TheatreSkipMessage: //Only the leader can skip the registration
				if msg.OptPlayerToken == loop.LeaderPlayerToken {
					break Loop
				}
			default:
				zap.L().Warn("TheatreInMessageType " + string(msg.Type) + " is not accepted at this stage.")
			}
		case <-time.After(loop.RegistrationTimeout):
			break Loop
		}
	}
	if err := loop.AfterRegistration(); err != nil {
		return err
	}
	return nil
}

func (loop *EventLoop) AfterLessonStart() error {
	zap.L().Info("AfterLessonStart")
	return nil
}

func (loop *EventLoop) LessonStart() error {
	zap.L().Info("LessonStart")
	if err := loop.BeforeLessonStart(); err != nil {
		return err
	}
	if err := loop.DoRegistration(); err != nil {
			return err
}
	if err := loop.AfterLessonStart(); err != nil {
		return err
	}
	return nil
}

func (loop *EventLoop) BeforeLessonEnd() error {
	zap.L().Info("BeforeLessonEnd")
	return nil
}

func (loop *EventLoop) AfterLessonEnd() error {
	zap.L().Info("AfterLessonEnd")
	return nil
}

func (loop *EventLoop) LessonEnd() error {
	zap.L().Info("LessonEnd")
	if err := loop.BeforeLessonEnd(); err != nil {
		return err
	}
	//TODO close down some clients... TODO TODO possibly send a load of stuff - ideally really send a load of end game messages
	//TODO how to SEND to a channel? TODO TODO
	//TODO so far so ok...
	//TODO but how to BROADCAST to channel? Apparently this is impossible... TODO so we're going to need a list
	for consoleChannel := range loop.ConsoleChannelOuts {
		consoleChannel <- ConsoleMessageOut{Message{
			nil, //This is for all players
			loop.CurrentRound,
			json.RawMessage{}, //No extra info required
		}, ConsoleEndGameMessage}
	}
	if err := loop.AfterLessonEnd(); err != nil {
		return err
	}
	return nil
}

func (loop *EventLoop) BeforeHandleConsoleEvent(msg ConsoleMessageIn) error {
	zap.L().Debug("BeforeHandleConsoleEvent")
	return nil
}

func (loop *EventLoop) AfterHandleConsoleEvent(msg ConsoleMessageIn) error {
	zap.L().Debug("AfterHandleConsoleEvent")
	return nil
}

func (loop *EventLoop) HandleConsoleEvent(msg ConsoleMessageIn) error {
	zap.L().Debug("HandleConsoleEvent: " + string(msg.Type))
	if err := loop.BeforeHandleConsoleEvent(msg); err != nil { return err }
	//TODO some key logic here
	switch msg.Type {
		case ConsoleSkipMessage:
			//TODO we wanna send a next round message if the leader
			break
		default:
			zap.L().Warn("Unexpected ConsoleMessageType: " + string(msg.Type))
	}
	if err := loop.AfterHandleConsoleEvent(msg); err != nil { return err }
	return nil
}

func (loop *EventLoop) BeforeHandleTheatreEvent(msg TheatreMessageIn) error {
	zap.L().Debug("BeforeHandleTheatreEvent")
	return nil
}

func (loop *EventLoop) AfterHandleTheatreEvent(msg TheatreMessageIn) error {
	zap.L().Debug("AfterHandleTheatreEvent")
	return nil
}

func (loop *EventLoop) HandleTheatreEvent(msg TheatreMessageIn) error {
	zap.L().Debug("HandleTheatreEvent: " + string(msg.Type))
	if err := loop.BeforeHandleTheatreEvent(msg); err != nil { return err }
	//TODO some key logic here

	if err := loop.AfterHandleTheatreEvent(msg); err != nil { return err }
	return nil
}

func (loop *EventLoop) BeforeHandleScreenEvent(event ShowScreenEvent) error {
	zap.L().Debug("BeforeHandleScreenEvent")
	return nil
}

func (loop *EventLoop) AfterHandleScreenEvent(event ShowScreenEvent) error {
	zap.L().Debug("AfterHandleScreenEvent")
	return nil
}

//TODO I think the idea was we had a preload kinda idea
//TODO and then we had a go to round event

//TODO and we have a preload(N) on the controller, which should just return "what it can", even if incomplete

//TODO this probably requires returning a new round (or something similar) so we can trigger the "new" round type stuff
//TODO - for now let's keep it simple...
func (loop *EventLoop) HandleScreenEvent(event ShowScreenEvent) error {
	zap.L().Debug("HandleScreenEvent")
	if err := loop.BeforeHandleScreenEvent(event); err != nil { return err }
	//TODO some key logic here - let's just review what's in a ShowScreenEvent and that it's good to go...

	/*Round RoundIdx
	View RenderedView
	InputRequest map[PlayerToken]ConsoleMessageOut*/
	//TODO answers will update... there may be a preloading chance... working out when we need to is gonna be hard
	//TODO better to just send to all theatres the screen to be ready to load
	//TODO then send the message to consoles
	//TODO then send the navigation to the theatres probably about the order of affairs.

	//TODO a little bit of logic to think about, especially with preloading in mind... how to de-couple?
	//TODO how to tell when a slide CAN be preloaded? Perhaps we always send preload?

	if err := loop.AfterHandleScreenEvent(event); err != nil { return err }
	return nil
}

func (loop *EventLoop) HandleEvents() error {
	//TODO we need some natural break conditions here too
	Loop:
		for {
			select {
				case msg := <-loop.ConsoleChannelIn:
					if err := loop.HandleConsoleEvent(msg); err!= nil { return err }
				case msg := <-loop.TheatreChannelIn:
					if err := loop.HandleTheatreEvent(msg); err!= nil { return err }
				case event := <-loop.ControllerChannelIn:
					if err := loop.HandleScreenEvent(event); err!= nil { return err }
				case <-time.After(loop.InactivityTimeout):
					zap.L().Fatal("Main event loop is breaking out due to inactivity timeout")
					break Loop
			}
		}
	return nil
}

func (loop *EventLoop) Run() error {
	zap.L().Info("Running Event Loop")
	if err := loop.LessonStart(); err != nil { return err }
	if err := loop.HandleEvents(); err != nil { return err }
	if err := loop.LessonEnd(); err != nil { return err }
	return nil
}