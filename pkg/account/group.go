package account

import "github.com/alsenz/esl-games/pkg/model"

type Group struct {
	model.Base
	Name string		`gorm:"unique,uniqueIndex"`
	Admins []User	`gorm:"many2many:group_admins,"`
}
