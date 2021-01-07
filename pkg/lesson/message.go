package lesson

import (
	"encoding/json"
	"github.com/gofrs/uuid"
)

type MessageType string
const (
	SkipMessage MessageType = "Skip"
	ScreenMessage MessageType = "Screen" //Render a screen
)
var TheatreMessages = map[MessageType]bool{SkipMessage: true, ScreenMessage: true}
var ConsoleMessages = map[MessageType]bool{}

type Message struct {
	OptPlayerID uuid.UUID	`json:"playerId,omitempty"`
	Type MessageType		`json:"type"`
	Round RoundIdx			`json:"round,omitempty"` //Helps idempotency and avoids races
	Data json.RawMessage	`json:"data,omitempty"`
}
