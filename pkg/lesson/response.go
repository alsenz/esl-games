package lesson

//TODO need to break out certain responses here... TODO TODO

type PlayerResponse struct {
	//TODO! this is gonna need to be pretty generic
	Response interface{}
}

type PlayerResponses map[Player]PlayerResponse
