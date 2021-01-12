package lesson

import (
	"time"
)

type Register struct {
	Timeout time.Duration
	Done bool
	RequireLogin bool
	OptDomain *string
	LessonCode string
}