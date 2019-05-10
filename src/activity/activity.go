package activity

import (
	"github.com/inclee/gokit/src/util/stablity"
	"time"
)

type ActivityOptions = func(Opts *Options) string

type Elem struct {
	Value string
}
type Verb struct {
	verb int
}
type Options struct {
	Actor *Elem
	Verb *Verb
	Object *Elem
	Target *Elem
	Time time.Time
	Extra string
}

type BaseActivty struct {
	options *Options
}

func NewActivity(opts... ActivityOptions)*BaseActivty{
	act := new(BaseActivty)
	act.options = new(Options)
	for _,opt := range opts{
		opt(act.options)
	}
	stablity.AssertNil(act.options.Actor,"actor can't nil")
	stablity.AssertNil(act.options.Verb,"verb can't nil")
	stablity.AssertNil(act.options.Object,"object can't nil")
	return act
}

