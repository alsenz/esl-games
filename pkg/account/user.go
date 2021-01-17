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

func (user *User) CanReadObject(userObject *UserObject) bool {
	if userObject.Permissions.Others == ACLRead || userObject.Permissions.Others == ACLReadWrite {
		// If other permissions are read, continue regardless
		return true
	}
	if userObject.OwnerID == user.ID {
		// We are the owner - continue append if we owner permissions
		if userObject.Permissions.Owner == ACLRead || userObject.Permissions.Group == ACLReadWrite {
			return true
		}
	}
	// Final case - are we in the group for this?
	for _, tGrpID := range user.GroupIDs() {
		if tGrpID == userObject.GroupID {
			//check the group case
			if userObject.Permissions.Group == ACLRead || userObject.Permissions.Group == ACLReadWrite {
				return true
			}
			break
		}
	}
	return false
}

func (user *User) CanWriteObject(userObject *UserObject) bool {
	if userObject.Permissions.Others == ACLReadWrite {
		// If other permissions are read, continue regardless
		return true
	}
	if userObject.OwnerID == user.ID {
		// We are the owner - continue append if we owner permissions
		if userObject.Permissions.Group == ACLReadWrite {
			return true
		}
	}
	// Final case - are we in the group for this?
	for _, tGrpID := range user.GroupIDs() {
		if tGrpID == userObject.GroupID {
			//check the group case
			if userObject.Permissions.Group == ACLReadWrite {
				return true
			}
			break
		}
	}
	return false
}