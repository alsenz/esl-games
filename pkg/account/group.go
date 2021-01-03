package account

import "github.com/alsenz/esl-games/pkg/model"

const (
	CreatorsGroupName string = "creators"
	AdminsGroupName string = "admin"
	SuperAdminsGroupName string = "super-admin"
	PlayersGroupName string = "players"
)

type Group struct {
	model.Base
	Name string		`gorm:"unique,uniqueIndex" json:"name"`
	Admins []User	`gorm:"many2many:group_admins," json:"-"`
}
