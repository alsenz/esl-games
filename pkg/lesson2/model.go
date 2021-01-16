package lesson2

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"html/template"
)

//TODO if we're allowing for CLIENT-SPECIFIC renderings...
//TODO TODO we should potentially have a ClientID variable and have a kinda
//TODO individual client mode where we render many times for each client...
//TODO TODO makes quite a lot of sense...
type Model struct {
	CurrentRound Round
	// Note that this can naturally be grouped just by inspecting
	// ... the single, grouped upon key
	Frame map[Round]map[TeamName]map[ClientID]Record
}

//TODO implement these and then revisit the circle-back
//TODO TODO

// Splits the model into many sub models each with just one ClientID
func (mdl *Model) ForEachPlayer() []Model {
	//TODO just produce a nice map and then produce a list of models
}

// Splits the model into many sub models each with just one Round
func (mdl *Model) ForEachRound() []Model {
	//TODO just produce a nice map and then produce a list of models
}

func (mdl *Model) ForEachTeam() []Model {
	//TODO
}

func (mdl *Model) SortEntriesByScore(scoreName string) Model {
	//TODO
}

func (mdl *Model) SortEntriesByAccumulatedScore(scoreName string) Model {
	//TODO
}

func (mdl *Model) GetEntriesByPlayerName(playerName string) Model {
	//TODO
}

func (mdl *Model) GetEntriesByClientID(clientID string) Model {
	//TODO
}

func (mdl *Model) GetEntriesByAct(act uint64) Model {
	//TODO
}

func (mdl *Model) GetEntriesByActScene(act uint64, scene uint64) Model {
	//TODO
}

func (mdl *Model) GetEntriesByRound(act uint64, scene uint64, repeat uint64) Model {
	//TODO
}

func (mdl *Model) GetEntriesByTeam(teamName string) Model {
	//TODO
}

func (mdl *Model) GetQuestions() []Question {
	//TODO
}

func (mdl *Model) GetPlayers() []Player {
	//TODO
}

func (mdl *Model) GetScoreCards() []ScoreCards {
	//TODO
}

func (mdl *Model) GetScores(scoreName string) []ScoreCard {
	//TODO
}

func (mdl *Model) GetScoreCardsByPlayerName(playerName string) []ScoreCards {
	//TODO
}

func (mdl *Model) GetQuestionsByPlayerName(playerName string) []Question {
	//TODO
}

//TODO circle back to question after this...
func (mdl *Model) Eval(templateStr string) (string, error) {
	// Just translate the expression into a parseable if expression
	tpl, err := template.New("base").Funcs(sprig.HermeticHtmlFuncMap()).Parse(
		templateStr)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err = tpl.Execute(buf, mdl); err != nil {
		return "", err
	}
	return buf.String(), nil
}

