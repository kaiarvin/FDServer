// MsgCommn
package FDServer

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"
)

const (
	Const_ServerHead_Len = int32(4)
	Const_ClientHead_Len = int32(4)
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
	E_ACK_ENTERSERVER
	E_ACK_EXIT
	E_ACK_CHATDATA
	E_ACK_SENDUSERNAMELIST
)

func MsgJsonEncode(msg interface{}) ([]byte, bool) {

	pkgBody, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("MsgJsonEncode Error:", err)
		return nil, true
	}

	headLen := int32(len(pkgBody))
	pkg := bytes.NewBuffer([]byte{})
	binary.Write(pkg, binary.BigEndian, headLen)
	binary.Write(pkg, binary.BigEndian, pkgBody)
	return pkg.Bytes(), false
}

func MsgJsonDecode(RecBytes []byte) ([]byte, []byte) {
	if 0 == len(RecBytes) {
		return nil, nil
	}
	fmt.Println("MsgJsonDecode")
	headData := RecBytes[:Const_ClientHead_Len]
	var bodyLen int32
	pkgHead := bytes.NewBuffer(headData)
	binary.Read(pkgHead, binary.BigEndian, &bodyLen)
	fmt.Println("MsgJsonDecode bodyLen:", bodyLen)

	if (Const_ClientHead_Len + bodyLen) > int32(len(RecBytes)) {
		fmt.Println("MsgJsonDecode Const_ClientHead_Len + bodyLen:", Const_ClientHead_Len+bodyLen)
		fmt.Println("MsgJsonDecode len(RecBytes):", len(RecBytes))

		return RecBytes, nil
	}

	pkgBody := RecBytes[Const_ClientHead_Len:(Const_ClientHead_Len + bodyLen)]
	fmt.Println("client.RecMsgByte:", string(RecBytes))
	fmt.Println("MsgJsonDecode pkgBody:", string(pkgBody))
	RecBytes = RecBytes[Const_ClientHead_Len+bodyLen:]

	return RecBytes, pkgBody
}

func MsgByteToJson(buf []byte, msg interface{}) {
	fmt.Println(string(buf))
	err := json.Unmarshal(buf, msg)
	if err != nil {
		return
	}
}

type MsgHead struct {
	Id int32
}

type TestMsg struct {
	MsgHead
	A int32
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
	AccountId int
	Gname     string
	Level     int
	Sex       int8
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

type AckChatData struct {
	MsgHead
	FromAccountID int
	ToAccountID   int
	FromName      string
	ToName        string
	Data          string
}

type AckUserNameList struct {
	MsgHead
	FromName     string
	UserNameList map[string]int
}
