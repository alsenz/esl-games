package lesson

import (
	"database/sql/driver"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
)

type Layout string
const (
	LayoutTitle Layout = "title"
	LayoutSectionTitle Layout = "section title"
	LayoutTitleContent Layout = "content with title (optional)"
	LayoutHeroTop Layout = "hero"
	LayoutTwoColumnContent Layout = "two column content"
	LayoutThreeColumnContent Layout = "three column content"
	LayoutQuadrants Layout = "four quadrants content"
	LayoutContentAndImage = "content with image"
	//TODO extend this list a bit further
)
var Layouts = map[Layout]bool{LayoutTitle: true, LayoutSectionTitle: true, LayoutTitleContent: true, LayoutHeroTop: true,
	LayoutTwoColumnContent: true, LayoutThreeColumnContent: true, LayoutQuadrants: true, LayoutContentAndImage: true}

type LayoutFlag string
const (
	FillPage LayoutFlag = "fill page"
	VariationTop LayoutFlag = "top variation"
	VariationMiddle LayoutFlag = "middle variation"
	VariationBottom LayoutFlag = "bottom variation"
	VariationLeft LayoutFlag = "left variation"
	VariationRight LayoutFlag = "right variation"
	VariationImageStretch LayoutFlag = "image stretch"
	VariationImageCover LayoutFlag = "image cover"
	VariationImageContain LayoutFlag = "image contain"
	VariationImageRepeat LayoutFlag = "image repeat"
)
type LayoutFlagsStruct map[LayoutFlag]bool
var LayoutFlags = LayoutFlagsStruct{FillPage: true, VariationTop: true, VariationMiddle: true, VariationBottom: true,
	VariationLeft: true, VariationRight: true, VariationImageStretch: true, VariationImageCover: true,
	VariationImageContain: true, VariationImageRepeat: true}

type View struct {
	Layout Layout	`json:"layout"`
	LayoutFlags LayoutFlagsStruct `json:"layoutFlags,omitempty"`
	Title string				`json:"title"`
	Byline string				`json:"byline,omitempty"`
	Header string				`json:"header,omitempty"`
	Footer string				`json:"footer,omitempty"`
	Image uuid.UUID				`json:"image,omitempty"`
	Caption string				`json:"caption,omitempty"`
	Content []string			`json:"content,omitempty"`
	Logo uuid.UUID				`json:"logo,omitempty"`
	Gallery []uuid.UUID			`json:"gallery,omitempty"`
	Classes	[]string			`json:"classes,omitempty"` //This is our theming hook
}

type Scene struct {
	Template View                  `json:"template"`
	OptQuestionFilter *QuestionSet `json:"questionFilter,omitempty"`
}
type Scenes []Scene

func (scs *Scenes) Value() (driver.Value, error) {
	if raw, err := json.Marshal(scs); err != nil {
		return nil, err
	} else {
		return datatypes.JSON(raw).Value()
	}
}

func (scs *Scenes) Scan(src interface{}) error {
	jsn := &datatypes.JSON{}
	if err := jsn.Scan(src); err != nil {
		return err
	}
	return json.Unmarshal(*jsn, scs)
}

type RenderedView struct {
	View
}

func (view View) Render(mdl *Model) (* RenderedView, error) {
	var err error
	var title, byline, header, footer, caption string
	if title, err = mdl.Eval(view.Title); err != nil {
		return nil, err
	}
	if byline, err = mdl.Eval(view.Byline); err != nil {
		return nil, err
	}
	if header, err = mdl.Eval(view.Header); err != nil {
		return nil, err
	}
	if footer, err = mdl.Eval(view.Footer); err != nil {
		return nil, err
	}
	if caption, err = mdl.Eval(view.Caption); err != nil {
		return nil, err
	}
	content := make([]string, len(view.Content))
	for i, val := range content {
		if content[i], err = mdl.Eval(val); err != nil {
			return nil, err
		}
	}
	return &RenderedView{
		View{
			view.Layout,
			view.LayoutFlags,
			title,
			byline,
			header,
			footer,
			view.Image, //TODO find some way to make this templateable
			caption,
			content,
			view.Logo, //TODO find some way to make templatable
			view.Gallery, //TODO ditto
			view.Classes,
		},
	}, nil
}