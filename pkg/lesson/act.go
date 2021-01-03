package lesson

import (
	"github.com/alsenz/esl-games/pkg/account"
	"gorm.io/datatypes"
)

type RepeatCondition struct {
	Op Op
	LHS string //Templates to be resolved
	RHS string
}

type RepeatLogic struct {
	Fixed uint16 //65k repeats is enough :)
	ConditionCNF [][]RepeatCondition //Takes precedence if it exists. CNF.
}

//TODO RepeatLogic needs an Eval() based on Round plus info
//TODO TODO

type Act struct {
	account.UserObject
	ScenesJSON datatypes.JSON	//TODO at these stage these need to be JSON with a method to get them.
	RepeatLogic RepeatLogic //TODO gorm embed this thing
	//TODO we need a repeat condition.
}

func (a *Act) Scenes() []Scene {
	//TODO we need to convert this old thing
}
