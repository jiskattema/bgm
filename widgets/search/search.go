package search

import (
	// "reflect"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/list"
	"git.sr.ht/~rockorager/vaxis/vxfw/text"

	"ash/bgm/remote"
)

var items []string

type Search struct {
	Entry  vxfw.Widget 
	// The content of the widget
	list    list.Dynamic
	cursor  int
	remote  remote.Remote
}

func (r *Search) Init() {
	r.list = list.Dynamic{
		Builder:              getWidget,
		DrawCursor:           false,
		Gap:                  0,
		// DisableEventHandlers: true,
	}
	r.Entry = &r.list
}

func (r *Search) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	switch ev := ev.(type) {
	case vaxis.Key:
		// Ctrl-p : move filter up
		if ev.Matches('p', vaxis.ModCtrl) {
			if r.cursor > 0 {
				// s := reflect.Swapper()
				// s(r.cursor-1, r.cursor)
				r.cursor -= 1
			}
		}
		// // Ctrl-n : move filter down
		// if ev.Matches('n', vaxis.ModCtrl) {
		// 	if r.cursor < len(r.Filters)-1 {
		// 		s := reflect.Swapper(r.Filters)
		// 		s(r.cursor, r.cursor+1)
		// 		r.cursor += 1
		// 	}
		// }
		// // Ctrl-t : move filter to top
		// if ev.Matches('t', vaxis.ModCtrl) {
		// 	s := reflect.Swapper(r.Filters)
		// 	for p := r.cursor; p > 0; p -= 1 {
		// 		s(p, p-1)
		// 	}
		// 	r.cursor = 0
		// }
		// // Ctrl-b : move filter to bottom
		// if ev.Matches('b', vaxis.ModCtrl) {
		// 	s := reflect.Swapper(r.Filters)
		// 	for p := r.cursor; p < len(r.Filters)-1; p += 1 {
		// 		s(p, p+1)
		// 	}
		// }
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
		// 	constraints := make([]remote.Constraint, r.cursor+1)
		// 	for i := 0; i <= r.cursor; i++ {
		// 		f := r.Filters[i]
		// 		if f.Cursor >= 0 && len(f.Matches) > f.Cursor {
		// 			// exact match of value at cursor
		// 			constraints = append(constraints, remote.Constraint{
		// 				Tag:   f.Label,
		// 				Op:    "==",
		// 				Query: f.Matches[f.Cursor],
		// 			})
		// 		} else {
		// 			// search for value using 'contains'
		// 			constraints = append(constraints, remote.Constraint{
		// 				Tag:   f.Label,
		// 				Op:    "contains",
		// 				Query: f.Value,
		// 			})
		// 		}
		// 	}

		// 	// fire-off a query for the current Filter
		// 	filter := r.Filters[r.cursor]
		// 	filter.Current_query = r.remote.PostQuery(filter.Label, constraints)
		// }
	}

	return vxfw.RedrawCmd{}, nil
}

func (w *Search) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	surf := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, w)

	subsurf, err := w.list.Draw(ctx)
	if err != nil {
		print("Error!")
		return surf, err
	}
	surf.AddChild(0, 0, subsurf)

	return surf, nil
}

func getWidget(i uint, cursor uint) vxfw.Widget {
    if i>0 && i < 10 {
	  return &text.Text{
		  Content: "Here we go again",
	  }
	}
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
