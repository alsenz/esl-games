package lesson

import uuid "github.com/satori/go.uuid"

type PlayerToken uuid.UUID

type RegistrationEvent struct {
	Name string
	Avatar []byte
	Token PlayerToken //Generated unique to securely disambiguate players //TODO - align with JWT //TODO (probably in the handler itself)
}

// A stripped-down version of user that can be baked into templates - not even an email address!
type Player struct {
	RegistrationEvent
	ScoreCard ScoreCard
	AllScoreCards map[RoundIdx]ScoreCard
	Responses map[RoundIdx]Response //Links the players response to the question for them generated at that round
}

type Team struct {
	Name string
	Avatar []byte
	Players []Player
	ScoreCard ScoreCard
	AllScoreCards map[RoundIdx]ScoreCard
}