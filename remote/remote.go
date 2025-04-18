package remote

type QueryType int

const (
	Songs QueryType = iota
	Playlist
)

type Constraint struct {
	Tag   string
	Op    string
	Query string
}

type Query struct {
	Query_id    int
	Type        QueryType
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
	PostQuery(tpe QueryType, tag string, constraints []Constraint) int
}
