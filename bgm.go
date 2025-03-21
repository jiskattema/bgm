package main

import (
	//"encoding/json"
	//"fmt"
	"log"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/textfield"
)

type Bgm struct {
	Filters []*Filter
	input   *textfield.TextField
	cursor  int
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
		if ev.Matches('j') {
			if b.cursor < len(b.Filters)-1 {
				b.cursor += 1
			}
		}
		// G : go to bottom
		if ev.Matches('G') {
			b.cursor = len(b.Filters) - 1
		}
		// k : up
		if ev.Matches('k') {
			if b.cursor > 0 {
				b.cursor -= 1
			}
		}
		// g : go to top
		if ev.Matches('g') {
			b.cursor = 0
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

			// fire-off a query
			mpd_remote.chQuery <- mpd_query{
				query_id: qid,
				tag:      filter.Label,
				query:    filter.Value,
			}
		}
	}
	for pos, filter := range b.Filters {
		filter.Active = (pos == b.cursor)
	}

	return vxfw.RedrawCmd{}, nil
}

func (b *Bgm) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	// check if there are updates from the mpd_remote
Poll:
	for {
		select {
		case r := <-mpd_remote.chResult:
			for _, filter := range b.Filters {
				if filter.current_query == r.result_id {
					filter.Count = len(r.result)
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

	// Note: the example creates a new context to set surface sizes
	s, err := b.input.Draw(ctx)
	if err != nil {
		return root, err
	}
	// root.AddChild(0, len(b.Filters), s)
	root.AddChild(0, int(ctx.Max.Height-1), s)

	return root, nil
}

var app *vxfw.App
var mpd_remote MpdRemote

func main() {
	mpd_remote.Dial()
	defer mpd_remote.HangUp()

	app, err := vxfw.NewApp(vaxis.Options{})
	if err != nil {
		log.Fatalf("Couldn't create a new app: %v", err)
	}

	root := &Bgm{
		cursor:  1,
		Filters: MpdFilters[:],
		input:   &textfield.TextField{},
	}

	root.Filters[root.cursor].Active = true

	app.Run(root)
}
