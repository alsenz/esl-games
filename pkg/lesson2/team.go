package lesson2

type TeamLogic string
const (
	TeamLogicNoTeams	TeamLogic = "No Teams"
	TeamLogicBySize		TeamLogic = "By Size"
	TeamLogicByNumberOf	TeamLogic = "By Number Of"
)

type TeamRules struct {
	TeamLogic TeamLogic
	Count int
}

type TeamName string