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
	RegistrationSync chan bool //This is to listen for the "early ready" event of someone ending the registration period mnaually
}