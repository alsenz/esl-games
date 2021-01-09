package lesson

type Controller struct {
	Planner Planner
	Model Model

	//Channels...
	// Channels for communicating with the event loop
	NextRoundChannelIn <-chan MoveOnEvent
	PlayerResponseIn <-chan PlayerResponseEvent
	PlayerRegisterIn <-chan RegistrationEvent
	ShowScreenChannel chan<- ShowScreenEvent
}

func (c Controller) Run() {
	//TODO basically run handle event a lot!

}
