package lesson2

type EventLoopEvent interface {
	HandleFromController(* EventLoop) error
	HandleFromConsole(* EventLoop, string) error //Second argument is clientID
	HandleFromTheatre(* EventLoopEvent) error
}