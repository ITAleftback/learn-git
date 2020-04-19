package str

type Lesson struct {
	Stunum   string
	Lesson  string
	Class   string
	Typee   string
	Teacher string
	Time    string
	Place   string
}

type Student struct {
	Name    string
	Stunum   string
	Lesson []Lesson
}
type Selectlesson struct {
	Lesson   string
	Typee   string
	Teacher string
}