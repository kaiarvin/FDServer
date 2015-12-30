// ChatServer
package FDServer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"qbs"
	"strings"
)

const (
	Const_Loger_sys    = 0
	COnst_Loger_Player = 1
)
const (
	Const_Log_NoLog        = 0
	Const_Log_level_Debug  = 1
	Const_Log_level_Realse = 2
	Const_Log_CanWrite     = Const_Log_level_Realse
)

type Server struct {
	ID              int                  //服务器ID
	Port            int                  //监听Client端口号
	Host            string               //host地址
	UserList        map[net.Conn]*Client //用户列表
	AccountList     map[int]net.Conn     //account 列表
	NameList        map[string]net.Conn  //name列表
	NameAccountList map[string]int       //account&name列表
	DataChannl      chan byte            //中转接受数据
	ChatConfigData  map[string]string
	RecDataSize     uint64
	AckDataSize     uint64
	TatleTime       uint64
	DB              *qbs.Qbs //character数据库
	IsCloseServer   bool
	Loger           []*FDLog
}

func (this *Server) processConf(args []string) {
	if len(args) != 2 {
		return
	} else {
		this.ChatConfigData[args[0]] = args[1]
	}
	//	fmt.Println(args, this.ChatConfigData)
}

func (this *Server) ReadConf(name string) (err error) {
	this.initChatServerData()
	f, err := os.Open(name)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		buf := strings.Split(line, " = ")

		this.processConf(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func (this *Server) initChatServerData() {
	this.ID = 0
	this.Port = 0
	this.Host = ""
	this.DataChannl = make(chan byte, 1024)
	this.ChatConfigData = make(map[string]string, 0)
	this.UserList = make(map[net.Conn]*Client, 1000)
	this.NameAccountList = make(map[string]int, 1000)
	this.AccountList = make(map[int]net.Conn, 1000)
	this.NameList = make(map[string]net.Conn, 1000)
	this.IsCloseServer = false
	this.Loger = make([]*FDLog, 10, 20)
	//this.ChannlList = make(map[int]chan byte, 0)
}

func (this *Server) InitServer() error {
	err := this.ReadConf("./Config/Server.conf")
	if err != nil {
		fmt.Println(err)
		return err
	}

	serverhost := ":" + this.ChatConfigData["ChatPort"]

	server, err := net.Listen("tcp4", serverhost)

	defer server.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}

	DBType := this.ChatConfigData["ChatDBType"]
	DBUser := this.ChatConfigData["ChatDBUser"]
	DBPw := this.ChatConfigData["ChatDBPw"]
	DBHost := this.ChatConfigData["ChatDBHost"]
	DBPort := this.ChatConfigData["ChatDBPort"]
	DBName := this.ChatConfigData["ChatDBName"]
	DBParam := this.ChatConfigData["ChatDBParam"]

	dberr := this.initQbs(DBType, DBUser, DBPw, DBHost, DBPort, DBName, DBParam)
	if dberr != nil {
		fmt.Println(dberr)
		return err
	}
	defer this.GetServerQbs().Close()

	this.InitLoger()
	//监听连接
	go func() {
		for {
			fmt.Println("Wait Accept.")
			client_socket, err := server.Accept()
			if err != nil {
				return
			}
			this.newClient(&client_socket)
			fmt.Println("Accept success")
		}

	}()

	for {
		if this.IsCloseServer {
			return nil
		}
		fmt.Println("Input Command:\n gm  . Game Command\n sys . Server Command\n")
		r := bufio.NewReader(os.Stdin)
		line, _, _ := r.ReadLine()
		cmd := string(line)
		switch cmd {
		case "gm":
		case "sys":
			{
				fmt.Println("Input Server Command:")
				r = bufio.NewReader(os.Stdin)
				line, _, _ = r.ReadLine()
				cp := strings.ToUpper(string(line))
				if cp == "CLOSE" {
					this.IsCloseServer = true
				} else if cp == "LIST" {
					fmt.Println("LIST:", this.UserList)
				}
			}
		case "show":
			{
				var ByteSize float64
				ByteSize = float64(this.RecDataSize / this.TatleTime)
				ByteSize /= 1024
				fmt.Printf("接收数据:%.3f kb/s \n", ByteSize)
				ByteSize = float64(this.AckDataSize / this.TatleTime)
				ByteSize /= 1024
				fmt.Printf("发送数据:%.3f kb/s \n", ByteSize)
			}
		default:

		}

	}
	fmt.Println("End......")

	return nil
}

func (this *Server) newClient(n *net.Conn) {
	client := &Client{IsLive: true, Server: this, Client_Socket: *n}
	client.RecMsg = make(chan string, Const_MsgBody_Len)
	client.AckMsg = make(chan string, Const_MsgBody_Len)
	client.RecMsgByte = make([]byte, 0, Const_MsgBody_Len)
	client.AckMsgByte = make([]byte, 0, Const_MsgBody_Len)

	this.UserList[*n] = client
	//fmt.Println(this.UserList)
	client.IsLive = true
	//	fmt.Println(client)
	//	fmt.Println("UserList: ", len(this.UserList))

	//创建接受发送线程
	go func() {
		var buf []byte
		for {
			if this.IsCloseServer || !client.IsLive {
				this.CleanExitUser(&client.Client_Socket)
				client.Client_Socket.Close()
				close(client.RecMsg)
				close(client.AckMsg)
				return
			}

			//fmt.Println("Socket Waiting Read.")
			buf = make([]byte, Const_MsgBody_Len)
			n, err := client.Client_Socket.Read(buf)

			if err != nil {
				if err != io.EOF {
					fmt.Println("Read Msg:", err)
					client.IsLive = false
					continue
				}
			}

			if n != 0 {
				if len(client.RecMsgByte) != 0 {
					client.RecMsgByte = append(client.RecMsgByte, buf[:n]...)
				} else {
					client.RecMsgByte = buf[:n]
				}
				fmt.Println("Rec len: ", n)
				client.Server.RecDataSize += uint64(n)
				client.ClientRecMsgProcess()
			} else {
				fmt.Println("Rec len: ", n)
				fmt.Println("Read Msg:", err)
				client.IsLive = false
				continue
			}
		}
	}()

}

func (this *Server) initQbs(dbtype, dbuser, pw, dbhost, dbport, dbname, param string) error {
	dsn := dbuser + ":" + pw + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?" + param

	q, err := DBInit(dbtype, dsn, dbname)

	if nil != err {
		fmt.Println(err)
		return err
	}
	this.DB = q
	return nil
}

func (this *Server) GetServerQbs() *qbs.Qbs {
	return this.DB
}

func (this *Server) CleanExitUser(conn *net.Conn) {
	cl := this.UserList[*conn]
	fmt.Println(cl.Name, " EXIT")
	ack := new(ReqUserLogout)
	ack.Id = E_ACK_EXIT
	ack.AccountID = cl.AccountID
	ack.Name = cl.Name
	delete(this.AccountList, cl.AccountID)
	delete(this.NameList, cl.Name)
	delete(this.NameAccountList, cl.Name)
	delete(this.UserList, *conn)

	for _, client := range this.UserList {
		client.SendToSlef(ack)
	}
}

func (this *Server) WritList(cl *Client) {
	this.AccountList[cl.AccountID] = cl.Client_Socket
	this.NameList[cl.Name] = cl.Client_Socket
	this.NameAccountList[cl.Name] = cl.AccountID
}

func (this *Server) InitLoger() {
	fmt.Println("InitLoger")
	this.Loger[Const_Loger_sys] = new(FDLog)
	this.Loger[Const_Loger_sys].InitFDLog(Const_Log_level_Realse, "System.log")
}
