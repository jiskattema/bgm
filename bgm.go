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
	conn *mpd.Client
	lastQuery int
	chQuery chan mpd_query
	chResult chan mpd_result
}

func (m *MpdRemote) Dial() {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "192.168.1.110:6600")
	if err != nil {
		log.Fatalln(err)
	}
	m.conn = conn

	// setup up channels
	m.chQuery = make(chan mpd_query, 10)
	m.chResult = make(chan mpd_result, 10)

	// fire-off goroutine to do the querying
	go func() {
		for q := range(m.chQuery) {
			lines, err := mpd_remote.conn.List(q.query)
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
	m.conn.Close()
}

var mpd_remote MpdRemote

type Filter struct {
	Label string
	Count int
	Active bool
	current_query int
}

func NewFilter(label string) *Filter {
	return &Filter{
		Count: 0,
		Label: label,
		Active: false,
	}
}

// no-op for now
func (r *Filter) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	return nil, nil
}


func (r *Filter) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	count_text := strconv.Itoa(r.Count)
	spaces := int(ctx.Max.Width) - utf8.RuneCountInString(count_text) - utf8.RuneCountInString(r.Label)
	full_text := r.Label + strings.Repeat(" ", spaces) + strconv.Itoa(r.Count)

	chars := ctx.Characters(full_text)
	cells := make([]vaxis.Cell, 0, len(chars))

	style := InactiveFilter
	if r.Active {
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
		Widget: r,
		Cursor: nil,
		Buffer: cells,
		Children: []vxfw.SubSurface{},
	}, nil
}

type App struct {
	Filters [6]*Filter
	active bool
	active_filter int
}

func (a *App) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	switch ev := ev.(type) {
	case vaxis.Key:
		// Ctrl-C : quit
		if ev.Matches('c', vaxis.ModCtrl) {
			return vxfw.QuitCmd{}, nil
		}
		// Ctrl-u : halve a page up
		if ev.Matches('u', vaxis.ModCtrl) {
			if a.active_filter > 10 {
				a.active_filter -= 10
			} else {
				a.active_filter = 0
			}
		}
		// Ctrl-d : halve a page down
		if ev.Matches('d', vaxis.ModCtrl) {
			if a.active_filter < len(a.Filters) - 11 {
				a.active_filter += 10
			} else {
				a.active_filter = len(a.Filters) - 1
			}
		}
		// j : down
		if ev.Matches('j') {
			if a.active_filter < len(a.Filters) - 1 {
				a.active_filter +=1
			}
		}
		// G : go to bottom
		if ev.Matches('G') {
			a.active_filter = len(a.Filters) - 1
		}
		// k : up
		if ev.Matches('k') {
			if a.active_filter > 0 {
				a.active_filter -=1
			}
		}
		// g : go to top
		if ev.Matches('g') {
			a.active_filter = 0
		}
		// action on current row
		if ev.Matches(' ') {
			var qid = mpd_remote.lastQuery + 1
			mpd_remote.lastQuery += 1

			// fire-off a query
			mpd_remote.chQuery <- mpd_query{
				query_id: qid,
				query: a.Filters[a.active_filter].Label,
			}
			// save the query id to match against the result_id later
			a.Filters[a.active_filter].current_query = qid
		}
	}
	for pos, row := range(a.Filters) {
		row.Active = (pos == a.active_filter)
	}
	return nil, nil
}

func (a *App) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {

	// check if there are updates from the mpd_remote
	Poll:
	for {
		select {
		case r := <- mpd_remote.chResult :
			for _, row := range(a.Filters) {
				if row.current_query == r.result_id {
					row.Count = len(r.result)
				}
			}
		default:
			break Poll
		}
	}

	root := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, a)

	for pos, row := range(a.Filters) {
		surf, err := row.Draw(ctx)
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

	root := &App{
		Filters: [6]*Filter{
			NewFilter("Artist"),
			NewFilter("Album"),
			NewFilter("Track"),
			NewFilter("Title"),
			NewFilter("Label"),
			NewFilter("Date"),
		},
		active: true,
		active_filter: 2,
	}

	app.Run(root)
}
