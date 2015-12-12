// ClientMsgProcess
package ChatServer

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

type Client struct {
	AccountID     int
	Name          string
	Icon          int
	Level         int
	Area          int
	Viplevel      int
	IsLive        bool
	RecMsg        chan string
	AckMsg        chan string
	ChatServer    *ChatServer
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
			buffer := []byte(buf)
			var Head Msg
			err := json.Unmarshal(buffer, &Head)
			if err != nil {
				fmt.Println(err)
				return
			}
			this.ProcessMsg(buffer, &Head)

			fmt.Println("In <-this.RecMsg...")
			fmt.Println(buf)
			fmt.Println("Out <-this.RecMsg...")
		}
	case data := <-this.AckMsg:
		{
			fmt.Println("In <-this.AckMsg...")
			fmt.Println("data:", data)
			buf := []byte(data)
			n, err := this.Client_Socket.Write(buf)
			if err != nil {
				if err != io.EOF {
					return
				}
			}
			fmt.Println("Ack len: ", n)
			if n != 0 {
				this.ChatServer.AckDataSize += uint64(n)
			}
			fmt.Println("In <-this.AckMsg...")
		}
	}
}

func (this *Client) SendToSelfClient(buf []byte) {
	data := string(buf)
	this.AckMsg <- data
	this.ClientMsgProcess()
}

func (this *Client) MsgToJson(MsgData interface{}) ([]byte, int) {
	buf, err := json.Marshal(MsgData)
	if err != nil {
		fmt.Println(err)
		return nil, 0
	}
	length := len(buf)

	return buf, length
}

func (this *Client) ProcessMsg(data []byte, Head *Msg) {
	//length 取出后吧Msg的length设置为0然后从新获取然后作比较
}
