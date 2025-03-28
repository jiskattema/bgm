package search

import (
	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/list"
	"git.sr.ht/~rockorager/vaxis/vxfw/textfield"
)


type Search struct {
	// The content of the widget
	Filters []*Filter
	list    list.Dynamic
	input   *textfield.TextField
}

func New(content string) *Search{
	return &Search{}
}

// Noop for text
func (r *Search) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	return nil, nil
}

func (w *Search) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	s := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, w)

	for pos, filter := range w.Filters {
		surf, err := filter.Draw(ctx)
		if err != nil {
			return s, err
		}
		s.AddChild(0, pos, surf)
	}

	panel_size := vxfw.Size{Width: ctx.Max.Width, Height: 1}

	// The commandline
	panel_size.Height = 1
	s, err := w.input.Draw(vxfw.DrawContext{Min: panel_size, Max: panel_size, Characters: ctx.Characters})
	if err != nil {
		return s, err
	}
	s.AddChild(0, int(ctx.Max.Height-1), s)

	// full item list
	panel_size.Height = ctx.Max.Height - uint16(len(w.Filters)) - 1
	s, err = w.list.Draw(vxfw.DrawContext{Min: panel_size, Max: panel_size, Characters: ctx.Characters})
	if err != nil {
		return s, err
	}
	s.AddChild(0, len(w.Filters), s)

	return s, nil
}

// Verify we meet the Widget interface
var _ vxfw.Widget = &Search{}
