package lesson2

type EventLoopEvent interface {
	HandleFromConsole(* EventLoop, string) error //Second argument is clientID
	HandleFromTheatre(* EventLoopEvent) error
}