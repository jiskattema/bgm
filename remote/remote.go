package remote

type Constraint struct {
	Tag   string
	Op    string
	Query string
}

type Query struct {
	Query_id    int
	Tag         string
	Constraints []Constraint
}

type Result struct {
	Result_id int
	Result    []string
}

type Remote interface {
	Dial() 
	HangUp()
	PostQuery(tag string, constraints []Constraint) int
}

var Tags = [6]string {
	"Artist", "Album", "Track", "Title", "Label", "Date",
}
