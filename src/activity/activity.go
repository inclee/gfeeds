package activity

import (
	"strconv"
	"time"
)

type Verb int

type BaseActivty struct {
	Actor string
	Verb Verb
	Object string
	Target string
	Time time.Time
	Extra string
}

func NewActivity(actor string,verb Verb,object string,target string,time time.Time ,extra string)*BaseActivty{
	act := new(BaseActivty)
	act.Actor  = actor
	act.Verb = verb
	act.Object = object
	act.Time = time
	act.Extra = extra
	return act
}

func (self *BaseActivty)SerializeId()string{
	return strconv.Itoa(int(self.Time.Unix()))
}