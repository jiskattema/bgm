package search

import (
	"reflect"
	// "slices"
	"ash/bgm/base"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/list"
	"git.sr.ht/~rockorager/vaxis/vxfw/text"

	"ash/bgm/remote"
)

var songs []base.Attrs

type SearchTag struct {
	Display string
	Name    string
	Value   string
}

var filters []SearchTag

type Search struct {
	filters list.Dynamic
	songs   list.Dynamic
	remote  remote.Remote
	app     *vxfw.App
}

func New(remote remote.Remote, app *vxfw.App) *Search {
	filters = make([]SearchTag, 0, len(base.ServerTags))
	for _, tag := range base.ServerTags {
		if tag.Visible {
			filters = append(filters, SearchTag{
				Display: tag.Display,
				Name:    tag.Name,
			})
		}
	}

	return &Search{
		filters: list.Dynamic{
			Builder:              formatFilter,
			DrawCursor:           true,
			Gap:                  0,
			DisableEventHandlers: false,
		},
		songs: list.Dynamic{
			Builder:              formatSong,
			DrawCursor:           false,
			Gap:                  0,
			DisableEventHandlers: false,
		},
		remote: remote,
		app:    app,
	}
}

func setFilter(q, a string) {
	for i := range len(filters) {
		if filters[i].Name == q {
			filters[i].Value = a
		}
	}
}

// Noop for text
func (r *Search) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	switch ev := ev.(type) {
	case vaxis.FocusIn:
		return vxfw.FocusWidgetCmd(&r.filters), nil
	case vaxis.Key:
		// Enter : set filter
		if ev.Matches(vaxis.KeyEnter) {
			r.app.PostEvent(base.Question{
				Prompt:   "Enter " + filters[r.filters.Cursor()].Display,
				Key:      filters[r.filters.Cursor()].Name,
				Value:    filters[r.filters.Cursor()].Value,
				Callback: setFilter,
			})
		}

		// // Ctrl-u : half a page up
		// if ev.Matches('u', vaxis.ModCtrl) {
		// 	if r.cursor > 10 {
		// 		r.cursor -= 10
		// 	} else {
		// 		r.cursor = 0
		// 	}
		// }
		// // Ctrl-d : half a page down
		// if ev.Matches('d', vaxis.ModCtrl) {
		// 	if r.cursor < len(r.Filters)-11 {
		// 		r.cursor += 10
		// 	} else {
		// 		r.cursor = len(r.Filters) - 1
		// 	}
		// }
		// // G : go to bottom
		// if ev.Matches('G') || ev.Matches(vaxis.KeyEnd) {
		// 	r.cursor = len(r.Filters) - 1
		// }

		// Ctrl-p : move filter up
		if ev.Matches('p', vaxis.ModCtrl) {
			c := int(r.filters.Cursor())
			if c > 0 {
				s := reflect.Swapper(filters)
				s(c-1, c)
				r.filters.SetCursor(uint(c - 1))
			}
		}
		// Ctrl-n : move filter down
		if ev.Matches('n', vaxis.ModCtrl) {
			c := int(r.filters.Cursor())
			if c < len(filters)-1 {
				s := reflect.Swapper(filters)
				s(c, c+1)
				r.filters.SetCursor(uint(c + 1))
			}
		}
		// Ctrl-t : move filter to top
		if ev.Matches('t', vaxis.ModCtrl) {
			c := int(r.filters.Cursor())
			s := reflect.Swapper(filters)
			s(c, 0)
			r.filters.SetCursor(0)
		}
		// Ctrl-b : move filter to bottom
		if ev.Matches('b', vaxis.ModCtrl) {
			c := int(r.filters.Cursor())
			s := reflect.Swapper(filters)
			for c < len(filters)-1 {
				s(c, c+1)
				c = c + 1
			}
		}
		// // Enter : show matches for current filter in bottom pane
		// if ev.Matches(vaxis.KeyEnter) {
		// 	items = r.Filters[r.cursor].Matches[:]
		// 	r.list.SetCursor(0)
		// }
		// // ] : select next match
		// if ev.Matches(']') {
		// 	cursor := r.Filters[r.cursor].Cursor + 1
		// 	if cursor >= len(r.Filters[r.cursor].Matches) {
		// 		cursor = -1
		// 	}
		// 	r.Filters[r.cursor].Cursor = cursor
		// }
		// // [ : select previous match
		// if ev.Matches('[') {
		// 	cursor := r.Filters[r.cursor].Cursor - 1
		// 	if cursor < -1 {
		// 		cursor = len(r.Filters[r.cursor].Matches) - 1
		// 	}
		// 	r.Filters[r.cursor].Cursor = cursor
		// }
		// // action on current filter
		// if ev.Matches(' ') {
		// 	// create tag + query pairs for each filter up-to and including the cursor
		// 	constraints := make([]base.Constraint, r.cursor+1)
		// 	for i := 0; i <= r.cursor; i++ {
		// 		f := r.Filters[i]
		// 		if f.Cursor >= 0 && len(f.Matches) > f.Cursor {
		// 			// exact match of value at cursor
		// 			constraints = append(constraints, base.Constraint{
		// 				Tag:   f.Label,
		// 				Op:    "==",
		// 				Value: f.Matches[f.Cursor],
		// 			})
		// 		} else {
		// 			// search for value using 'contains'
		// 			constraints = append(constraints, base.Constraint{
		// 				Tag:   f.Label,
		// 				Op:    "contains",
		// 				Value: f.Value,
		// 			})
		// 		}
		// 	}

		// 	// fire-off a query for the current Filter
		// 	filter := r.Filters[r.cursor]
		// 	filter.Current_query = r.remote.PostQuery(filter.Label, constraints)
		// }

		// Tab : focus on bottom panel
		if ev.Matches(vaxis.KeyTab) {
			return vxfw.FocusWidgetCmd(&r.songs), nil
		}
		// Esc : focus on top panel
		if ev.Matches(vaxis.KeyEsc) {
			return vxfw.FocusWidgetCmd(&r.filters), nil
		}
	}

	return vxfw.RedrawCmd{}, nil
}

func (w *Search) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	s := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, w)

	panel_size := vxfw.Size{Width: ctx.Max.Width, Height: 1}

	// filters
	panel_size.Height = uint16(len(filters))
	surf1, err := w.filters.Draw(vxfw.DrawContext{Min: panel_size, Max: panel_size, Characters: ctx.Characters})
	if err != nil {
		return s, err
	}
	s.AddChild(0, 0, surf1)

	// songs
	panel_size.Height = ctx.Max.Height - uint16(len(filters)) - 1
	surf2, err := w.songs.Draw(vxfw.DrawContext{Min: panel_size, Max: panel_size, Characters: ctx.Characters})
	if err != nil {
		return s, err
	}
	s.AddChild(0, len(filters), surf2)

	return s, nil
}

func formatFilter(i uint, cursor uint) vxfw.Widget {
	if i >= uint(len(filters)) {
		return nil
	}

	var style vaxis.Style
	if i == cursor {
		// style.Attribute = vaxis.AttrReverse
		style.Foreground = vaxis.HexColor(0x000000)
		style.Background = vaxis.HexColor(0xffffff)
	} else {
		style.Foreground = vaxis.HexColor(0x00ffff)
		style.Background = vaxis.HexColor(0xff0fff)
	}

	return &text.Text{
		Content: filters[i].Display + " " + filters[i].Value,
		Style:   style,
	}
}

func formatSong(i uint, cursor uint) vxfw.Widget {
	var style vaxis.Style
	if i >= 0 && i < 10 {
		return &text.Text{
			Content: "We have a widget",
			Style:   style,
		}
	}
	if i >= uint(len(songs)) {
		return nil
	}
	if i == cursor {
		style.Attribute = vaxis.AttrReverse
	}

	var display_text string
	if title, ok := songs[i]["Title"]; ok {
		display_text = title
	} else {
		display_text = "[Unknown]"
	}

	return &text.Text{
		Content: display_text,
		Style:   style,
	}
}

// Verify we meet the Widget interface
var _ vxfw.Widget = &Search{}
