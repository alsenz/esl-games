package lesson2


type Record struct {
	Player Player //TODO build these out and keep them simple!
	Round Round //TODO
	ScoreCards ScoreCards //TODO - needs to contain score (on round), player ranking (on round), team ranking (on round), total ranking, total score
	Question *QuestionFoo //Optional //TODO
	Team TeamName
}
