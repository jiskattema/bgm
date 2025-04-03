package root

import (
	"log"
	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/textfield"
	"ash/bgm/remote"
	"ash/bgm/widgets/search"
)

type Root struct {
	// The content of the widget
	input   textfield.TextField
	search  search.Search
	remote  remote.Remote
	app     *vxfw.App 

	routes   map[string]vxfw.Widget
}

func New(remote remote.Remote, app *vxfw.App) *Root {
	root := &Root{
		remote:  remote,
		app: app,
		routes: make(map[string]vxfw.Widget),
	}

	root.search.Init()
	root.addRoute("search", root.search.Entry)
	root.addRoute("command", &root.input)

	return root
}

func (r *Root) Navigate(route string) (vxfw.Command, error) {
  w, ok := r.routes[route]
  if ok {
	print("navigating to", route)
	return []vxfw.Command{
	  	vxfw.FocusWidgetCmd(w),
		vxfw.RedrawCmd{},
	  }, nil 
  } else {
	log.Fatalf("Route not defined: %v", route)
	return nil, nil
  }
}

func (r *Root) addRoute(route string, widget vxfw.Widget) {
  r.routes[route] = widget
}

func (r *Root) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	switch ev := ev.(type) {
	case vaxis.Key:
		// 1 : quit
		if ev.Matches('1') {
			return r.Navigate("search")
		}
		// Ctrl-C : quit
		if ev.Matches('c', vaxis.ModCtrl) {
			return vxfw.QuitCmd{}, nil
		}
		// / : search
		if ev.Matches(':') {
			// Prepare commandline
			r.input.Reset()
			r.input.InsertStringAtCursor("hello world")
			r.input.OnSubmit = func(line string) (vxfw.Command, error) {
				return r.Navigate("search")
			}
			// Focus the input widget
			return r.Navigate("command")
		}
		return vxfw.RedrawCmd{}, nil
	}
	return nil, nil
}

func (r *Root) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	s := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, r)

	// Add the search panel
	{
		subcontext := vxfw.DrawContext{
			Min: vxfw.Size{Width: ctx.Max.Width, Height: ctx.Max.Height - 1},
			Max: vxfw.Size{Width: ctx.Max.Width, Height: ctx.Max.Height - 1},
			Characters: ctx.Characters,
		}
		surf, err := r.search.Draw(subcontext)
		if err != nil {
			return s, err
		}
		s.AddChild(0, 0, surf)
	}

	// Add the commandline at the bottom
	{
		subcontext := vxfw.DrawContext{
			Min: vxfw.Size{Width: ctx.Max.Width, Height: 1},
			Max: vxfw.Size{Width: ctx.Max.Width, Height: 1},
			Characters: ctx.Characters,
		}
		surf, err := r.input.Draw(subcontext)
		if err != nil {
			return s, err
		}
		s.AddChild(0, int(s.Size.Height) - 1, surf)
	}

	return s, nil
}

// Verify we meet the Widget interface
var _ vxfw.Widget = &Root{}
