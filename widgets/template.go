package main

import (
	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
)


type Root struct {
	// The content of the widget
}

func New(content string) *Root{
	return &Root{}
}

// Noop for text
func (r *Root) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	return nil, nil
}

func (r *Root) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	s := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, r)
	return s, nil
}

// Verify we meet the Widget interface
var _ vxfw.Widget = &Root{}
