package lesson

import (
	"database/sql/driver"
	"encoding/json"
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
	//TODO let's get a list of layouts
)
var Layouts = map[string]bool{}

type LayoutFlag string
const (
	FillPage LayoutFlag = "fill page"
	VariationTop LayoutFlag = "top variation"
	VariationMiddle LayoutFlag = "middle variation"
	VariationBottom LayoutFlag = "bottom variation"
	VariationLeft LayoutFlag = "left variation"
	VariationRight LayoutFlag = "right variation"
	//TODO image modes centre crop, fill etc.
)
type LayoutFlags map[string]bool
var LayoutFlagsSet = LayoutFlags{}

//TODO need some way to provide "theming" to the layout elements

type SceneClasses struct {
	//TODO
}

type Scene struct {
	OptionalQuestionFilter *QuestionFilter //May be a template
	Layout LayoutFlags
	Title string
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