package lesson

import uuid "github.com/satori/go.uuid"

//TODO this to probably become ClientID!
type ClientID uuid.UUID

// A stripped-down version of user that can be baked into templates - not even an email address!
type Player struct {
	Name          string
	Avatar        []byte
	Token         ClientID //Generated unique to securely disambiguate players //TODO - align with JWT //TODO (probably in the handler itself)
	ScoreCard     ScoreCard
	AllScoreCards map[RoundIdx]ScoreCard
	Responses     map[RoundIdx]Response //Links the players response to the question for them generated at that round
}

type Team struct {
	Name string
	Avatar []byte
	Players []Player
	ScoreCard ScoreCard
	AllScoreCards map[RoundIdx]ScoreCard
}