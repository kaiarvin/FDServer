// ChatServer
package ChatServer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Client struct {
	AccountID    int
	Name         string
	Icon         int
	Level        int
	Area         int
	Viplevel     int
	RecMsg       chan string
	ChatChannlID int
	UserRelation
}

type UserRelation struct {
	FriendInfo map[uint64]string
	IgnorID    map[uint64]string
}

type ChatServer struct {
	ID         int               //服务器ID
	UserList   map[uint64]Client //用户列表
	ChannlList map[int]chan byte //通道列表
	DataChannl chan byte         //中转接受数据
}

func processLine(line string) {
	s := string(line)
	//os.Stdout.Write(line)
	fmt.Println(s)
}

func (this *ChatServer) ReadConf(name string, handler func(string)) (err error) {
	f, err := os.Open(name)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	buf := bufio.NewReader(f)

	for {
		//line, err := buf.ReadBytes('\n')
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func (this *ChatServer) InitServer() {}

func (this *ChatServer) InitDB(addr string, host string, DBname string) {}

func (this *ChatServer) ClientMsgProcess() {}

func (this *ChatServer) ServerMsgProcess() {}
