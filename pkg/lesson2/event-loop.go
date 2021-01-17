package lesson2

type EventLoop struct {
	Ctrl Controller
	CurrentRound Round
	ConsoleInChannel <-chan EventLoopEvent
	TheatreInChannel <-chan EventLoopEvent
	CtrlInChannel <-chan EventLoopEvent //TODO rarely used - really just for serious async
	ConsoleOutChannels map[ClientID] chan<- EventLoopEvent
	TheatreOutChannels map[ClientID] chan<- EventLoopEvent
}

//TODO new event loop! - mostly copy this in and set up the handle logic

//TODO this MUST be built up, roughly speaking, FROM the websocket interface.