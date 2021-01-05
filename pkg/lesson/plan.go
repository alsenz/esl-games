package lesson

import "github.com/alsenz/esl-games/pkg/account"

type Plan struct {
	account.UserObject
	Name string
	Description string
	TeamStructures TeamStructures //This is embedded
	Acts []Act // This is also embedded.
}

//TODO any utilities here that would help
