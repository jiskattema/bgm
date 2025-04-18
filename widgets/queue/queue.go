package queue

import (
	"ash/bgm/remote"
	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/list"
	"git.sr.ht/~rockorager/vaxis/vxfw/text"
)

var items []string

type Queue struct {
	remote *remote.Remote
	list   list.Dynamic
}

func New(remote remote.Remote) *Queue {
	return &Queue{
		list: list.Dynamic{
			Builder:              getWidget,
			DrawCursor:           true,
			Gap:                  0,
			DisableEventHandlers: false,
		},
	}
}

func (q *Queue) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	switch ev := ev.(type) {
	case vaxis.FocusIn:
		return vxfw.FocusWidgetCmd(&q.list), nil
	case vaxis.Key:
		if ev.Matches('q') {
			print("Q pressed")
		}
	}
	return nil, nil
}

func (q *Queue) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	surf, err := q.list.Draw(ctx)
	return surf, err
}

func getWidget(i uint, cursor uint) vxfw.Widget {
	var style vaxis.Style
	if i >= 0 && i < 10 {
		return &text.Text{
			Content: "We have a widget",
			Style:   style,
		}
	}
	if i >= uint(len(items)) {
		return nil
	}
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
