package lesson2

type ScoreName string

type RoundScore struct {
	ScoreName ScoreName
	Score float64
	Ranking uint64
}

type ScoreCard struct {
	RoundScore
	AccumulatedScore float64
	AccumulatedRanking uint64
}

type ScoreCards map[ScoreName]ScoreCard

type ScoreAccumulatorExpr string
var PresetScoreAccumulatorExprs = map[string]ScoreAccumulatorExpr{
	//TODO let's load sprig funcs
	"sum": "addf .Accumulated.Score .Next.Score",
	"max": "maxf .Accumulated.Score .Next.Score",
}

type ScoreDefinitions map[ScoreName]ScoreAccumulatorExpr

func (expr ScoreAccumulatorExpr) Eval(score RoundScore, card ScoreCard) float64 {
	//TODO need to load spring funcs
	//TODO let's revisit evalling a go template... remembering extensions here
	//TODO TODO
}

//TODO for the Act repeater logic, we need some kind of accumulated scoring expression here
//TODO TODO
