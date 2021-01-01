package lesson

import uuid "github.com/satori/go.uuid"

//TODO this is
type Lesson struct {
	Plan Plan
	QABindings map[uuid.UUID] uuid.UUID			//Maps Question.ID's to Response.ID's. If not in the map, we haven't gone and go it yet.
	ResolvedQuestions map[uuid.UUID] Question	//Maps resolved questions for each round
	//TODO is this right? Because a question could come up any time...? TODO TODO And there might be multiple questions in a round
	ResolvedQuestionRounds map[Round] []Question  //One-many map from rounds to questions that are bound to it
	RoundResolvedQuestions map[uuid.UUID] *Round //Many-one map from question ids to rounds that mention them
	//TODO we probably wanna map responses (by id) to rounds they were bound to too - and only once...
}

func NewLesson() * Lesson {
	lesson := &Lesson{}
	return lesson
}

//TODO this needs a receiving channel for client events and an outgoing channel for clients (broadcast) and shows (broadcast)
// Go routine lesson worker
func RunLesson() bool {
	//Go-routine state
	lesson := NewLesson()
	_ = lesson.Plan //TODO work this through
	return false
}
