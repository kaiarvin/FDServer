// MsgCommn
package ChatServer

const (
	E_NONE = iota
	E_HEARTBEAT
)

type Msg struct {
	Id     int
	Length int64
}

type TestMsg struct {
	Msg
	A int
	B string
}
