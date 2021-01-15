package lesson2


type Model struct {
	CurrentRound Round
	LatestRound Round
	// Note that this can naturally be grouped just by inspecting
	// ... the single, grouped upon key
	Frame map[Round]map[ClientID]Record
}

// Splits the model into many sub models each with just one ClientID
func (mdl *Model) ForEachPlayer() []Model {
	//TODO just produce a nice map and then produce a list of models
}

// Splits the model into many sub models each with just one Round
func (mdl *Model) ForEachRound() []Model {
	//TODO just produce a nice map and then produce a list of models
}

//TODO latest score, current score etc
//TODO score deets on own
//TODO get question, latest question, current question, last question etc.

//TODO get player deets

