package lesson2

//TODO we should allow go-templates in text I think

type Renderer struct {
	slideTemplate string
}

func NewRenderer(slideTemplate string) *Renderer {
	return &Renderer{slideTemplate}
}

func (renderer *Renderer) Render() (string, error) {
	//TODO need to walk and expand the dom using the expansion logic
	//TODO need to use templates on text content so that we can have some power there too
	//TODO TODO use net/http extensively! And basically just walk the dom.
	return "", nil
}


