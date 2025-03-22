package main

import (
	//"encoding/json"
	"log"
	"strings"

	"github.com/fhs/gompd/v2/mpd"
)

var MpdFilters = [6]*Filter{
	{Label: "Artist"},
	{Label: "Album"},
	{Label: "Track"},
	{Label: "Title"},
	{Label: "Label"},
	{Label: "Date"},
}

type tag_query struct {
	tag   string
	op    string
	query string
}

type mpd_query struct {
	query_id    int
	tag         string
	constraints []tag_query
}

type mpd_result struct {
	result_id int
	result    []string
}

// struct to bundle access to the MPD daemon
type MpdRemote struct {
	lastQuery int
	chQuery   chan mpd_query
	chResult  chan mpd_result
}

func (m *MpdRemote) newQueryId() int {
	qid := m.lastQuery + 1
	m.lastQuery += 1
	return qid
}

func (m *MpdRemote) Dial() {
	// setup up channels
	m.chQuery = make(chan mpd_query, 10)
	m.chResult = make(chan mpd_result, 10)

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
			var sb strings.Builder
			ccount := 0

			sb.WriteString("(")
			for _, tq := range q.constraints {
				if tq.query != "" {
					if ccount > 0 {
						sb.WriteString(" AND ")
					}
					sb.WriteString("(" + tq.tag + " " + tq.op + " '" + illegalChars.Replace(tq.query) + "')")
					ccount += 1
				}
			}
			sb.WriteString(")")

			var lines []string

			if ccount > 0 {
				lines, err = conn.List(q.tag, sb.String())
			} else {
				lines, err = conn.List(q.tag)
			}

			if err != nil {
				log.Printf("Failed: %s", sb.String())
				log.Fatalf("MPD error: %v", err)
			}
			m.chResult <- mpd_result{
				result_id: q.query_id,
				result:    lines,
			}
		}
		close(m.chResult)
	}()
}

func (m *MpdRemote) HangUp() {
	close(m.chQuery)
}
