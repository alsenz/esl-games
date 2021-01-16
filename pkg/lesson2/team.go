package lesson2

type TeamLogic string
const (
	TeamLogicBySize		TeamLogic = "By Size"
	TeamLogicByNumberOf	TeamLogic = "By Number Of"
)

type TeamRules struct {
	TeamLogic TeamLogic
	Count uint8
}

type TeamName string