package root

import (
	"ash/bgm/base"
	"ash/bgm/remote"
	"ash/bgm/widgets/queue"
	"ash/bgm/widgets/search"
	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/text"
	"git.sr.ht/~rockorager/vaxis/vxfw/textfield"
)

type BgmState int

const (
	Searching BgmState = iota
	Queueing
)

type Root struct {
	// The content of the widget
	message      *text.Text
	input        *textfield.TextField
	search       *search.Search
	queue        *queue.Queue
	app          *vxfw.App
	state        BgmState
	messageModal bool
}

func New(remote remote.Remote, app *vxfw.App) *Root {
	return &Root{
		search:       search.New(remote, app),
		queue:        queue.New(remote),
		app:          app,
		state:        Searching,
		messageModal: false,
		message:      text.New("Prompt me"),
		input:        textfield.New(),
	}
}

func (r Root) currentWidget() vxfw.Widget {
	switch r.state {
	case Searching:
		return r.search
	case Queueing:
		return r.queue
	}
	return &r
}

func (r *Root) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	switch ev := ev.(type) {
	case base.Question:
		return r.ask(ev)
	case base.Playlist:
		r.queue.UpdateFromRemote(ev)
		return vxfw.RedrawCmd{}, nil
	case vaxis.Key:
		// 1 : Search panel
		if ev.Matches('1') {
			r.state = Searching
			return vxfw.FocusWidgetCmd(r.currentWidget()), nil
		}
		// 2 : Play queue
		if ev.Matches('2') {
			r.state = Queueing
			return vxfw.FocusWidgetCmd(r.currentWidget()), nil
		}
		// Ctrl-C : quit
		if ev.Matches('c', vaxis.ModCtrl) {
			return vxfw.QuitCmd{}, nil
		}
		// : : command
		if ev.Matches(':') {
			return r.ask(base.Question{
				Prompt: "command",
				Key: "command",
			})
		}
	}
	return nil, nil
}

func (r *Root) ask(what base.Question) (vxfw.Command, error) {
	// configure the messageModal
	r.messageModal = true
	r.message.Content = what.Prompt
	r.input.Reset()
	r.input.InsertStringAtCursor(what.Value)

	// Set callback
	r.input.OnSubmit = func(answer string) (vxfw.Command, error) {
		if what.Callback != nil {
			what.Callback(what.Key, answer)
		}
		r.messageModal = false
		return vxfw.FocusWidgetCmd(r.currentWidget()), nil
	}

	// Focus the messageModal 
	return vxfw.FocusWidgetCmd(r.input), nil
}

func (r *Root) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	s := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, r)

	// Add the main widget
	{
		subcontext := vxfw.DrawContext{
			Min:        vxfw.Size{Width: ctx.Max.Width, Height: ctx.Max.Height - 1},
			Max:        vxfw.Size{Width: ctx.Max.Width, Height: ctx.Max.Height - 1},
			Characters: ctx.Characters,
		}
		mainWidget := r.currentWidget()
		surf, err := mainWidget.Draw(subcontext)
		if err != nil {
			return s, err
		}
		s.AddChild(0, 0, surf)
	}

	// Add the messageModal
	if r.messageModal {
		subcontext := vxfw.DrawContext{
			Min:        vxfw.Size{Width: ctx.Max.Width, Height: 1},
			Max:        vxfw.Size{Width: ctx.Max.Width, Height: 1},
			Characters: ctx.Characters,
		}
		messagesurf, err := r.message.Draw(subcontext)
		if err != nil {
			return s, err
		}
		s.AddChild(0, int(s.Size.Height)-2, messagesurf)

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
