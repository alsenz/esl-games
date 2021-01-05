package lesson

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/datatypes"
)

type TeamLogic string
const (
	TeamLogicBySize		TeamLogic = "By Size"
	TeamLogicByNumberOf	TeamLogic = "By Number Of"
)

type TeamStructure struct {
	TeamLogic TeamLogic		`json:"logic"`
	Count uint8				`json:"count"`
}

type TeamStructures map[string]TeamLogic

func (tsrcts *TeamStructures) Value() (driver.Value, error) {
	if raw, err := json.Marshal(tsrcts); err != nil {
		return nil, err
	} else {
		return datatypes.JSON(raw).Value()
	}
}

func (tsrcts *TeamStructures) Scan(src interface{}) error {
	jsn := &datatypes.JSON{}
	if err := jsn.Scan(src); err != nil {
		return err
	}
	return json.Unmarshal(*jsn, tsrcts)
}
