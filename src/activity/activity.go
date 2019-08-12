package activity

import (
	"encoding/json"
	"strconv"
	"time"
)

type Verb int


const(
	ActObject_Type_Body = 1
	ActObject_Type_Thing = 2
)

type  ActObject struct{
	Id int `json:id`
	Type int `json:type`
}

type BaseActivty struct {
	Actor  int `json:"actor"`
	Verb   Verb      `json:"verb"`
	Object ActObject `json:"object"`
	Target int `json:"target"`
	Time   time.Time `json:"time"`
	Extra  string    `json:"extra"`
}

func NewActivity(actor int, verb Verb, object int, target int, time time.Time, extra string) *BaseActivty {
	act := new(BaseActivty)
	act.Actor = actor
	act.Verb = verb
	act.Object = object
	act.Target = target
	act.Time = time
	act.Extra = extra
	return act
}

func (self *BaseActivty) SerializeId() string {
	return strconv.Itoa(int(self.Time.Unix()))
}

func (self *BaseActivty) Serialize() ([]byte, error) {
	return json.Marshal(self)
}
