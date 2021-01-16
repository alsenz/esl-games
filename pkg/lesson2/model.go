package lesson2

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"html/template"
	"sort"
)

//TODO unit tests deffo necessary

// Note that this can naturally be grouped just by inspecting
// ... the single, grouped upon key
type Model map[Round]map[TeamName]map[Player]Record

// A small extension specific to a particular client at a particular round
type ClientModel struct {
	Model
	CurrentPlayer Player
	CurrentRound Round
}

func (mdl *Model) AddRecord(rec Record) {
	teamMap, found := (*mdl)[rec.Round]
	if !found {
		(*mdl)[rec.Round] = make(map[TeamName]map[Player]Record)
		teamMap = (*mdl)[rec.Round]
	}
	playerMap, found := teamMap[rec.Team]
	if !found {
		teamMap[rec.Team] = make(map[Player]Record)
		playerMap = teamMap[rec.Team]
	}
	playerMap[rec.Player] = rec
}

//TODO: Should have something for updating records according to ctrler

func (mdl *Model) Records() []Record {
	records := make([]Record, 0, len(*mdl))
	//Fixme: This should have better capacity reservation
	for _, teamMap := range *mdl {
		for _, playerMap := range teamMap {
			for _, rec := range playerMap {
				records = append(records, rec)
			}
		}
	}
	return records
}


// Splits the model into many sub models each with just one ClientID
func (mdl *Model) ForEachPlayer() []Model {
	grouped := make(map[Player]Model)
	for _, rec := range mdl.Records() {
		val, found := grouped[rec.Player]
		if !found {
			grouped[rec.Player] = make(Model)
			val = grouped[rec.Player]
		}
		val.AddRecord(rec)
	}
	mdlSlice := make([]Model, 0, len(grouped))
	for _, val := range grouped {
		mdlSlice = append(mdlSlice, val)
	}
	return mdlSlice
}

func (mdl *Model) ForEachPlayerByRank(scoreName ScoreName) []Model {
	models := mdl.ForEachPlayer()
	// Sorts from smallest ranking to largest
	sort.Slice(models, func(i, j int) bool {
		// If there's no ranking yet, then this is "very big"
		iScores := models[i].GetScores(string(scoreName))
		if len(iScores) == 0 {
			return false
		}
		jScores := models[j].GetScores(string(scoreName))
		if len(jScores) == 0 {
			return true
		}
		return iScores[0].Ranking < jScores[0].Ranking
	})
	return models
}

func (mdl *Model) ForEachPlayerByAccumulatedRank(scoreName ScoreName) []Model {
	models := mdl.ForEachPlayer()
	// Sorts from smallest ranking to largest
	sort.Slice(models, func(i, j int) bool {
		// If there's no ranking yet, then this is "very big"
		iScores := models[i].GetScores(string(scoreName))
		if len(iScores) == 0 {
			return false
		}
		jScores := models[j].GetScores(string(scoreName))
		if len(jScores) == 0 {
			return true
		}
		return iScores[0].AccumulatedRanking < jScores[0].AccumulatedRanking
	})
	return models
}

// Splits the model into many sub models each with just one Round
func (mdl *Model) ForEachRound() []Model {
	grouped := make(map[Round]Model)
	for _, rec := range mdl.Records() {
		val, found := grouped[rec.Round]
		if !found {
			grouped[rec.Round] = make(Model)
			val = grouped[rec.Round]
		}
		val.AddRecord(rec)
	}
	mdlSlice := make([]Model, 0, len(grouped))
	for _, val := range grouped {
		mdlSlice = append(mdlSlice, val)
	}
	return mdlSlice
}

func (mdl *Model) ForEachTeam() []Model {
	grouped := make(map[TeamName]Model)
	for _, rec := range mdl.Records() {
		val, found := grouped[rec.Team]
		if !found {
			grouped[rec.Team] = make(Model)
			val = grouped[rec.Team]
		}
		val.AddRecord(rec)
	}
	mdlSlice := make([]Model, 0, len(grouped))
	for _, val := range grouped {
		mdlSlice = append(mdlSlice, val)
	}
	return mdlSlice
}

//TODO implement these and then revisit the circle-back

func (mdl *Model) GetEntriesByPlayerName(playerName string) Model {
	newMdl := make(Model, len(*mdl))
	for _, rec := range mdl.Records() {
		if rec.Player.Name == playerName {
			newMdl.AddRecord(rec)
		}
	}
	return newMdl
}

func (mdl *Model) GetEntriesByClientID(clientID string) Model {
	newMdl := make(Model, len(*mdl))
	for _, rec := range mdl.Records() {
		if string(rec.Player.ClientID) == clientID {
			newMdl.AddRecord(rec)
		}
	}
	return newMdl
}

func (mdl *Model) GetEntriesByAct(act uint64) Model {
	newMdl := make(Model, len(*mdl))
	for _, rec := range mdl.Records() {
		if rec.Round.Act == act {
			newMdl.AddRecord(rec)
		}
	}
	return newMdl
}

func (mdl *Model) GetEntriesByActScene(act uint64, scene uint64) Model {
	newMdl := make(Model, len(*mdl))
	for _, rec := range mdl.Records() {
		if rec.Round.Act == act && rec.Round.Scene == scene {
			newMdl.AddRecord(rec)
		}
	}
	return newMdl
}

func (mdl *Model) GetEntriesByRound(act uint64, scene uint64, repeat uint64) Model {
	newMdl := make(Model, len(*mdl))
	for _, rec := range mdl.Records() {
		if rec.Round.Act == act && rec.Round.Scene == scene && rec.Round.Rep == repeat {
			newMdl.AddRecord(rec)
		}
	}
	return newMdl
}

// If there is no previous scene in the act
func (mdl *ClientModel) GetEntriesByPreviousScene() ClientModel {
	if mdl.CurrentRound.Scene <= 1 { // Cannot have a previous scene in the act
		return ClientModel{
			mdl.GetEntriesByRound(0, 0 , 0),
			mdl.CurrentPlayer,
			Round{0, 0, 0},
		}
	}
	return ClientModel{
		mdl.GetEntriesByRound(mdl.CurrentRound.Act, mdl.CurrentRound.Scene - 1,
			mdl.CurrentRound.Rep),
		mdl.CurrentPlayer,
		mdl.CurrentRound,
	}
}

//Fixme: I'm not sure semantically a Round is a Round as much as a Rep?
func (mdl *ClientModel) GetEntriesByPreviousRep() ClientModel {
	if mdl.CurrentRound.Rep <= 1 { //There are no previous reps...
		return ClientModel{
			mdl.GetEntriesByRound(0, 0 , 0),
			mdl.CurrentPlayer,
			Round{0, 0, 0},
		}
	}
	return ClientModel{
		mdl.GetEntriesByRound(mdl.CurrentRound.Act, mdl.CurrentRound.Scene,
			mdl.CurrentRound.Rep-1),
		mdl.CurrentPlayer,
		mdl.CurrentRound,
	}
}

func (mdl *ClientModel) GetEntriesByCurrentPlayer() ClientModel {
	return ClientModel{
		mdl.GetEntriesByClientID(string(mdl.CurrentPlayer.ClientID)),
		mdl.CurrentPlayer,
		mdl.CurrentRound,
	}
}

func (mdl *ClientModel) GetEntriesByLastRoundWithQuestion() ClientModel {
	maxRound := Round{0, 0, 0}
	// This is two-pass: find the max round and then find all the records that are the max
	for _, rec := range mdl.Records() {
		if rec.Question != nil && maxRound.LessThan(rec.Round) {
			maxRound = rec.Round
		}
	}
	return ClientModel{
		mdl.GetEntriesByRound(maxRound.Act, maxRound.Scene, maxRound.Rep),
		mdl.CurrentPlayer,
		mdl.CurrentRound}
}

func (mdl *Model) GetEntriesByTeam(teamName string) Model {
	//TODO this is next. TODO TODO
	//TODO pick this up after coffee...
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

//TODO we probably need a GetResponses TODO TODO

func (mdl *Model) GetScores(scoreName string) []ScoreCard {
	//TODO
}

func (mdl *Model) GetScoreCardsByPlayerName(playerName string) []ScoreCards {
	//TODO
}

func (mdl *Model) GetQuestionsByPlayerName(playerName string) []Question {
	//TODO
}


//TODO circle back to the renderer after this! - make a note to circle back to question!
//TODO circle back to Question after this... look for the Eval!
func (mdl *Model) Eval(currentPlayer Player, currentRound Round, templateStr string) (string, error) {
	clientMdl := ClientModel{
		*mdl,
		currentPlayer,
		currentRound,
	}
	// Just translate the expression into a parseable if expression
	tpl, err := template.New("base").Funcs(sprig.HermeticHtmlFuncMap()).Parse(
		templateStr)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err = tpl.Execute(buf, clientMdl); err != nil {
		return "", err
	}
	return buf.String(), nil
}

