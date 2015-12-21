// ClientMsgProcess
package FDServer

import (
	"fmt"
	"io"
	"net"
)

type Client struct {
	AccountID     int
	Name          string
	Sex           int8
	Level         int
	Area          int
	Viplevel      int
	IsLive        bool
	RecMsg        chan string
	AckMsg        chan string
	Server        *Server
	Client_Socket net.Conn
	UserRelation
}

type UserRelation struct {
	FriendInfo map[uint64]string
	IgnorID    map[uint64]string
}

func (this *Client) ClientMsgProcess() {
	fmt.Println("In ClientMsgProcess...")
	select {
	case buf := <-this.RecMsg:
		{
			fmt.Println("In <-this.RecMsg...")
			data := []byte(buf)

			buf, Head := MsgJsonDecode(data)
			if buf == nil || Head == nil {
				fmt.Print("RecMsg buf==", buf, "||", "Head==", Head)
				return
			}
			this.ProcessMsg(buf, Head)
			fmt.Println("Out <-this.RecMsg...")
		}
	case data := <-this.AckMsg:
		{
			fmt.Println("In <-this.AckMsg...")
			buf := []byte(data)
			n, err := this.Client_Socket.Write(buf)
			if err != nil {
				if err != io.EOF {
					return
				}
			}
			fmt.Println("Ack len: ", n)
			if n != 0 {
				this.Server.AckDataSize += uint64(n)
			}
			fmt.Println("Out <-this.AckMsg...")
		}
	}
}

func (this *Client) SendDataToChann(buf []byte) {
	data := string(buf)
	fmt.Println("SendDataToChann:", data)
	this.AckMsg <- data
	this.ClientMsgProcess()
}

func (this *Client) SendToSlef(msg interface{}) {
	buf, err := MsgJsonEncode(msg)
	if err {
		fmt.Println("SendToSlef Error:", err)
		return
	}
	this.SendDataToChann(buf)
}

func (this *Client) SendToOtherByAccount(account int, msg interface{}) {
	conn := this.Server.AccountList[account]
	cl := this.Server.UserList[conn]

	buf, err := MsgJsonEncode(msg)
	if err {
		return
	}

	cl.SendDataToChann(buf)
}

func (this *Client) SendToOtherByName(name string, msg interface{}) {
	conn := this.Server.NameList[name]
	cl := this.Server.UserList[conn]

	buf, err := MsgJsonEncode(msg)
	if err {
		return
	}

	cl.SendDataToChann(buf)
}

func (this *Client) ProcessMsg(data []byte, Head *MsgHead) {
	//length 取出后吧Msg的length设置为0然后从新获取然后作比较
	switch Head.Id {
	case E_REQ_NONE:
		{
			testMsg := new(TestMsg)
			MsgByteToJson(data, testMsg)
			fmt.Println(testMsg)

		}
	case E_REQ_PING:
		{
			heartBeat := new(HeartBeat)
			MsgByteToJson(data, heartBeat)
			fmt.Println(heartBeat)
		}
	case E_REQ_REGIST:
		{
			regist := new(ReqUserRegist)
			MsgByteToJson(data, regist)
			account := &Account{UserName: regist.Uname, UserPw: regist.Pw}
			rel := ProcessRegistAccountId(this.Server.DB, account, regist.Gname)

			ack := new(AckUserRegist)
			ack.Id = E_ACK_REGIST
			ack.Result = rel

			this.SendToSlef(ack)
		}
	case E_REQ_ENTERSERVER:
		{
			enter := new(ReqUserEnterServer)
			MsgByteToJson(data, enter)
			fmt.Println(enter)
			rel := this.EnterWorld(enter)
			fmt.Print("Enter Server:", enter, "Resutl:", rel)
			ack := new(AckUserEnterServer)
			ack.Id = E_ACK_ENTERSERVER
			if rel {
				this.Server.WritList(this)
				ack.AccountId = this.AccountID
				ack.Gname = this.Name
				ack.Level = this.Level
				ack.Sex = this.Sex
				ack.UserList = this.Server.NameAccountList
				this.SendToSlef(ack)
				fmt.Print("ACKENTER True:", ack)
			} else {
				ack := new(AckUserEnterServer)
				ack.AccountId = 0
				this.SendToSlef(ack)
				fmt.Print("ACKENTER False:", ack)
			}
		}
	case E_REQ_CHATDATA:
		{
			say := new(ReqChatData)
			MsgByteToJson(data, say)
			fmt.Println(say)
		}
	case E_REQ_EXITSERVER:
		this.CloseClient()
	default:
		{

		}
	}
}

func (this *Client) EnterWorld(data *ReqUserEnterServer) bool {
	account := CheckUserEnterWorld(this.Server.DB, data)
	if account == nil || account.Id == 0 {
		return false
	}

	if account.UserPw != data.Pw || account.UserName != data.Uname {
		return false
	}

	character := GetDBCharacter(this.Server.DB, account.Id)
	this.AccountID = character.Id
	this.Level = character.Level
	this.Name = character.Name
	this.Sex = character.Sex
	return true
}
func (this *Client) CloseClient() {
	//this.Server.CleanExitUser(&this.Client_Socket)
	this.IsLive = false
}
