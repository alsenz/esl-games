package lesson2

import "encoding/json"

//TODO the list should now be VERY, VERY SIMPLE.

//TODO FROM console: register, input-response, request-round
//TODO TO console: input-request, show-loading, show-idle, pause-lesson, end-lesson
//TODO FROM theatre: request-round
//TODO TO theatre: load-slide, show-slide, pause-lesson end-lesson

//TODO we should have the pattern where we have an apply-to(EventLoop *)!
//TODO note these are mostly super super simple.

//TODO i'm really not quite sure how this sits together...

type WebsocketMessage struct {
	Type string
	ClientID ClientID
	Data json.RawMessage
}

func (wsm *WebsocketMessage) send(to chan<- EventLoopEvent) error {
	//TODO this is key as it's a nice way of applying to an even tloop
	return nil
}