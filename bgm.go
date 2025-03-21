package main

import (
	//"encoding/json"
	//"fmt"
	"log"
	"reflect"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/list"
	"git.sr.ht/~rockorager/vaxis/vxfw/text"
	"git.sr.ht/~rockorager/vaxis/vxfw/textfield"
)

type Bgm struct {
	Filters []*Filter
	input   *textfield.TextField
	cursor  int
	list    list.Dynamic
}

func (b *Bgm) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	switch ev := ev.(type) {
	case vaxis.Key:
		// Ctrl-C : quit
		if ev.Matches('c', vaxis.ModCtrl) {
			return vxfw.QuitCmd{}, nil
		}
		// Ctrl-u : half a page up
		if ev.Matches('u', vaxis.ModCtrl) {
			if b.cursor > 10 {
				b.cursor -= 10
			} else {
				b.cursor = 0
			}
		}
		// Ctrl-d : half a page down
		if ev.Matches('d', vaxis.ModCtrl) {
			if b.cursor < len(b.Filters)-11 {
				b.cursor += 10
			} else {
				b.cursor = len(b.Filters) - 1
			}
		}
		// j : down
		if ev.Matches('j') || ev.Matches(vaxis.KeyDown) {
			if b.cursor < len(b.Filters)-1 {
				b.cursor += 1
			}
		}
		// G : go to bottom
		if ev.Matches('G') || ev.Matches(vaxis.KeyEnd) {
			b.cursor = len(b.Filters) - 1
		}
		// k : up
		if ev.Matches('k') || ev.Matches(vaxis.KeyUp) {
			if b.cursor > 0 {
				b.cursor -= 1
			}
		}
		// g : go to top
		if ev.Matches('g') || ev.Matches(vaxis.KeyHome) {
			b.cursor = 0
		}
		// Ctrl-p : move filter up
		if ev.Matches('p', vaxis.ModCtrl) {
			if b.cursor > 0 {
				s := reflect.Swapper(b.Filters)
				s(b.cursor-1, b.cursor)
				b.cursor -= 1
			}
		}
		// Ctrl-n : move filter down
		if ev.Matches('n', vaxis.ModCtrl) {
			if b.cursor < len(b.Filters)-1 {
				s := reflect.Swapper(b.Filters)
				s(b.cursor, b.cursor+1)
				b.cursor += 1
			}
		}
		// Ctrl-t : move filter to top
		if ev.Matches('t', vaxis.ModCtrl) {
			s := reflect.Swapper(b.Filters)
			for p := b.cursor; p > 0; p -= 1 {
				s(p, p-1)
			}
			b.cursor = 0
		}
		// Ctrl-b : move filter to bottom
		if ev.Matches('b', vaxis.ModCtrl) {
			s := reflect.Swapper(b.Filters)
			for p := b.cursor; p < len(b.Filters)-1; p += 1 {
				s(p, p+1)
			}
		}
		// Enter : show matches for current filter in bottom pane
		if ev.Matches(vaxis.KeyEnter) {
			items = b.Filters[b.cursor].matches[:]
			b.list.SetCursor(0)
		}
		// ] : select next match
		if ev.Matches(']') {
			cursor := b.Filters[b.cursor].cursor + 1
			if cursor >= len(b.Filters[b.cursor].matches) {
				cursor = -1
			}
			b.Filters[b.cursor].cursor = cursor
		}
		// [ : select previous match
		if ev.Matches('[') {
			cursor := b.Filters[b.cursor].cursor - 1
			if cursor < -1 {
				cursor = len(b.Filters[b.cursor].matches) - 1
			}
			b.Filters[b.cursor].cursor = cursor
		}
		// / : search
		if ev.Matches('/') {
			// Set callback
			b.input.OnSubmit = func(line string) (vxfw.Command, error) {
				b.Filters[b.cursor].Value = line
				b.input.OnSubmit = nil
				return vxfw.FocusWidgetCmd(b), nil
			}
			// Focus the input widget
			b.input.Reset()
			b.input.InsertStringAtCursor(b.Filters[b.cursor].Value)
			return vxfw.FocusWidgetCmd(b.input), nil
		}
		// action on current filter
		if ev.Matches(' ') {
			var qid = mpd_remote.newQueryId()

			// save the query id to match against the result_id later
			filter := b.Filters[b.cursor]
			filter.current_query = qid

			// create tag + query pairs for each filter up-to and including the cursor
			tq := make([]tag_query, b.cursor+1)
			for i := 0; i <= b.cursor; i++ {
				f := b.Filters[i]
				if f.cursor >= 0 && len(f.matches) > f.cursor {
					// exact match of value at cursor
					tq = append(tq, tag_query{
						tag:   f.Label,
						op:    "==",
						query: f.matches[f.cursor],
					})
				} else {
					// search for value using 'contains'
					tq = append(tq, tag_query{
						tag:   f.Label,
						op:    "contains",
						query: f.Value,
					})
				}
			}

			// fire-off a query
			mpd_remote.chQuery <- mpd_query{
				query_id:    qid,
				tag:         filter.Label,
				constraints: tq,
			}
		}
		// Tab : focus on bottom panel
		if ev.Matches(vaxis.KeyTab) {
			return vxfw.FocusWidgetCmd(&b.list), nil
		}
		// Esc : focus on top panel
		if ev.Matches(vaxis.KeyEsc) {
			return vxfw.FocusWidgetCmd(b), nil
		}
	}
	for pos, filter := range b.Filters {
		filter.Active = (pos == b.cursor)
	}

	return vxfw.RedrawCmd{}, nil
}

func getWidget(i uint, cursor uint) vxfw.Widget {
	if i >= uint(len(items)) {
		return nil
	}
	var style vaxis.Style
	if i == cursor {
		style.Attribute = vaxis.AttrReverse
	}
	var display_text string
	if items[i] == "" {
		display_text = "[Unknown]"
	} else {
		display_text = items[i]
	}
	return &text.Text{
		Content: display_text,
		Style:   style,
	}
}

func (b *Bgm) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	// check if there are updates from the mpd_remote
Poll:
	for {
		select {
		case r := <-mpd_remote.chResult:
			for _, filter := range b.Filters {
				if filter.current_query == r.result_id {
					filter.matches = r.result
					filter.cursor = -1
				}
			}
		default:
			break Poll
		}
	}

	root := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, b)

	for pos, filter := range b.Filters {
		surf, err := filter.Draw(ctx)
		if err != nil {
			return root, err
		}
		root.AddChild(0, pos, surf)
	}

	panel_size := vxfw.Size{Width: ctx.Max.Width, Height: 1}

	// The commandline
	panel_size.Height = 1
	s, err := b.input.Draw(vxfw.DrawContext{Min: panel_size, Max: panel_size, Characters: ctx.Characters})
	if err != nil {
		return root, err
	}
	root.AddChild(0, int(ctx.Max.Height-1), s)

	// full item list
	panel_size.Height = ctx.Max.Height - uint16(len(b.Filters)) - 1
	s, err = b.list.Draw(vxfw.DrawContext{Min: panel_size, Max: panel_size, Characters: ctx.Characters})
	if err != nil {
		return root, err
	}
	root.AddChild(0, len(b.Filters), s)

	return root, nil
}

var app *vxfw.App
var mpd_remote MpdRemote
var items []string

func main() {
	mpd_remote.Dial()
	defer mpd_remote.HangUp()

	app, err := vxfw.NewApp(vaxis.Options{})
	if err != nil {
		log.Fatalf("Couldn't create a new app: %v", err)
	}

	root := &Bgm{
		cursor:  0,
		Filters: MpdFilters[:],
		input:   &textfield.TextField{},
		list: list.Dynamic{
			Builder:              getWidget,
			DrawCursor:           false,
			Gap:                  0,
			DisableEventHandlers: true,
		},
	}

	app.Run(root)
}
