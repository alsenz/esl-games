package account

import (
	"github.com/alsenz/esl-games/pkg/model"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type User struct {
	model.Base
	Email string			`gorm:"unique;uniqueIndex"`
	Name string
	Groups []Group			`gorm:"many2many:user_groups;"`
	AdminOfGroups []Group	`gorm:"many2many:group_admins;"`
	Domain string //Likely just to be part of the email
	LastIdProvider string
}

func CheckAuth(r *http.Request) (* User, error) {
	return nil, nil
}

func (user *User) GroupIDs() []uuid.UUID {
	result := make([]uuid.UUID, len(user.Groups))
	for _, grp := range user.Groups {
		result = append(result, grp.ID)
	}
	return result
}