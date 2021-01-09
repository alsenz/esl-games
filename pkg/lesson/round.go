package lesson

import (
	"time"
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
	return RoundStats{rs.QuestionNumber, rs.StartTime, time.Now()}
}

func (rs *RoundStats) PlayTime() time.Duration {
	return time.Since(rs.StartTime)
}

type RoundIdx struct {
	ActNumber uint				// Number left-to-right - almost like an Act
	SceneNumber uint			// Going down which number
	Repetition uint				// A round may have multiple repetitions of that scene
}

type Round struct {
	RoundIdx
	Stats           RoundStats
	PreviousRound   *Round
	LinkedQuestions QuestionDraw
}

// Get the linked question if there is one, or a random linked question if there are many
func (round *Round) LinkedQuestion() (* Question) {
	for _, val := range round.LinkedQuestions {
		for _, val := range val {
			return &val
		}
		break
	}
	return nil
}

func FirstRound() * Round {
	return &Round{RoundIdx{1, 1, 0}, NewRoundStats(), nil,
		make(QuestionDraw)}
}

func (round* Round) NextAct(nextLinkedQuestions QuestionDraw) * Round {
	rIdx := RoundIdx{round.RoundIdx.ActNumber + 1, 1, 0}
	return &Round{rIdx, round.Stats.RefreshRoundStats(), round, nextLinkedQuestions}
}

func (round* Round) NextScene(nextLinkedQuestions QuestionDraw) * Round {
	rIdx := RoundIdx{round.RoundIdx.ActNumber, round.SceneNumber + 1, round.Repetition}
	return &Round{rIdx, round.Stats.RefreshRoundStats(), round, nextLinkedQuestions}
}

func (round* Round) RepeatAct(nextLinkedQuestions QuestionDraw) * Round {
	rIdx := RoundIdx{round.RoundIdx.ActNumber, 1, round.Repetition + 1}
	return &Round{rIdx, round.Stats.RefreshRoundStats(), round, nextLinkedQuestions}
}