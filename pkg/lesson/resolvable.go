package lesson

import "gorm.io/gorm"

type Resolvable interface {
	IsTemplated() bool
	Resolve(*Round, *gorm.DB) *Question
}
