package lesson

import (
	"bytes"
	"google.golang.org/protobuf/cmd/protoc-gen-go/testdata/proto2"
	"text/template"
)

type QuestionLink string //This might sometimes have readable form, but more likely it won't do and the link is
						//... managed in lesson.go

type Context struct {
	Round Round
	AllRounds map[RoundIdx]Round
	Players []Player
	Teams []Team
	//Note plan has to hold these links and lesson has to maintain them
	LinkedQuestionRounds map[RoundIdx]Round // Back links to (multiple linked question rounds, if they exist
	LinkedQuestions map[QuestionLink]Question
	LinkedResponses map[QuestionLink]PlayerResponses
	Scores Scores
}

func NewContext() *Context {
	//TODO need NewScores... TODO TODO once done this pick up from here...
	return &Context{*FirstRound(), make(map[RoundIdx]Round), make([]Player, 0), make([]Team, 0),
		make(map[RoundIdx]Round), make(map[QuestionLink]Question),
		make(map[QuestionLink]PlayerResponses), NewScores()}
}


func (ctx *Context) Eval(str string) (string, error) {
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