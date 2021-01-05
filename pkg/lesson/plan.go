package lesson

import "github.com/alsenz/esl-games/pkg/account"

type Plan struct {
	account.UserObject
	Name string
	Description string
	TeamStructures map[string]TeamLogic //TODO embed this (check for some better gorm logic here) TODO TODO
	Acts []Act //TODO one-many mapping to acts
}

//TODO any utilities here that would help
