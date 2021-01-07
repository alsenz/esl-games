package lesson

import "github.com/gofrs/uuid"

//TODO need to break out certain responses here... TODO TODO

type Response struct {
	Question *Question
	//TODO! this is gonna need to be pretty generic - make this work some'ow ;) Would be
	Data interface{}
}

type PlayerResponseEvent struct {
	PlayerToken uuid.UUID
	QuestionLink QuestionLink
	Data interface{}
}


