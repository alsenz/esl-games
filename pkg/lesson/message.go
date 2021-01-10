package lesson

import (
	"encoding/json"
)

type Message struct {
	PlayerToken PlayerToken     `json:"playerId,omitempty"`
	Round          RoundIdx        `json:"round,omitempty"` //Helps idempotency and avoids races
	Data           json.RawMessage `json:"data,omitempty"`
}

type ConsoleMessageInType string
const (
	ConsoleSkipMessage ConsoleMessageInType = "skip"
	ConsoleRegisterMessage ConsoleMessageInType = "register"
)

type ConsoleMessageIn struct {
	Message
	Type ConsoleMessageInType	`json:"type"`
	OptOutChan *chan<- ConsoleMessageOut
}

type ConsoleMessageOutType string
const (
	ConsoleRequestInputMessage ConsoleMessageOutType = "request_input"
	ConsoleShowLoadingMessage ConsoleMessageOutType = "show_loading"
	ConsoleShowIdle ConsoleMessageOutType = "show_idle"
	ConsoleEndGameMessage ConsoleMessageOutType = "end_game"
)

type ConsoleMessageOut struct {
	Message
	Type ConsoleMessageOutType	`json:"type"`
}

type TheatreMessageInType string
const (
	TheatreSkipMessage TheatreMessageInType = "skip"
)

type TheatreMessageIn struct {
	Message
	Type TheatreMessageInType `json:"type"`
}

type TheatreMessageOutType string
const (
	TheatreUpdateScreenMessage TheatreMessageOutType = "update_screen" //Render a screen
	TheatreGoToRoundMessage TheatreMessageOutType = "go_to_round"
	TheatreEndGameMessage TheatreMessageOutType = "end_game"
)

type TheatreMessageOut struct {
	Message
	Type TheatreMessageOutType `json:"type"`
}