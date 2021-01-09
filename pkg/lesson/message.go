package lesson

import (
	"encoding/json"
)

type Message struct {
	OptPlayerToken PlayerToken     `json:"playerId,omitempty"`
	Round          RoundIdx        `json:"round,omitempty"` //Helps idempotency and avoids races
	Data           json.RawMessage `json:"data,omitempty"`
}

type ConsoleMessageInType string
const (
	ConsoleSkipMessage ConsoleMessageInType = "skip"
	RegisterMessage ConsoleMessageInType = "register"
)

type ConsoleMessageIn struct {
	Message
	Type ConsoleMessageInType	`json:"type"`
}

type ConsoleMessageOutType string
const (
	//TODO add some real messages here
	DummyMessage ConsoleMessageOutType = "dummy"
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
	ScreenMessage TheatreMessageOutType = "screen" //Render a screen
)

type TheatreMessageOut struct {
	Message
	Type TheatreMessageOutType `json:"type"`
}