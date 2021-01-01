package account

import (
	"github.com/alsenz/esl-games/pkg/model"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type ACLFlag uint8

const (
	ACLNone ACLFlag = 0
	ACLRead ACLFlag = 1
	ACLReadWrite ACLFlag = 2
)

type ACL struct {
	Owner ACLFlag
	Group ACLFlag
	Others ACLFlag
}

func NewACL() *ACL {
	acl := ACL{ACLReadWrite, ACLNone, ACLNone}
	return &acl
}

type UserObject struct {
	model.Base
	OwnerID uuid.UUID
	Owner User //Gorm should use the field above to create an association by foreign key.
	GroupID uuid.UUID
	Group Group
	Permissions ACL		`gorm:"embedded;embeddedPrefix:acl_"`
}

func (userObject *UserObject) BeforeCreate(db *gorm.DB) error {
	userObject.ID = uuid.NewV4()
	userObject.Permissions = *NewACL()
	return nil
}