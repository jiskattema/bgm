package main

import (
	//"encoding/json"
	//"fmt"
	"log"

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

	// fire-off goroutine to do the querying
	go func() {
		// Connect to MPD server
		conn, err := mpd.Dial("tcp", "192.168.1.110:6600")
		if err != nil {
			log.Fatalln(err)
		}
		defer conn.Close()

		for q := range m.chQuery {
			lines, err := conn.List(q.query)
			if err != nil {
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
