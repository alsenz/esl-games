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

// Makes a client model with the same player and round but no records
func (mdl *ClientModel) EmptyCopy() ClientModel {
	return ClientModel{
		make(Model),
		mdl.CurrentPlayer,
		mdl.CurrentRound,
	}
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


// Splits the model into many sub ClientModels each with just one ClientID
func (mdl *ClientModel) ForEachPlayer() []ClientModel {
	grouped := make(map[Player]ClientModel)
	for _, rec := range mdl.Records() {
		val, found := grouped[rec.Player]
		if !found {
			grouped[rec.Player] = mdl.EmptyCopy()
			val = grouped[rec.Player]
		}
		val.AddRecord(rec)
	}
	mdlSlice := make([]ClientModel, 0, len(grouped))
	for _, val := range grouped {
		mdlSlice = append(mdlSlice, val)
	}
	return mdlSlice
}

func (mdl *ClientModel) ForEachPlayerByRank(scoreName ScoreName) []ClientModel {
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

func (mdl *ClientModel) ForEachPlayerByAccumulatedRank(scoreName ScoreName) []ClientModel {
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

//TODO we're gonna need an ability to get scores for teams
//TODO TODO

// Splits the model into many sub models each with just one Round
func (mdl *ClientModel) ForEachRound() []ClientModel {
	grouped := make(map[Round]ClientModel)
	for _, rec := range mdl.Records() {
		val, found := grouped[rec.Round]
		if !found {
			grouped[rec.Round] = mdl.EmptyCopy()
			val = grouped[rec.Round]
		}
		val.AddRecord(rec)
	}
	mdlSlice := make([]ClientModel, 0, len(grouped))
	for _, val := range grouped {
		mdlSlice = append(mdlSlice, val)
	}
	// We sort by round since this is likely to be expected behaviour
	//By convention, everything in the models has a round
	sort.Slice(mdlSlice, func(i int, j int) bool {
		// Just look at the round of the first record
		var iRound, jRound Round
		for round, _ := range mdlSlice[i].Model {
			iRound = round
			break
		}
		for round, _ := range mdlSlice[j].Model {
			jRound = round
			break
		}
		return iRound.LessThan(jRound)
	})
	return mdlSlice
}

func (mdl *ClientModel) ForEachTeam() []ClientModel {
	grouped := make(map[TeamName]ClientModel)
	for _, rec := range mdl.Records() {
		val, found := grouped[rec.Team]
		if !found {
			grouped[rec.Team] = mdl.EmptyCopy()
			val = grouped[rec.Team]
		}
		val.AddRecord(rec)
	}
	mdlSlice := make([]ClientModel, 0, len(grouped))
	for _, val := range grouped {
		mdlSlice = append(mdlSlice, val)
	}
	return mdlSlice
}

//TODO implement these and then revisit the circle-back

func (mdl *ClientModel) GetEntriesByPlayerName(playerName string) ClientModel {
	newMdl := mdl.EmptyCopy()
	for _, rec := range mdl.Records() {
		if rec.Player.Name == playerName {
			newMdl.AddRecord(rec)
		}
	}
	return newMdl
}

func (mdl *ClientModel) GetEntriesByClientID(clientID string) ClientModel {
	newMdl := mdl.EmptyCopy()
	for _, rec := range mdl.Records() {
		if string(rec.Player.ClientID) == clientID {
			newMdl.AddRecord(rec)
		}
	}
	return newMdl
}

func (mdl *ClientModel) GetEntriesByAct(act uint64) ClientModel {
	newMdl := mdl.EmptyCopy()
	for _, rec := range mdl.Records() {
		if rec.Round.Act == act {
			newMdl.AddRecord(rec)
		}
	}
	return newMdl
}

func (mdl *ClientModel) GetEntriesByActScene(act uint64, scene uint64) ClientModel {
	newMdl := mdl.EmptyCopy()
	for _, rec := range mdl.Records() {
		if rec.Round.Act == act && rec.Round.Scene == scene {
			newMdl.AddRecord(rec)
		}
	}
	return newMdl
}

func (mdl *ClientModel) GetEntriesByRound(act uint64, scene uint64, repeat uint64) ClientModel {
	newMdl := mdl.EmptyCopy()
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
		mdl.GetEntriesByRound(0, 0 , 0)
	}
	mdl.GetEntriesByRound(mdl.CurrentRound.Act, mdl.CurrentRound.Scene - 1,
		mdl.CurrentRound.Rep)
}

//Fixme: I'm not sure semantically a Round is a Round as much as a Rep?
func (mdl *ClientModel) GetEntriesByPreviousRep() ClientModel {
	if mdl.CurrentRound.Rep <= 1 { //There are no previous reps...
		return mdl.GetEntriesByRound(0, 0 , 0)
	}
	return mdl.GetEntriesByRound(mdl.CurrentRound.Act, mdl.CurrentRound.Scene,
		mdl.CurrentRound.Rep-1)
}

func (mdl *ClientModel) GetEntriesByCurrentPlayer() ClientModel {
	return mdl.GetEntriesByClientID(string(mdl.CurrentPlayer.ClientID))
}

func (mdl *ClientModel) GetEntriesByLastRoundWithQuestion() ClientModel {
	maxRound := Round{0, 0, 0}
	// This is two-pass: find the max round and then find all the records that are the max
	for _, rec := range mdl.Records() {
		if rec.Question != nil && maxRound.LessThan(rec.Round) {
			maxRound = rec.Round
		}
	}
	return mdl.GetEntriesByRound(maxRound.Act, maxRound.Scene, maxRound.Rep)
}

func (mdl *ClientModel) GetEntriesByTeam(teamName string) ClientModel {
	result := mdl.EmptyCopy()
	for _, rec := range mdl.Records() {
		if string(rec.Team) == teamName {
			result.AddRecord(rec)
		}
	}
	return result
}

func (mdl *ClientModel) GetQuestions() []Question {
	questions := make([]Question, 0, len(mdl.Model))
	for _, rec := range mdl.Records() {
		if rec.Question != nil {
			questions = append(questions, *rec.Question)
		}
	}
	return questions
}

func (mdl *ClientModel) GetPlayers() []Player {
	players := make([]Player, 0, len(mdl.Model))
	for _, rec := range mdl.Records() {
		if rec.Question != nil {
			players = append(players, rec.Player)
		}
	}
	return players
}

func (mdl *ClientModel) GetScoreCards() []ScoreCards {
	scoreCards := make([]ScoreCards, 0, len(mdl.Model))
	for _, rec := range mdl.Records() {
		if rec.Question != nil {
			scoreCards = append(scoreCards, rec.ScoreCards)
		}
	}
	return scoreCards
}

//TODO TODO TODO this still needs to even be thought about, let alone implemented
//TODO we probably need a GetResponses TODO TODO

func (mdl *ClientModel) GetScores(scoreName string) []ScoreCard {
	scores := make([]ScoreCard, 0, len(mdl.Model))
	for _, rec := range mdl.Records() {
		if rec.Question != nil {
			if score, found := rec.ScoreCards[ScoreName(scoreName)]; found {
				scores = append(scores, score)
			}
		}
	}
	return scores
}

func (mdl *Model) Eval(currentPlayer Player, currentRound Round, templateStr string) (string, error) {
	clientMdl := ClientModel{
		*mdl,
		currentPlayer,
		currentRound,
	}
	return clientMdl.Eval(templateStr)
}

func (mdl *ClientModel) Eval(templateStr string) (string, error) {
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

