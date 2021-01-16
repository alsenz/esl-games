package lesson2

type Scene struct {
	Template string
	ScoresActive []ScoreName
	QuestionSet *QuestionSet //TODO add these filters across.
}
