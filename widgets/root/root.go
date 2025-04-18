package root

import (
	"ash/bgm/remote"
	"ash/bgm/widgets/queue"
	"ash/bgm/widgets/search"
	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/textfield"
)

type BgmState int

const (
	Searching BgmState = iota
	Queueing
)

type Root struct {
	// The content of the widget
	input  *textfield.TextField
	search *search.Search
	queue  *queue.Queue
	remote remote.Remote
	app    *vxfw.App
	state  BgmState
}

func New(remote remote.Remote, app *vxfw.App) *Root {
	return &Root{
		input:  textfield.New(),
		search: search.New(remote),
		queue:  queue.New(remote),
		remote: remote,
		app:    app,
		state:  Searching,
	}
}

func (r *Root) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	switch ev := ev.(type) {
	case vaxis.Key:
		// 1 : Search panel
		if ev.Matches('1') {
			r.state = Searching
			return vxfw.FocusWidgetCmd(r.search), nil
		}
		// 2 : Play queue
		if ev.Matches('2') {
			r.state = Queueing
			return vxfw.FocusWidgetCmd(r.queue), nil
		}
		// Ctrl-C : quit
		if ev.Matches('c', vaxis.ModCtrl) {
			return vxfw.QuitCmd{}, nil
		}
		// / : search
		if ev.Matches(':') {
			// Set callback
			r.input.OnSubmit = func(line string) (vxfw.Command, error) {
				return vxfw.FocusWidgetCmd(r), nil
			}
			// Focus the input widget
			r.input.Reset()
			r.input.InsertStringAtCursor("hello world")
			return vxfw.FocusWidgetCmd(r.input), nil
		}
	}
	return nil, nil
}

func (r *Root) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	var mainWidget vxfw.Widget

	switch r.state {
	case Searching:
		mainWidget = r.search
	case Queueing:
		mainWidget = r.queue
	}

	s := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, r)

	// Add the main widget
	{
		subcontext := vxfw.DrawContext{
			Min:        vxfw.Size{Width: ctx.Max.Width, Height: ctx.Max.Height - 1},
			Max:        vxfw.Size{Width: ctx.Max.Width, Height: ctx.Max.Height - 1},
			Characters: ctx.Characters,
		}
		surf, err := mainWidget.Draw(subcontext)
		if err != nil {
			return s, err
		}
		s.AddChild(0, 0, surf)
	}

	// Add the commandline at the bottom
	{
		subcontext := vxfw.DrawContext{
			Min:        vxfw.Size{Width: ctx.Max.Width, Height: 1},
			Max:        vxfw.Size{Width: ctx.Max.Width, Height: 1},
			Characters: ctx.Characters,
		}
		surf, err := r.input.Draw(subcontext)
		if err != nil {
			return s, err
		}
		s.AddChild(0, int(s.Size.Height)-1, surf)
	}

	return s, nil
}

// Verify we meet the Widget interface
var _ vxfw.Widget = &Root{}
