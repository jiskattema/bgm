package base

// To ask the user to enter a string; the callback is called with both the question and the users' answer.
type Question struct {
	Question string
	Callback func(question, anwser string)
}

// A search criteria: TAG OP VALUE , 'Title' '=' 'One'
type Constraint struct {
	Tag   string
	Op    string
	Value string
}

// For querying a remote (ie. MPD) for all Tags of songs matching the criteria.
type Query struct {
	Query_id    int
	Tag         string
	Constraints []Constraint
}

// The result of a query
type Result struct {
	Query_id int
	Result   []string
}

// the attributes of a song
type Attrs map[string]string

// a slice containing songs
type Playlist []Attrs
