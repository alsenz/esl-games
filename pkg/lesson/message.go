package lesson

import (
	"encoding/json"
)

//TODO we can simplify these a bit. Call them WebsocketMessage instead.
//TODO and we don't need the different between the In and the Out, although a console message vs theatre message a good idea.

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
}

//TODO
//OptOutChan *chan<- ConsoleMessageOut - the register event will have a callback channel

type TheatreMessageInType string
const (
	TheatreSkipMessage TheatreMessageInType = "skip"
)

type TheatreMessageIn struct {
	WebsocketMessage
	Type TheatreMessageInType `json:"type"`
}

type TheatreMessageOutType string
const (
	TheatreUpdateScreenMessage TheatreMessageOutType = "update_screen" //Render a screen
	TheatreGoToRoundMessage TheatreMessageOutType = "go_to_round"
	TheatreEndGameMessage TheatreMessageOutType = "end_game"
)

type TheatreMessageOut struct {
	WebsocketMessage
	Type TheatreMessageOutType `json:"type"`
}