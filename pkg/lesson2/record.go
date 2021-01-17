package lesson2


type Record struct {
	Player Player
	Round Round
	ScoreCards ScoreCards
	Question *Question //Optional
	Team TeamName
	Response *Response //TODO need to fold this through model
}
