package lesson

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/alsenz/esl-games/pkg/account"
	"gorm.io/datatypes"
)

type RepeatCondition struct {
	Op  Op	`json:"op"`
	LHS string `json:"lhs"` //Templates to be resolved
	RHS string `json:"rhs"`
}

func (rc RepeatCondition) Eval(ctx *Context) (bool, error) {
	if lhs, err := ctx.Eval(rc.LHS); err == nil {
		if rhs, err := ctx.Eval(rc.RHS); err == nil {
			return rc.Op.Eval(lhs, rhs), nil
		} else {
			return false, err
		}
	} else {
		return false, err
	}
}

type RepeatUntilLogic struct {
	Fixed        uint16 `json:"fixed,omitempty"`              //65k repeats is enough :)
	CNF [][]RepeatCondition `json:"cnf,omitempty"` //Takes precedence if it exists. CNF.
}

func (rul *RepeatUntilLogic) Value() (driver.Value, error) {
	if raw, err := json.Marshal(rul); err != nil {
		return nil, err
	} else {
		return datatypes.JSON(raw).Value()
	}
}

func (rul *RepeatUntilLogic) Scan(src interface{}) error {
	jsn := &datatypes.JSON{}
	if err := jsn.Scan(src); err != nil {
		return err
	}
	return json.Unmarshal(*jsn, rul)
}

// Eval returns true if we should break out and continue (i.e. no longer repeat) - the repeat condition is satisfied
func (rl RepeatUntilLogic) Eval(ctx *Context) (bool, error) {
	if len(rl.CNF) == 0 {
		return ctx.Round.Repetition >= uint(rl.Fixed), nil
	}
	allCjsTrue := true
	for _, disjunct := range rl.CNF {
		anyDjTrue := false
		for _, literal := range disjunct {
			if eval, err := literal.Eval(ctx); err != nil {
				return false, err
			} else if eval {
				anyDjTrue = true
				break
			}
		}
		if !anyDjTrue {
			allCjsTrue = false
			break
		}
	}
	return allCjsTrue, nil
}

type Act struct {
	account.UserObject
	Scenes  Scenes `json:"scenes"`	//This should get gorm embedded as json
	RepeatUntil *RepeatUntilLogic `json:"repeatLogic,omitempty"` //This should get gorm embedded as json
}

func NewAct(owner account.User, group account.Group) *Act {
	act := &Act{}
	act.Owner = owner
	act.OwnerID = owner.ID
	act.Group = group
	act.GroupID = group.ID
	act.Permissions = *account.DefaultPermissions
	act.Scenes = make(Scenes, 0)
	act.RepeatUntil = nil
	return act
}

func (a *Act) CouldFinishAct(ctx *Context) (bool, error) {
	return a.RepeatUntil.Eval(ctx)
}