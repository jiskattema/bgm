package main

import (
	//"encoding/json"
	//"fmt"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"

	"github.com/fhs/gompd/v2/mpd"
)

var InactiveFilter = vaxis.Style{
	Foreground: vaxis.HexColor(0x00ffff),
	Background: vaxis.HexColor(0xff0fff),
}

var ActiveFilter = vaxis.Style{
	Foreground: vaxis.HexColor(0x000000),
	Background: vaxis.HexColor(0xffffff),
}

type mpd_query struct {
	query_id int
	query string
}

type mpd_result struct {
	result_id int
	result []string
}

// struct to bundle access to the MPD daemon
type MpdRemote struct {
	lastQuery int
	chQuery chan mpd_query
	chResult chan mpd_result
}

func (m *MpdRemote) newQueryId() (int) {
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

		for q := range(m.chQuery) {
			lines, err := conn.List(q.query)
			if err != nil {
				log.Fatalf("MPD error: %v", err)
			}
			m.chResult <- mpd_result{
				result_id: q.query_id,
				result: lines,
			}
		}
		close(m.chResult)
	}()
}

func (m *MpdRemote) HangUp() {
	close(m.chQuery)
}

var mpd_remote MpdRemote
var app *vxfw.App

type Filter struct {
	Label string
	Count int
	Active bool
	current_query int
}

// no-op for now
func (r *Filter) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	return nil, nil
}


func (f *Filter) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	count_text := strconv.Itoa(f.Count)
	spaces := int(ctx.Max.Width) - utf8.RuneCountInString(count_text) - utf8.RuneCountInString(f.Label)
	full_text := f.Label + strings.Repeat(" ", spaces) + strconv.Itoa(f.Count)

	chars := ctx.Characters(full_text)
	cells := make([]vaxis.Cell, 0, len(chars))

	style := InactiveFilter
	if f.Active {
	  style = ActiveFilter
	}

	var w int
	for _, char := range chars {
		cell := vaxis.Cell{
			Character: char,
			Style: style, 
		}
		cells = append(cells, cell)
		w += char.Width
	}

	return vxfw.Surface{
		Size: vxfw.Size{Width: uint16(w), Height: 1}, 
		Widget: f,
		Cursor: nil,
		Buffer: cells,
		Children: []vxfw.SubSurface{},
	}, nil
}

type Bgm struct {
	Filters [6]Filter
	active bool
	cursor int
}

func (b *Bgm) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	switch ev := ev.(type) {
	case vaxis.Key:
		// Ctrl-C : quit
		if ev.Matches('c', vaxis.ModCtrl) {
			return vxfw.QuitCmd{}, nil
		}
		// Ctrl-u : halve a page up
		if ev.Matches('u', vaxis.ModCtrl) {
			if b.cursor > 10 {
				b.cursor -= 10
			} else {
				b.cursor = 0
			}
		}
		// Ctrl-d : halve a page down
		if ev.Matches('d', vaxis.ModCtrl) {
			if b.cursor < len(b.Filters) - 11 {
				b.cursor += 10
			} else {
				b.cursor = len(b.Filters) - 1
			}
		}
		// j : down
		if ev.Matches('j') {
			if b.cursor < len(b.Filters) - 1 {
				b.cursor +=1
			}
		}
		// G : go to bottom
		if ev.Matches('G') {
			b.cursor = len(b.Filters) - 1
		}
		// k : up
		if ev.Matches('k') {
			if b.cursor > 0 {
				b.cursor -=1
			}
		}
		// g : go to top
		if ev.Matches('g') {
			b.cursor = 0
		}
		// action on current filter
		if ev.Matches(' ') {
			var qid = mpd_remote.newQueryId()

			// save the query id to match against the result_id later
			filter := b.Filters[b.cursor]
			filter.current_query = qid

			// fire-off a query
			mpd_remote.chQuery <- mpd_query{
			query_id: qid,
				query: filter.Label,
			}
		}
	}
	for pos, filter := range(b.Filters) {
		filter.Active = (pos == b.cursor)
	}
	return vxfw.RedrawCmd{}, nil
}

func (b *Bgm) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	// check if there are updates from the mpd_remote
	Poll:
	for {
		select {
		case r := <- mpd_remote.chResult :
			for _, filter := range(b.Filters) {
				if filter.current_query == r.result_id {
					filter.Count = len(r.result)
				}
			}
		default:
			break Poll
		}
	}

	root := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, b)

	for pos, filter := range(b.Filters) {
		surf, err := filter.Draw(ctx)
		if err != nil {
			return vxfw.Surface{}, err
		}
		root.AddChild(0, pos, surf)
	}

	return root, nil
}

func main() {
	mpd_remote.Dial()
	defer mpd_remote.HangUp()

	app, err := vxfw.NewApp()
	if err != nil {
		log.Fatalf("Couldn't create a new app: %v", err)
	}

	root := &Bgm{
		Filters: [6]Filter{
			{ Label: "Artist", },
			{ Label: "Album", },
			{ Label: "Track", },
			{ Label: "Title", },
			{ Label: "Label", },
			{ Label: "Date", },
		},
		active: true,
		cursor: 1,
	}

	root.Filters[root.cursor].Active = true

	app.Run(root)
}
