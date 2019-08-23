package activity

import (
	"encoding/json"
	"strconv"
	"time"
)

type Verb int

type ActObject struct {
	Id   int `json:id`
	Type int `json:type`
}

type BaseActivty struct {
	Actor  int       `json:"actor"`
	Verb   Verb      `json:"verb"`
	Object ActObject `json:"object"`
	Target int       `json:"target"`
	Time   time.Time `json:"time"`
	Private bool 	 `json:"private"`
	Allow  []int     `json:"allow"`
	Deny   []int     `json:deny`
	Extra  string    `json:"extra"`
}

func NewActivity(actor int, verb Verb, object ActObject, target int, time time.Time,private bool,allow []int,deny []int, extra string) *BaseActivty {
	act := new(BaseActivty)
	act.Actor = actor
	act.Verb = verb
	act.Object = object
	act.Target = target
	act.Time = time
	act.Private = private
	act.Allow = allow
	act.Deny = deny
	act.Extra = extra
	return act
}

func (self *BaseActivty) SerializeId() string {
	return strconv.Itoa(int(self.Time.Unix()))
}

func (self *BaseActivty) Serialize() ([]byte, error) {
	return json.Marshal(self)
}

func  (self *BaseActivty) DeepCopy() (*BaseActivty) {
	return &BaseActivty{
		Actor: self.Actor,
		Verb: self.Verb,
		Object: self.Object,
		Target: self.Target,
		Time : self.Time,
		Private: self.Private,
		Allow: self.Allow,
		Deny: self.Deny,
		Extra: self.Extra,
	}
}
