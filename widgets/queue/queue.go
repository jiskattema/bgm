package queue

import (
	"ash/bgm/base"
	"ash/bgm/remote"
	"fmt"
	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/list"
	"git.sr.ht/~rockorager/vaxis/vxfw/text"
)

// MPD uses https://mpd.readthedocs.io/en/latest/protocol.html#tags
var items base.Playlist

type Queue struct {
	remote remote.Remote
	list   list.Dynamic
}

func New(remote remote.Remote) *Queue {
	return &Queue{
		remote: remote,
		list: list.Dynamic{
			Builder:              getWidget,
			DrawCursor:           true,
			Gap:                  0,
			DisableEventHandlers: false,
		},
	}
}

func (q *Queue) UpdateFromRemote(newQueue base.Playlist) {
	items = newQueue
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
	if i >= uint(len(items)) {
		return nil
	}
	if i == cursor {
		style.Attribute = vaxis.AttrReverse
	}
	display_text := fmt.Sprintf("%2s ~ %10s ~ %30s", items[i]["Track"], items[i]["Artist"], items[i]["Title"])

	return &text.Text{
		Content: display_text,
		Style:   style,
	}
}
