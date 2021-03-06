package lesson2

type Round struct {
	Act uint
	Scene uint
	Rep uint
}

func (round *Round) LessThan(rhs Round) bool {
	return (round.Act < rhs.Act) ||
		(round.Act == rhs.Act && round.Scene < rhs.Scene) ||
		(round.Act == rhs.Act && round.Scene == rhs.Scene && round.Rep < rhs.Rep)
}