package lesson2

import "github.com/alsenz/esl-games/pkg/account"

type Plan struct {
	account.UserObject
	Name        string
	Description string
	Tags Tags //TODO jsonify
	TeamRules   TeamRules //TODO embed
	ScoreDefinitions ScoreDefinitions //TODO jsonify
	Acts        Acts      //TODO jsonify
}

