package remote

import (
	"fmt"
	"log"
	"strings"

	"github.com/fhs/gompd/v2/mpd"

	"git.sr.ht/~rockorager/vaxis/vxfw"
)

// struct to bundle access to the MPD daemon
type MpdRemote struct {
	lastQuery int
	chQuery   chan Query
	app       *vxfw.App
}

func (m *MpdRemote) PostQuery(tpe QueryType, tag string, constraints []Constraint) int {
	qid := m.lastQuery + 1
	m.lastQuery += 1

	m.chQuery <- Query{
		Type:        tpe,
		Query_id:    qid,
		Tag:         tag,
		Constraints: constraints,
	}
	return qid
}

func (m *MpdRemote) Dial() {
	// setup up channels
	m.chQuery = make(chan Query, 10)

	illegalChars := strings.NewReplacer("'", `\'`, `"`, `\"`)

	// fire-off goroutine to do the querying
	go func() {
		// Connect to MPD server
		conn, err := mpd.Dial("tcp", "localhost:6600")
		if err != nil {
			log.Fatalln(err)
		}
		defer conn.Close()

		for q := range m.chQuery {
			switch q.Type {
			case Songs:
				var sb strings.Builder
				ccount := 0

				sb.WriteString("(")
				for _, constraint := range q.Constraints {
					if constraint.Query != "" {
						if ccount > 0 {
							sb.WriteString(" AND ")
						}
						sb.WriteString("(" + constraint.Tag + " " + constraint.Op + " '" + illegalChars.Replace(constraint.Query) + "')")
						ccount += 1
					}
				}
				sb.WriteString(")")

				var lines []string

				if ccount > 0 {
					lines, err = conn.List(q.Tag, sb.String())
				} else {
					lines, err = conn.List(q.Tag)
				}

				if err != nil {
					log.Printf("Failed: %s", sb.String())
					log.Fatalf("MPD error: %v", err)
				}
				// TODO: use vx.PostCommand
				m.app.PostEvent(Result{
					Result_id: q.Query_id,
					Result:    lines,
				})
			case Playlist:
				attrs, err := conn.PlaylistInfo(-1, -1)
				if err != nil {
					log.Printf("Failed: PlaylistInfo")
					log.Fatalf("MPD error: %v", err)
				}

				var lines []string
				for att := range attrs {
					lines = append(lines, fmt.Sprintf("%v", att))
				}

				m.app.PostEvent(Result{
					Result_id: q.Query_id,
					Result:    lines,
				})
			default:
				panic("Unknown Query")
			}
		}
	}()
}

func (m *MpdRemote) HangUp() {
	close(m.chQuery)
}

func New(app *vxfw.App) *MpdRemote {
	return &MpdRemote{
		app: app,
	}
}
