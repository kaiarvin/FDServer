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
	RecMsgByte    []byte
	AckMsgByte    []byte
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
	for {
		if this.Server.IsCloseServer || !this.IsLive {
			return
		}
		var pkgbody []byte
		this.RecMsgByte, pkgbody = MsgJsonDecode(this.RecMsgByte)

		if pkgbody == nil {
			return
		} else {
			this.RecMsg <- string(pkgbody)
		}

		select {
		case buff := <-this.RecMsg:
			{
				fmt.Println("In <-this.RecMsg...")
				data := []byte(buff)
				Head := new(MsgHead)
				MsgByteToJson(data, Head)

				fmt.Println("Head:", Head)
				this.ProcessMsg(data, Head)
				fmt.Println("Out <-this.RecMsg...")
			}
		case data := <-this.AckMsg:
			{
				fmt.Println("In <-this.AckMsg...")
				buf := []byte(data)
				n, err := this.Client_Socket.Write(buf)
				if err != nil {
					if err != io.EOF {
						break
					}
				}
				fmt.Println("Ack len: ", n)
				if n != 0 {
					this.Server.AckDataSize += uint64(n)
				}
				fmt.Println("Out <-this.AckMsg...")
			}
		default:
			{
			}
		}
	}
}

func (this *Client) SendDataToChann(buf []byte) {
	data := string(buf)
	fmt.Println("SendDataToChann:", data)
	this.AckMsg <- data
	//this.ClientMsgProcess()
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
	if conn == nil {
		return
	}
	cl := this.Server.UserList[conn]
	if cl == nil {
		return
	}
	fmt.Println("SendToOtherByAccount:", account, this.Server.AccountList, this.Server.UserList)
	buf, err := MsgJsonEncode(msg)
	if err {
		return
	}

	cl.SendDataToChann(buf)
}

func (this *Client) SendToOtherByName(name string, msg interface{}) {
	conn := this.Server.NameList[name]
	if conn == nil {
		return
	}

	cl := this.Server.UserList[conn]
	if cl == nil {
		return
	}
	buf, err := MsgJsonEncode(msg)
	if err {
		return
	}

	cl.SendDataToChann(buf)
}

func (this *Client) SendToAllPlayer(msg interface{}) {
	buf, err := MsgJsonEncode(msg)
	if err {
		fmt.Println("SendToSlef Error:", err)
		return
	}

	for _, v := range this.Server.UserList {
		v.SendDataToChann(buf)
	}
}

func (this *Client) SendToAllExceptMe(msg interface{}) {
	buf, err := MsgJsonEncode(msg)
	if err {
		fmt.Println("SendToSlef Error:", err)
		return
	}

	for _, v := range this.Server.UserList {
		if v.Name == this.Name {
			continue
		}
		v.SendDataToChann(buf)
	}
}

func (this *Client) ProcessMsg(data []byte, head *MsgHead) {
	//length 取出后吧Msg的length设置为0然后从新获取然后作比较
	switch head.Id {
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
			fmt.Print(this.Name, " Enter Server. ", "Resutl:", rel)
			ack := new(AckUserEnterServer)
			ack.Id = E_ACK_ENTERSERVER
			//this.SendToSlef(ack)
			if rel {
				this.Server.WritList(this)
				ack.AccountId = this.AccountID
				ack.Gname = this.Name
				ack.Level = this.Level
				ack.Sex = this.Sex
				this.SendToSlef(ack)

				acklist := new(AckUserNameList)
				acklist.Id = E_ACK_SENDUSERNAMELIST
				acklist.FromName = this.Name
				acklist.UserNameList = this.Server.NameAccountList
				this.SendToAllPlayer(acklist)
			} else {
				ack := new(AckUserEnterServer)
				ack.AccountId = 0
				this.SendToSlef(ack)
			}
		}
	case E_REQ_CHATDATA:
		{
			say := new(ReqChatData)
			MsgByteToJson(data, say)
			ack := new(AckChatData)
			ack.Id = E_ACK_CHATDATA
			ack.FromAccountID = say.FromAccountID
			ack.FromName = say.FromName
			ack.ToAccountID = say.ToAccountID
			ack.ToName = say.ToName
			ack.Data = say.Data
			if -1 == say.ToAccountID {
				this.SendToAllPlayer(ack)
				break
			}
			fmt.Println(ack.ToAccountID)
			if 0 != ack.ToAccountID {
				fmt.Println("this.SendToOtherByAccount(ack.ToAccountID, ack)", ack)
				this.SendToOtherByAccount(ack.ToAccountID, ack)
				break
			} else if "" != ack.ToName {
				fmt.Println("this.SendToOtherByName(ack.ToName, ack)", ack)
				this.SendToOtherByName(ack.ToName, ack)
				break
			}
			fmt.Println("Cant Find Player!")
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
