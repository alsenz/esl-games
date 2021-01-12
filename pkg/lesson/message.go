package lesson

import (
	"encoding/json"
	"errors"
)

//TODO we can simplify these a bit. Call them WebsocketMessage instead.
//TODO and we don't need the different between the In and the Out, although a console message vs theatre message a good idea.

type EventLoopTransformable interface {
	MakeEventLoopEvent() EventLoopEvent
}

type WebsocketMessage struct {
	PlayerToken ClientID        `json:"playerId,omitempty"`
	Round       RoundIdx        `json:"round,omitempty"` //Helps idempotency and avoids races
	Data        json.RawMessage `json:"data,omitempty"`
}

type ConsoleMessageType string
const (
	ConsoleSkipMessage ConsoleMessageType = "skip_round"
	ConsoleRegisterMessage ConsoleMessageType = "register"
	ConsoleRequestInputMessage ConsoleMessageType = "request_input"
	ConsoleShowLoadingMessage ConsoleMessageType = "show_loading"
	ConsoleShowIdle ConsoleMessageType = "show_idle"
	ConsoleEndGameMessage ConsoleMessageType = "end_game"
)

type ConsoleMessage struct {
	WebsocketMessage
	Type ConsoleMessageType	`json:"type"`
	ClientOutChannel *chan<- ConsoleMessage //Optional - the output chanel associated with this console message for this client
}

func (msg ConsoleMessage) MakeEventLoopEvent() (EventLoopEvent, error) {
	switch msg.Type {
	case ConsoleRegisterMessage:
		if msg.ClientOutChannel == nil {
			return nil, errors.New("Unable to make EventLoopEvent for ConsoleMessage of type " + string(msg.Type) + "... no out channel provided")
		}
		return RegisterEvent{msg.PlayerToken, *msg.ClientOutChannel}, nil
	default:
		//TODO error log... for the others.
	}
}

type TheatreMessageType string
const (
	TheatreSkipMessage TheatreMessageType = "skip_round"
	TheatreUpdateScreenMessage TheatreMessageType = "update_screen" //Render a screen
	TheatreGoToRoundMessage TheatreMessageType = "go_to_round"
	TheatreEndGameMessage TheatreMessageType = "end_game"
)

type TheatreMessage struct {
	WebsocketMessage
	Type TheatreMessageType `json:"type"`
}

//TODO we need a MakeEventLoopEvent() for TheatreMessage.... TOOD TODO