package lesson2

// A checkpoint contains summary information about the state of the game currently
// It is used in eval-ing act repetitions and question sets
type Checkpoint struct {
	MaxAccumulatedScores map[ScoreName]float64
	SumAccumulatedScores map[ScoreName]float64
	Round Round
	//TODO this will need to be made richer using things like previous answers
	//TODO to previous questions for things like "chose a category" situations
	//TODO TODO
}
