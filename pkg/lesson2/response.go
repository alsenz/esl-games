package lesson2

type ResponseContentType string
const (
	ContentTypeText		ResponseContentType = "text"
	ContentTypeCanvas	ResponseContentType = "canvas"
	ContentTypeImage	ResponseContentType = "image"
	ContentTypeVideo	ResponseContentType = "video"
	ContentTypeNumber	ResponseContentType = "number"
	ContentTypeEmoji	ResponseContentType = "emoji"
	ContentTypeBoolean	ResponseContentType = "boolean"
	ContentTypeMathjax	ResponseContentType = "mathjax"
	ContentTypeAudio	ResponseContentType = "audio"
	ContentTypeChart	ResponseContentType = "chart"
)
var QuestionContentTypes = [...]ResponseContentType{ContentTypeText, ContentTypeCanvas, ContentTypeImage,
	ContentTypeVideo, ContentTypeNumber, ContentTypeEmoji, ContentTypeBoolean, ContentTypeMathjax, ContentTypeAudio,
	ContentTypeChart}

type Response struct {
	Player Player
	ContentType ResponseContentType
	Data interface{}
	Logic QuestionLogic
}

//TODO need to connect event loop up with websocket channels before understanding how this looks
func (response *Response) Validate() bool {
	//TODO needs to check that data is the right type on the way up
	//TODO TODO
	return false
}
