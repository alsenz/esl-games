package account

import "github.com/alsenz/esl-games/pkg/model"

type User struct {
	model.Base
	Email string			`gorm:"unique;uniqueIndex"`
	Name string
	Groups []Group			`gorm:"many2many:user_groups;"`
	AdminOfGroups []Group	`gorm:"many2many:group_admins;"`
	Domain string //Likely just to be part of the email
	LastIdProvider string
}
