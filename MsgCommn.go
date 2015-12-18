// MsgCommn
package FDServer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

const (
	E_REQ_NONE = iota
	E_REQ_PING
	E_REQ_REGIST
	E_REQ_ENTERSERVER
	E_REQ_EXITSERVER
	E_REQ_CHATDATA
)

const (
	E_ACK_NONE = iota
	E_ACK_PONG
	E_ACK_REGIST
	E_ACK_EXIT
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

	cpzero := *Head
	cpzero.Length = 0
	cpstr, err := json.Marshal(Head)
	cpzerostr, err := json.Marshal(cpzero)
	lensub := len(cpstr) - len(cpzerostr)
	if Head.Length != int64(len(data)-lensub) {
		fmt.Println("Head.Length:", Head.Length, "len:", int64(len(data)-lensub))
		return nil, nil
	}

	return data, Head
}

func MsgByteToJson(buf []byte, msg interface{}) {
	err := json.Unmarshal(buf, msg)
	if err != nil {
		return
	}
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

type HeartBeat struct {
	MsgHead
	RecTime time.Time
	AckTime time.Time
}

type ReqUserRegist struct {
	MsgHead
	Uname string
	Pw    string
	Gname string
}

type AckUserRegist struct {
	MsgHead
	Result int
}

type ReqUserEnterServer struct {
	MsgHead
	Uname string
	Pw    string
}

type AckUserEnterServer struct {
	MsgHead
	UserList map[int]string
	Gname    string
	Level    int
	Sex      int8
}

type ReqUserLogout struct {
	MsgHead
	AccountID int
	Name      string
}

type ReqChatData struct {
	MsgHead
	FromAccountID int
	ToAccountID   int
	FromName      string
	ToName        string
	Data          string
}
