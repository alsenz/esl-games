package lesson2

//TODO these need to be super simple. Add: player, round, upsert response, score round etc.
//TODO we'll also need to gen new scenes on a rep of round etc. Not quite sure how that'll work.


type ControllerEvent interface {
	Handle(* Controller) error
}

type CtrlAddPlayerEvent struct {
	ClientID string
	PlayerName  string
}

//TOOD add a bunch of events