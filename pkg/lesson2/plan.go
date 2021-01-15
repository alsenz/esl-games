package lesson2

import "github.com/alsenz/esl-games/pkg/account"

type Plan struct {
	account.UserObject
	Name        string
	Description string
	TeamRules   TeamRules //TODO embed
	Acts        Acts      //TODO jsonify
}

