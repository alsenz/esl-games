package lesson

type ScoreName string

const MainScoreName ScoreName = "main"

type ScoreCard map[ScoreName]float64

type PlayerScoreSheet struct {
	Player Player
	AllRounds map[RoundIdx]ScoreCard
	Rounds map[RoundIdx]float64
	Cumulative ScoreCard
}

type PlayerRanking struct {
	ScoreName ScoreName
	Player Player
	Score float64
}

type PlayerRankings []PlayerRanking

type GameScores struct {
	ScoreSheets []PlayerScoreSheet
	AllRankings map[ScoreName]PlayerRankings
	Rankings PlayerRankings //For the default "main" score name
	AllLeaders map[ScoreName]Player
	Lead Player //For the default "main" score name
	AllLosers map[ScoreName]Player
	Loser Player //For the default "main" score name
}

type ScoringLogic string
const (
	Simple 			 ScoringLogic = "simple"
	AddScores		 ScoringLogic = "add-scores" //Adds two scores together
	InverseFrequency ScoringLogic = "inverse-frequency"
)

type ScoringRules map[ScoreName]ScoringLogic

//TODO we need some nice scoring routines here once we know more about what these look like.
