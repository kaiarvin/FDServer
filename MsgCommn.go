// MsgCommn
package ChatServer

import (
	"encoding/json"
	"reflect"
)

const (
	E_NONE = iota
	E_HEARTBEAT
)

func MsgJsonEncode(msg interface{}) ([]byte, bool) {
	buf, err := json.Marshal(msg)
	if err != nil {
		return nil, true
	}

	Length := len(buf)
	v := reflect.ValueOf(msg).Elem().Field(0).FieldByName("Length")
	if !v.CanSet() {
		return nil, true
	}

	v.SetInt(int64(Length))
	buf, err = json.Marshal(msg)
	if err != nil {
		return nil, true
	}
	return buf, false
}

func MsgJsonDecode(data []byte) ([]byte, *MsgHead) {
	var Head *MsgHead
	err := json.Unmarshal(data, &Head)
	if err != nil {
		return nil, nil
	}

	cpzero := Head
	cpzero.Length = 0
	cpstr, err := json.Marshal(Head)
	cpzerostr, err := json.Marshal(cpzero)
	lensub := len(cpstr) - len(cpzerostr)

	if Head.Length != int64(len(data)-lensub) {
		return nil, nil
	}

	return data, Head
}

type MsgHead struct {
	Id     int
	Length int64
}

type TestMsg struct {
	MsgHead
	A int
	B string
}
