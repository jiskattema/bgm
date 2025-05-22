package remote

import (
	"ash/bgm/base"
)

type Remote interface {
	Dial()
	HangUp()
	PostQuery(tag string, constraints []base.Constraint) int
}
