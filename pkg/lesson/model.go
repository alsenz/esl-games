package lesson

import (
	"bytes"
	"text/template"
)


type Model struct {
	Round Round
	AllRounds map[RoundIdx]Round
	Players []Player
	Teams []Team
	Scores Scores
}

func NewLessonModel() *Model {
	return &Model{*FirstRound(), make(map[RoundIdx]Round), make([]Player, 0),
		make([]Team, 0), NewScores()}
}


func (ctx *Model) Eval(str string) (string, error) {
	if tpl, err := template.New("template").Parse(str); err == nil {
		buf := &bytes.Buffer{}
		if err = tpl.Execute(buf, ctx); err != nil {
			return "", err
		}
		return buf.String(), nil
	} else {
		return "", err
	}
}

func ValidateCtxEval(str string) (error) {
	_, err := template.New("template").Parse(str)
	return err
}