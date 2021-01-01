package lesson

// A stripped-down version of user that can be baked into templates - not even an email address!
type Player struct {
	Name string
	Scores Scores
	Avatar []byte //This may or not contain an image... TODO TODO
}