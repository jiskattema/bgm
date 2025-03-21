package main

import (
	//"encoding/json"
	"fmt"
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

type mpd_query struct {
	query_id int
	tag      string
	query    string
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
			quoted := illegalChars.Replace(q.query)
			formatted := fmt.Sprintf("(%s contains '%s')", q.tag, quoted)
			lines, err := conn.List(q.tag, formatted)
			if err != nil {
				log.Printf("Failed: %s", formatted)
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
