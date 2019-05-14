package activity

import (
	"encoding/json"
	"strconv"
	"time"
)

type Verb int


type BaseActivty struct {
	Actor string `json:"actor"`
	Verb Verb `json:"verb"`
	Object string `json:"object"`
	Target string `json:"target"`
	Time time.Time `json:"time"`
	Extra string `json:"extra"`
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

func (self *BaseActivty)JsonSerialize()([]byte,error){
	return json.Marshal(self)
}