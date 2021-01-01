package lesson

import "github.com/alsenz/esl-games/pkg/account"

//TOOD Scene feels like it should be embedded...
type Scene struct {
	account.UserObject
	OptionalQuestion *Question //May be a template
	OptionalQuestionFilter *QuestionFilter //May be a template
}

func GetResolvable(scene *Scene) *Resolvable {
	// Quesiton filters take precedence
	if scene.OptionalQuestionFilter != nil {
		return scene.OptionalQuestionFilter
	}
	return nil
}

type Act struct {
	account.UserObject
	Scenes []Scene
}

type Plan struct {
	account.UserObject
	Name string
	Description string
	Acts []Act //TODO one-many
}

//TODO designing this out in gorm is gonna be a biggie...!
