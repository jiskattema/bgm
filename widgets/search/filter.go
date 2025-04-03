package search 

import (
	//"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
)

var InactiveFilter = vaxis.Style{
	Foreground: vaxis.HexColor(0x00ffff),
	Background: vaxis.HexColor(0xff0fff),
}

var ActiveFilter = vaxis.Style{
	Foreground: vaxis.HexColor(0x000000),
	Background: vaxis.HexColor(0xffffff),
}

type Filter struct {
	Label         string
	Value         string
	Active        bool
	Current_query int
	Matches       []string
	Cursor        int
}

// no-op for now
func (f *Filter) HandleEvent(ev vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	return nil, nil
}

func (f *Filter) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	var count_text string

	if len(f.Matches) == 1 {
		count_text = fmt.Sprintf("%s [only]", f.Matches[0])
	} else {
		if f.Cursor >= 0 && f.Cursor < len(f.Matches) {
			count_text = fmt.Sprintf("%s [%d of %d]", f.Matches[f.Cursor], f.Cursor, len(f.Matches))
		} else {
			count_text = strconv.Itoa(len(f.Matches))
		}
	}

	spaces := int(ctx.Max.Width) -
		utf8.RuneCountInString(count_text) -
		utf8.RuneCountInString(f.Value) -
		utf8.RuneCountInString(f.Label)
	full_text := f.Label + f.Value + strings.Repeat(" ", spaces) + count_text

	chars := ctx.Characters(full_text)
	cells := make([]vaxis.Cell, 0, len(chars))

	// style := InactiveFilter
	// if f.Active {
	// 	style = ActiveFilter
	// }

	var w int
	for _, char := range chars {
		cell := vaxis.Cell{
			Character: char,
//			Style:     style,
		}
		cells = append(cells, cell)
		w += char.Width
	}

	return vxfw.Surface{
		Size:     vxfw.Size{Width: uint16(w), Height: 1},
		Widget:   f,
		Cursor:   nil,
		Buffer:   cells,
		Children: []vxfw.SubSurface{},
	}, nil
}
