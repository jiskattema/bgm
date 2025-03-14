package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"

	"github.com/fhs/gompd/v2/mpd"
)

func ExampleDial() {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	line := ""
	line1 := ""
	// Loop printing the current status of MPD.
	for {
		status, err := conn.Status()
		if err != nil {
			log.Fatalln(err)
		}
		song, err := conn.CurrentSong()
		if err != nil {
			log.Fatalln(err)
		}
		if status["state"] == "play" {
			line1 = fmt.Sprintf("%s - %s", song["Artist"], song["Title"])
		} else {
			line1 = fmt.Sprintf("State: %s", status["state"])
		}
		if line != line1 {
			line = line1
			fmt.Println(line)
		}
		time.Sleep(1e9)
	}
}

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
	full_text := r.Label + " " + strconv.Itoa(r.Count)

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
	r1 *Row
	r2 *Row
	r3 *Row
	active bool
}

func (a *App) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	switch ev := ev.(type) {
	case vaxis.Key:
		if ev.Matches('c', vaxis.ModCtrl) {
			return vxfw.QuitCmd{}, nil
		}
		if ev.Matches('a') {
			a.r1.Active = true
			a.r2.Active = false
			a.r3.Active = false
			a.r1.Count += 1
		}
		if ev.Matches('b') {
			a.r1.Active = false
			a.r2.Active = true
			a.r3.Active = false
			a.r2.Count += 1
		}
		if ev.Matches('d') {
			a.r1.Active = false
			a.r2.Active = false
			a.r3.Active = true
			a.r3.Count += 1
		}
	}
	return nil, nil
}

func (a *App) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	root := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, a)

	s1, err := a.r1.Draw(ctx)
	s2, err := a.r2.Draw(ctx)
	s3, err := a.r3.Draw(ctx)
	if err != nil {
		return vxfw.Surface{}, err
	}
	root.AddChild(0, 0, s1)
	root.AddChild(0, 1, s2)
	root.AddChild(0, 2, s3)

	return root, nil
}

func main() {
	ExampleDial()
	app, err := vxfw.NewApp()
	if err != nil {
		log.Fatalf("Couldn't create a new app: %v", err)
	}

	root := &App{
	  r1: NewRow("Artist"),
	  r2: NewRow("Album"),
	  r3: NewRow("Track"),
	}

	app.Run(root)
}

