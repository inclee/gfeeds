package activity

import (
	"encoding/json"
	"strconv"
	"time"
)

type Verb int

type Object struct {
	Id   int `json:"id"`
	Type int `json:"type"`
}

func (o Object) Equal(ot Object) bool {
	return o.Id == ot.Id && o.Type == ot.Type
}

type BaseActivty struct {
	Actor     int           `json:"actor"`
	Verb      Verb          `json:"verb"`
	VerbObj   Object        `json:"verbobj"`
	Target    int           `json:"target"`
	TargetObj Object        `json:"targetobj"`
	Time      time.Time     `json:"time"`
	Private   bool          `json:"private"`
	Allow     []int         `json:"allow"`
	Deny      []int         `json:"deny"`
	Extra     string        `json:"extra"`
	Context   ActiveContext `json:"context"`
}

type ActiveContext struct {
	RootId   int `json:"root_id"`
	RootType int `json:"root_type"`
}

func NewActivity(actor int, verb Verb, verbobj Object, target int, targetobj Object, time time.Time, private bool, allow []int, deny []int, extra string, context ActiveContext) *BaseActivty {
	act := new(BaseActivty)
	act.Actor = actor
	act.Verb = verb
	act.VerbObj = verbobj
	act.Target = target
	act.TargetObj = targetobj
	act.Time = time
	act.Private = private
	act.Allow = allow
	act.Deny = deny
	act.Extra = extra
	act.Context = context
	return act
}

func (self *BaseActivty) SerializeId() string {
	return strconv.Itoa(int(self.Time.Unix()))
}

func (self *BaseActivty) Serialize() ([]byte, error) {
	return json.Marshal(self)
}

func (self *BaseActivty) DeepCopy() *BaseActivty {
	return &BaseActivty{
		Actor:     self.Actor,
		Verb:      self.Verb,
		VerbObj:   self.VerbObj,
		Target:    self.Target,
		TargetObj: self.TargetObj,
		Time:      self.Time,
		Private:   self.Private,
		Allow:     self.Allow,
		Deny:      self.Deny,
		Context:   self.Context,
		Extra:     self.Extra,
	}
}
