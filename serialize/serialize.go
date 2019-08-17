package serialize

type JsonAble interface {
	JsonSerialize()([]byte,error)
}