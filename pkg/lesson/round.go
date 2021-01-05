package lesson

import (
	"bytes"
	"time"
	"text/template"
)

type RoundStats struct {
	QuestionNumber uint 		// total number of scenes that were also questions
	StartTime time.Time
	Time time.Time				// the actual time
}

func NewRoundStats() RoundStats {
	now := time.Now()
	return RoundStats{0, now, now}
}

func (rs RoundStats) RefreshRoundStats() RoundStats {
	//TODO just refresh the Time to be now and return a copy.
}

func (rs *RoundStats) PlayTime() time.Duration {
	return time.Since(rs.StartTime) //TODO review the logic on this a bit
}

type RoundIdx struct {
	ActNumber uint				// Number left-to-right - almost like an Act
	SceneNumber uint			// Going down which number
	Repetition uint				// A round may have multiple repetitions of that scene
}

type Round struct {
	RoundIdx
	Stats RoundStats
	PreviousRound *Round
}

//TODO need total scene number TotalSceneCount

//TODO pick these up and tidy / fix.

func FirstRound() * Round {
	return &Round{RoundIdx{1, 1, 0}, NewRoundStats(), nil}
}

func (* Round) NextAct(currentRound * Round) * Round {
	//TODO need to refresh the stats.
	rIdx := RoundIdx{currentRound.RoundIdx.ActNumber + 1, 1, 0}
	return &Round{rIdx, currentRound.}
}

func (* Round) NextScene(currentRound * Round) * Round {
	now := time.Now()
	return &Round{currentRound.ActNumber, currentRound.SceneNumber+1, currentRound.Repetition,
		currentRound.TotalSceneNumber+1, currentRound.StartTime, time.Since(currentRound.StartTime),
		now, currentRound.Players, currentRound.CurrentLeader, currentRound.CurrentLoser,
		currentRound.CurrentMainLeader, currentRound.CurrentMainLoser, currentRound}
}

func (* Round) RepeatAct(currentRound * Round) * Round {
	now := time.Now()
	return &Round{currentRound.ActNumber, 1, currentRound.Repetition + 1,
		currentRound.TotalSceneNumber+1, currentRound.StartTime, time.Since(currentRound.StartTime),
		now, currentRound.Players, currentRound.CurrentLeader, currentRound.CurrentLoser,
		currentRound.CurrentMainLeader, currentRound.CurrentMainLoser, currentRound}
}



//TODO we need to make a round from next scene