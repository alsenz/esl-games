package lesson

import (
	"errors"
	"go.uber.org/zap"
)

//TODO how to build an API over lots of channels?
//TODO how to have an interface based on fields? What would it be? AddPlayer, AddResponse, RequestPreload

type ControllerEvent interface {
	Handle(* Controller) error
}

type CtrlAddPlayerEvent struct {
	PlayerToken ClientID
	PlayerName  string
}

type CtrlMakeTeamsEvent struct {}


type Controller struct {
	Planner Planner
	Model Model
	//Channels...
	// Channels for communicating with the event loop
	EventChannelIn <-chan ControllerEvent
	EventLoopOut chan<- EventLoopEvent
}

// Note: we assume that LoadPlan has already happened
func (c Controller) Run() error {
	if ! c.Planner.PlanIsLoaded() {
		return errors.New("Logic error - attempt to Run() controller before Plan is loaded on Planner")
	}
	for msg := range c.EventChannelIn {
		if err := msg.Handle(&c); err != nil {
			zap.L().Error("Controller loop handle received an error, shutting down controller")
			return err
		}
	}

}

func (event CtrlAddPlayerEvent) Handle(ctrl *Controller) error {
	for idx, existingPlayer := range ctrl.Model.Players {
		if existingPlayer.Token == event.PlayerToken {
			existingPlayer.Name = event.PlayerName
			ctrl.Model.Players[idx] = existingPlayer //TODO check if there's an inplace version here
			return nil
		}
	}
	ctrl.Model.Players = append(ctrl.Model.Players, Player{
		event.PlayerName,
		nil,
		event.PlayerToken,
		make(ScoreCard),
		make(map[RoundIdx]ScoreCard),
		make(map[RoundIdx]Response),
	})
	return nil
}

func (event CtrlMakeTeamsEvent) Handle(ctrl *Controller) error {
	//TODO we need to put people in teams here
	return nil
}

//TODO other controller events here? Build out a nice API?