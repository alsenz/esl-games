package lesson

type ScoreName string

const DefaultScoreName ScoreName = "default"

type ScoreCard map[ScoreName]float64

type PlayerRanking struct {
	ScoreName ScoreName
	Player Player
	Score float64
}
type PlayerRankings []PlayerRanking



type TeamRanking struct {
	ScoreName
	Team Team
	Score float64
}
type TeamRankings []TeamRanking

type PlayerScores map[ScoreName]PlayerRankings
type TeamScores map[ScoreName]TeamRankings

type Scores struct {
	PlayerScores PlayerScores
	TeamScores TeamScores
}

type ScoringLogic string
const (
	Simple 			 ScoringLogic = "simple"
	AddScores		 ScoringLogic = "add-scores" //Adds two scores together
	InverseFrequency ScoringLogic = "inverse-frequency"
	SimpleByTeam	 ScoringLogic = "simple-by-team"
	AddScoresByTeam	 ScoringLogic = "add-by-team"
)

type ScoringRules map[ScoreName]ScoringLogic

//TODO we need some nice scoring routines here once we know more about what these look like.
