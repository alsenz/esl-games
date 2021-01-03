package lesson

import (
	"time"
)

type RoundStats struct {
	TotalSceneNumber uint		// sum of scenes in previous acts plus SceneNumber
	StartTime time.Time
	PlayTime time.Duration		// time since play began - time - startTime
	Time time.Time				// the actual time
}

func NewRoundStats() RoundStats {
	now := time.Now()
	return RoundStats{0, now, time.Since(now), now}
}

type GameInfo struct {
	Players []Player
	PreviousRound *Round
}

func NewGameInfo(players []Player) GameInfo {

}

type RoundIdx struct {
	ActNumber uint				// Number left-to-right - almost like an Act
	SceneNumber uint			// Going down which number
	Repetition uint				// A round may have multiple repetitions of that scene
}

type Round struct {
	RoundIdx
	Stats RoundStats
	Info GameInfo
	Scores GameScores
}

func FirstRound(players []Player) * Round {
	now := time.Now()
	return &Round{1, 1, 1, 0, now,
		time.Since(now), now, players, make(map[string]Player),
		make(map[string]Player), nil, nil, nil}
}

func (* Round) NextAct(currentRound * Round) * Round {
	now := time.Now()
	return &Round{currentRound.ActNumber+1, 1, 1, currentRound.TotalSceneNumber+1,
		currentRound.StartTime, time.Since(currentRound.StartTime), now, currentRound.Players,
		currentRound.CurrentLeader, currentRound.CurrentLoser, currentRound.CurrentMainLeader,
		currentRound.CurrentMainLoser, currentRound}
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