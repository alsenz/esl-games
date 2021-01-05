package lesson

// A stripped-down version of user that can be baked into templates - not even an email address!
type Player struct {
	Name string
	Avatar []byte //This may or not contain an image...
	ScoreCard ScoreCard
	AllScoreCards map[RoundIdx]ScoreCard
	Responses map[RoundIdx]map[QuestionLink]PlayerResponse
}

type Team struct {
	Name string
	Avatar []byte
	Players []Player
	ScoreCard ScoreCard
	AllScoreCards map[RoundIdx]ScoreCard
}