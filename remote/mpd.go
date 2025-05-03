package remote

import (
	"log"
	"strings"
	"time"

	"github.com/fhs/gompd/v2/mpd"

	"git.sr.ht/~rockorager/vaxis/vxfw"
)

// struct to bundle access to the MPD daemon
type MpdRemote struct {
	lastQuery int
	chQuery   chan Query
	app       *vxfw.App
}

func (m *MpdRemote) PostQuery(tag string, constraints []Constraint) int {
	qid := m.lastQuery + 1
	m.lastQuery += 1

	m.chQuery <- Query{
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

		for {
			select {
			case q := <-m.chQuery:
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

			case <-time.After(2 * time.Second):
				songList, err := conn.PlaylistInfo(-1, -1)
				if err != nil {
					log.Printf("Failed: PlaylistInfo")
					log.Fatalf("MPD error: %v", err)
				}

				songs := make([]Attrs, 0, len(songList))
				for _, s := range songList {
					song := make(Attrs)
					for k, v := range(s) {
						song[k] = v
					}
					songs = append(songs, song)
				}
				m.app.PostEvent(songs)
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
