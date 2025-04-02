package search

import (
	"reflect"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/list"
	"git.sr.ht/~rockorager/vaxis/vxfw/text"

	"ash/bgm/remote"
)

var items []string

type Search struct {
	// The content of the widget
	Filters []*Filter
	list    list.Dynamic
	cursor  int
	remote  remote.Remote
}

func New() *Search{
	var filters []*Filter
	for _, name:= range(remote.Tags) {
		filters = append(
			filters,
			&Filter{Label: name},
		)
	}

	s := &Search{
		list: list.Dynamic{
			Builder:              getWidget,
			DrawCursor:           false,
			Gap:                  0,
			DisableEventHandlers: true,
		},
		Filters: filters,
	}

	return s
}

// Noop for text
func (r *Search) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	switch ev := ev.(type) {
	case vaxis.Key:
		// Ctrl-u : half a page up
		if ev.Matches('u', vaxis.ModCtrl) {
			if r.cursor > 10 {
				r.cursor -= 10
			} else {
				r.cursor = 0
			}
		}
		// Ctrl-d : half a page down
		if ev.Matches('d', vaxis.ModCtrl) {
			if r.cursor < len(r.Filters)-11 {
				r.cursor += 10
			} else {
				r.cursor = len(r.Filters) - 1
			}
		}
		// j : down
		if ev.Matches('j') || ev.Matches(vaxis.KeyDown) {
			if r.cursor < len(r.Filters)-1 {
				r.cursor += 1
			}
		}
		// G : go to bottom
		if ev.Matches('G') || ev.Matches(vaxis.KeyEnd) {
			r.cursor = len(r.Filters) - 1
		}
		// k : up
		if ev.Matches('k') || ev.Matches(vaxis.KeyUp) {
			if r.cursor > 0 {
				r.cursor -= 1
			}
		}
		// g : go to top
		if ev.Matches('g') || ev.Matches(vaxis.KeyHome) {
			r.cursor = 0
		}
		// Ctrl-p : move filter up
		if ev.Matches('p', vaxis.ModCtrl) {
			if r.cursor > 0 {
				s := reflect.Swapper(r.Filters)
				s(r.cursor-1, r.cursor)
				r.cursor -= 1
			}
		}
		// Ctrl-n : move filter down
		if ev.Matches('n', vaxis.ModCtrl) {
			if r.cursor < len(r.Filters)-1 {
				s := reflect.Swapper(r.Filters)
				s(r.cursor, r.cursor+1)
				r.cursor += 1
			}
		}
		// Ctrl-t : move filter to top
		if ev.Matches('t', vaxis.ModCtrl) {
			s := reflect.Swapper(r.Filters)
			for p := r.cursor; p > 0; p -= 1 {
				s(p, p-1)
			}
			r.cursor = 0
		}
		// Ctrl-b : move filter to bottom
		if ev.Matches('b', vaxis.ModCtrl) {
			s := reflect.Swapper(r.Filters)
			for p := r.cursor; p < len(r.Filters)-1; p += 1 {
				s(p, p+1)
			}
		}
		// Enter : show matches for current filter in bottom pane
		if ev.Matches(vaxis.KeyEnter) {
			items = r.Filters[r.cursor].Matches[:]
			r.list.SetCursor(0)
		}
		// ] : select next match
		if ev.Matches(']') {
			cursor := r.Filters[r.cursor].Cursor + 1
			if cursor >= len(r.Filters[r.cursor].Matches) {
				cursor = -1
			}
			r.Filters[r.cursor].Cursor = cursor
		}
		// [ : select previous match
		if ev.Matches('[') {
			cursor := r.Filters[r.cursor].Cursor - 1
			if cursor < -1 {
				cursor = len(r.Filters[r.cursor].Matches) - 1
			}
			r.Filters[r.cursor].Cursor = cursor
		}
		// action on current filter
		if ev.Matches(' ') {
			// create tag + query pairs for each filter up-to and including the cursor
			constraints := make([]remote.Constraint, r.cursor+1)
			for i := 0; i <= r.cursor; i++ {
				f := r.Filters[i]
				if f.Cursor >= 0 && len(f.Matches) > f.Cursor {
					// exact match of value at cursor
					constraints = append(constraints, remote.Constraint{
						Tag:   f.Label,
						Op:    "==",
						Query: f.Matches[f.Cursor],
					})
				} else {
					// search for value using 'contains'
					constraints = append(constraints, remote.Constraint{
						Tag:   f.Label,
						Op:    "contains",
						Query: f.Value,
					})
				}
			}

			// fire-off a query for the current Filter
			filter := r.Filters[r.cursor]
			filter.Current_query = r.remote.PostQuery(filter.Label, constraints)
		}
		// Tab : focus on bottom panel
		if ev.Matches(vaxis.KeyTab) {
			return vxfw.FocusWidgetCmd(&r.list), nil
		}
		// Esc : focus on top panel
		if ev.Matches(vaxis.KeyEsc) {
			return vxfw.FocusWidgetCmd(r), nil
		}
	}
	for pos, filter := range r.Filters {
		filter.Active = (pos == r.cursor)
	}

	return vxfw.RedrawCmd{}, nil
}

func (w *Search) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	s := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, w)

	panel_size := vxfw.Size{Width: ctx.Max.Width, Height: 1}

	// filters
	panel_size.Height = 1
	for pos, filter := range w.Filters {
		surf1, err := filter.Draw(vxfw.DrawContext{Min: panel_size, Max: panel_size, Characters: ctx.Characters})
		if err != nil {
			return s, err
		}
		s.AddChild(0, pos, surf1)
	}


	// full item list
	panel_size.Height = ctx.Max.Height - uint16(len(w.Filters)) - 1
	surf2, err := w.list.Draw(vxfw.DrawContext{Min: panel_size, Max: panel_size, Characters: ctx.Characters})
	if err != nil {
		return s, err
	}
	s.AddChild(0, len(w.Filters), surf2)

	return s, nil
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


// Verify we meet the Widget interface
var _ vxfw.Widget = &Search{}
