package lesson

type TeamLogic string
const (
	TeamLogicBySize		TeamLogic = "By Size"
	TeamLogicByNumberOf	TeamLogic = "By Number Of"
)

type TeamStructure struct {
	TeamLogic TeamLogic
	Count uint8
}
