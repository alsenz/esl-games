package lesson

type Register struct {

}

type Lesson struct {
	Plan Plan
	Registration Register
	ResolvedQuestions map[RoundIdx]map[QuestionLink] Question	//Maps resolved questions asked at each round
	PlayerResponses map[RoundIdx]map[QuestionLink] PlayerResponses //TODO players need the responses on them too
	Context Context //Used for state keeping and evaluation of strings
}

func NewLesson(plan Plan, register Register) * Lesson {
	//TODO check Context has a New...
	//TODO learn how to initialise a blank map
	lesson := &Lesson{plan, register, MAKE MAP?}
	return lesson
}

func (lesson *Lesson) Run() error {
	//TODO the main event loop
	return nil
}

func (lesson *Lesson) BeforeLessonStart() error {

}

func (lesson *Lesson) StartRegistration() error {
	//TODO we need some registration logic
}

func (lesson *Lesson) AfterLessonStart() error {

}

func (lesson *Lesson) LessonStart() error {
	if err := lesson.BeforeLessonStart(); err != nil {
		return err
	}
	if err := lesson.StartRegistration(); err != nil {
		return err
	}
	if err := lesson.AfterLessonStart(); err != nil {
		return err
	}
}