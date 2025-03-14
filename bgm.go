package main

import (
	"log"
	"strconv"
	"strings"
	"unicode/utf8"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"

	"github.com/fhs/gompd/v2/mpd"
)

var InactiveRow = vaxis.Style{
	Foreground: vaxis.HexColor(0x00ffff),
	Background: vaxis.HexColor(0xff0fff),
}

var ActiveRow = vaxis.Style{
	Foreground: vaxis.HexColor(0x000000),
	Background: vaxis.HexColor(0xffffff),
}

type Row struct {
	Label string
	Count int
	Active bool
}

func NewRow(label string) *Row {
	return &Row{
		Label: label,
		Active: false,
	}
}

// no-op for now
func (r *Row) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	return nil, nil
}


func (r *Row) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	count_text := strconv.Itoa(r.Count)
	spaces := int(ctx.Max.Width) - utf8.RuneCountInString(count_text) - utf8.RuneCountInString(r.Label)
	full_text := r.Label + strings.Repeat(" ", spaces) + strconv.Itoa(r.Count)

	chars := ctx.Characters(full_text)
	cells := make([]vaxis.Cell, 0, len(chars))

	style := InactiveRow
	if r.Active {
	  style = ActiveRow
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
	Rows [5]*Row
	active bool
	active_row int
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
			if a.active_row > 10 {
				a.active_row -= 10
			} else {
				a.active_row = 0
			}
		}
		// Ctrl-d : halve a page down
		if ev.Matches('d', vaxis.ModCtrl) {
			if a.active_row < len(a.Rows) - 11 {
				a.active_row += 10
			} else {
				a.active_row = len(a.Rows) - 1
			}
		}
		// j : down
		if ev.Matches('j') {
			if a.active_row < len(a.Rows) - 1 {
				a.active_row +=1
			}
		}
		// G : go to bottom
		if ev.Matches('G') {
			a.active_row = len(a.Rows) - 1
		}
		// k : up
		if ev.Matches('k') {
			if a.active_row > 0 {
				a.active_row -=1
			}
		}
		// g : go to top
		if ev.Matches('g') {
			a.active_row = 0
		}
		// action on current row
		if ev.Matches(' ') {
			a.Rows[a.active_row].Count += 1
		}
	}
	for pos, row := range(a.Rows) {
		row.Active = (pos == a.active_row)
	}
	return nil, nil
}

func (a *App) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	root := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, a)

	for pos, row := range(a.Rows) {
		surf, err := row.Draw(ctx)
		if err != nil {
			return vxfw.Surface{}, err
		}
		root.AddChild(0, pos, surf)
	}

	return root, nil
}

type MpdRemote struct {
	conn *mpd.Client
}

func (m *MpdRemote) Dial() {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "192.168.1.110:6600")
	if err != nil {
		log.Fatalln(err)
	}
	m.conn = conn
}

func (m *MpdRemote) HangUp() {
	m.conn.Close()
}

func main() {
	var mpd_remote MpdRemote
	
	mpd_remote.Dial()
	defer mpd_remote.HangUp()

	app, err := vxfw.NewApp()
	if err != nil {
		log.Fatalf("Couldn't create a new app: %v", err)
	}

	root := &App{
		Rows: [5]*Row{
			NewRow("Artist"),
			NewRow("Album"),
			NewRow("Track"),
			NewRow("Length"),
			NewRow("Year"),
		},
		active: true,
		active_row: 2,
	}

	app.Run(root)
}
